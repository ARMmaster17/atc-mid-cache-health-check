package mid_health_check

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	cmd = "/opt/trafficserver/bin/traffic_ctl metric match host_status"
	traffic_ctl = "/opt/trafficserver"
	traffic_monitors = [...]string{"tm.example.com"}
	api_path = "/api/cache-statuses"

	logLevel = zerolog.InfoLevel
	logLocation ="/var/log/mid-health-check/mhc.log"
)

func main() {
	logger := getLogger()
	hostStatus := getHostStatus(logger)
	getTrafficMonitorStatus(hostStatus, logger)
}

func getHostStatus(logger zerolog.Logger) (map[string]map[string]string) {
	var host_status = map[string]map[string]string{}
	logger.Trace().Str("cmd", cmd).Msg("executing traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("cmd", cmd).Msg("unable to execute traffic_ctl")
	}
	for i, line := range strings.Split(string(out), "\n") {
		logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmp_line := strings.Split(line, " ")
		fqdn := strings.Replace(tmp_line[0], "proxy.process.host_status.", "", -1)
		logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got FQDN from traffic_ctl output")
		host_status[hostname] = make(map[string]string)
		host_status[hostname]["fqdn"] = fqdn
		for j, s := range strings.Split(tmp_line[1], ",") {
			logger.Trace().Str("s", s).Str("tmp_line", tmp_line[1]).Msgf("processing substring #%d from traffic_ctl output", j)
			sSplit := strings.Split(s, ":")
			if len(sSplit) > 1 {
				logger.Trace().Str("s", s).Msg("split occurred")
				host_status[hostname][sSplit[0]] = sSplit[1]
			} else {
				logger.Trace().Str("s", s).Msg("split did not occur")
				host_status[hostname]["STATUS"] = strings.Split(sSplit[0], "_")[2]
			}
		}
	}
	return host_status
}

func getTrafficMonitorStatus(hostStatus map[string]map[string]string, logger zerolog.Logger) {
	logger.Debug().Str("url", traffic_monitors[0] + api_path).Msg("connecting to TM")
	r, err := http.Get(traffic_monitors[0] + api_path)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("url", traffic_monitors[0] + api_path).Msg("could not connect to TM")
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Msg("unable to read response from TM")
	}
	var data = map[string]map[string]string{}
	logger.Trace().Str("response", string(body)).Msg("unmarshalling response from TM")
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("response", string(body)).Msg("unable to parse response from TM")
	}
	for hostname, hostdata := range data {
		logger.Trace().Str("hostname", hostname).Str("hostdata", fmt.Sprint(hostdata)).Msg("processing host")
		cmd = ""
		hType := hostdata["type"]
		logger.Debug().Str("hostname", hostname).Str("type", hType).Msg("checking host type")
		available := hostdata["combined_available"]
		logger.Debug().Str("hostname", hostname).Str("available", available).Msg("checking host availability")
		_, hostExists := hostStatus[hostname]
		if hType == "MID" && hostExists {
			logger.Trace().Str("hostname", hostname).Msg("host is of type MID and exists")
			if available == "true" {
				logger.Trace().Str("hostname", hostname).Msg("host is available")
				status := "UP"
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports UP, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					cmd = fmt.Sprintf("%s host up %s", traffic_ctl, hostStatus[hostname]["fqdn"])
				}
			} else {
				status := "DOWN"
				logger.Trace().Str("hostname", hostname).Msg("host is not available")
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					cmd = fmt.Sprintf("%s host down %s", traffic_ctl, hostStatus[hostname]["fqdn"])
				}
			}
			if cmd != "" {
				logger.Debug().Str("hostname", hostname).Str("cmd", cmd).Msg("a command has been generated to modify host status")
				out, err := exec.Command(cmd).Output()
				fmt.Printf("%s %s", out, err)
			}
		}
	}
	return
}

func getLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)
	logFile, err := os.OpenFile(logLocation, os.O_RDWR, 0644)
	if err != nil {
		log.Warn().Msgf("unable to open '%s':\n%w", logLocation, err)
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	return zerolog.New(multi).With().Timestamp().Logger()
}