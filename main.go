package mid_health_check

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	cmd             = "/opt/trafficserver/bin/traffic_ctl metric match host_status"
	trafficCtl      = "/opt/trafficserver"
	trafficMonitors = [...]string{"tm.example.com"}
	apiPath         = "/api/cache-statuses"

	logLevel    = zerolog.InfoLevel
	logLocation = "/var/log/mid-health-check/mhc.log"
	logger      zerolog.Logger
)

func main() {
	initLogger()
	trafficCtlStatus := pollTrafficCtlStatus()
	hostStatus := getHostStatus(trafficCtlStatus)
	rawTMResponse := getStatusFromTrafficMonitor()
	tmStatus := parseTrafficMonitorStatus(rawTMResponse)
	getTrafficMonitorStatus(hostStatus, tmStatus)
}

func getHostStatus(trafficCtlStatus string) map[string]map[string]string {
	var hostStatus = map[string]map[string]string{}
	for i, line := range strings.Split(trafficCtlStatus, "\n") {
		logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmpLine := strings.Split(line, " ")
		fqdn := strings.Replace(tmpLine[0], "proxy.process.host_status.", "", -1)
		logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got FQDN from traffic_ctl output")
		hostStatus[hostname] = buildHostStatusStruct(fqdn, tmpLine[1])
	}
	return hostStatus
}

func buildHostStatusStruct(fqdn string, statusLine string) map[string]string {
	statusStruct := make(map[string]string)
	statusStruct["fqdn"] = fqdn
	for j, s := range strings.Split(statusLine, ",") {
		logger.Trace().Str("s", s).Str("tmp_line", statusLine).Msgf("processing substring #%d from traffic_ctl output", j)
		sSplit := strings.Split(s, ":")
		if len(sSplit) > 1 {
			logger.Trace().Str("s", s).Msg("split occurred")
			statusStruct[sSplit[0]] = sSplit[1]
		} else {
			logger.Trace().Str("s", s).Msg("split did not occur")
			statusStruct["STATUS"] = strings.Split(sSplit[0], "_")[2]
		}
	}
	return statusStruct
}

func pollTrafficCtlStatus() string {
	logger.Trace().Str("cmd", cmd).Msg("executing traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("cmd", cmd).Msg("unable to execute traffic_ctl")
	}
	return string(out)
}

func getTrafficMonitorStatus(hostStatus map[string]map[string]string, tmStatus map[string]map[string]string) {
	for hostname, hostdata := range tmStatus {
		logger.Trace().Str("hostname", hostname).Str("hostdata", fmt.Sprint(hostdata)).Msg("processing host")
		updateCmd := ""
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
					updateCmd = fmt.Sprintf("%s host up %s", trafficCtl, hostStatus[hostname]["fqdn"])
				}
			} else {
				status := "DOWN"
				logger.Trace().Str("hostname", hostname).Msg("host is not available")
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					updateCmd = fmt.Sprintf("%s host down %s", trafficCtl, hostStatus[hostname]["fqdn"])
				}
			}
			if updateCmd != "" {
				logger.Debug().Str("hostname", hostname).Str("cmd", updateCmd).Msg("a command has been generated to modify host status")
				out, err := exec.Command(updateCmd).Output()
				fmt.Printf("%s %s", out, err)
			}
		}
	}
	return
}

func getStatusFromTrafficMonitor() string {
	logger.Debug().Str("url", trafficMonitors[0]+apiPath).Msg("connecting to TM")
	r, err := http.Get(trafficMonitors[0] + apiPath)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("url", trafficMonitors[0]+apiPath).Msg("could not connect to TM")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Fatal().Err(err).Stack().Caller().Msg("unable to close connection with TM")
		}
	}(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Msg("unable to read response from TM")
	}
	return string(body)
}

func parseTrafficMonitorStatus(response string) map[string]map[string]string {
	var data = map[string]map[string]string{}
	logger.Trace().Str("response", response).Msg("unmarshalling response from TM")
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		logger.Fatal().Err(err).Stack().Caller().Str("response", response).Msg("unable to parse response from TM")
	}
	return data
}

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)
	logFile, err := os.OpenFile(logLocation, os.O_RDWR, 0644)
	if err != nil {
		log.Warn().Msgf("unable to open '%s':\n%w", logLocation, err)
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		return
	}
	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	logger = zerolog.New(multi).With().Timestamp().Logger()
}
