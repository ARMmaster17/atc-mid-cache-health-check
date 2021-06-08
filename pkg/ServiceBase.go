package atc_mid_health_check

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

var (
	cmd             string
	trafficCtl      string
	trafficMonitors []string
	apiPath         string

	LogLevel    zerolog.Level
	LogLocation = "/var/log/mid-health-check/mhc.log"
	Logger      zerolog.Logger
)

func StartServiceBase() {
	initVars()
	trafficCtlStatus := pollTrafficCtlStatus()
	hostStatus := getHostStatus(trafficCtlStatus)
	rawTMResponse := getStatusFromTrafficMonitor()
	tmStatus := parseTrafficMonitorStatus(rawTMResponse)
	cmds := getTrafficMonitorStatus(hostStatus, tmStatus)
	executeUpdateCommands(cmds)
}

func initVars() {
	LogLevel = zerolog.Level(viper.GetInt("LOG_LEVEL"))
	trafficCtl = viper.GetString("TRAFFIC_CTL_DIR")
	cmd = fmt.Sprintf("%s/bin/traffic_ctl metric match host_status", trafficCtl)
	trafficMonitors = strings.Split(viper.GetString("TM_HOSTS"), ",")
	apiPath = viper.GetString("TM_API_PATH")
}

func getHostStatus(trafficCtlStatus string) map[string]map[string]string {
	var hostStatus = map[string]map[string]string{}
	for i, line := range strings.Split(trafficCtlStatus, "\n") {
		Logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmpLine := strings.Split(line, " ")
		fqdn := strings.Replace(tmpLine[0], "proxy.process.host_status.", "", -1)
		Logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		Logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got FQDN from traffic_ctl output")
		hostStatus[hostname] = buildHostStatusStruct(fqdn, tmpLine[1])
	}
	return hostStatus
}

func buildHostStatusStruct(fqdn string, statusLine string) map[string]string {
	statusStruct := make(map[string]string)
	statusStruct["fqdn"] = fqdn
	for j, s := range strings.Split(statusLine, ",") {
		Logger.Trace().Str("s", s).Str("tmp_line", statusLine).Msgf("processing substring #%d from traffic_ctl output", j)
		sSplit := strings.Split(s, ":")
		if len(sSplit) > 1 {
			Logger.Trace().Str("s", s).Msg("split occurred")
			statusStruct[sSplit[0]] = sSplit[1]
		} else {
			Logger.Trace().Str("s", s).Msg("split did not occur")
			statusStruct["STATUS"] = strings.Split(sSplit[0], "_")[2]
		}
	}
	return statusStruct
}

func pollTrafficCtlStatus() string {
	Logger.Trace().Str("cmd", cmd).Msg("executing traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Str("cmd", cmd).Msg("unable to execute traffic_ctl")
	}
	return string(out)
}

func getTrafficMonitorStatus(hostStatus map[string]map[string]string, tmStatus map[string]map[string]string) []string {
	var updateCmds []string
	for hostname, hostdata := range tmStatus {
		Logger.Trace().Str("hostname", hostname).Str("hostdata", fmt.Sprint(hostdata)).Msg("processing host")
		updateCmd := ""
		hType := hostdata["type"]
		Logger.Debug().Str("hostname", hostname).Str("type", hType).Msg("checking host type")
		available := hostdata["combined_available"]
		Logger.Debug().Str("hostname", hostname).Str("available", available).Msg("checking host availability")
		_, hostExists := hostStatus[hostname]
		if hType == "MID" && hostExists {
			Logger.Trace().Str("hostname", hostname).Msg("host is of type MID and exists")
			if available == "true" {
				Logger.Trace().Str("hostname", hostname).Msg("host is available")
				status := "UP"
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports UP, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					updateCmd = fmt.Sprintf("%s host up %s", trafficCtl, hostStatus[hostname]["fqdn"])
				}
			} else {
				status := "DOWN"
				Logger.Trace().Str("hostname", hostname).Msg("host is not available")
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					updateCmd = fmt.Sprintf("%s host down %s", trafficCtl, hostStatus[hostname]["fqdn"])
				}
			}
		}
		if updateCmd != "" {
			updateCmds = append(updateCmds, updateCmd)
		}
	}
	return updateCmds
}

func executeUpdateCommands(cmds []string) {
	for i, cmd := range cmds {
		Logger.Debug().Str("cmd", cmd).Msgf("invoking traffic_ctl (%d/%d)", i, len(cmds))
		out, err := exec.Command(cmd).Output()
		fmt.Printf("%s %s", out, err)
	}
}

func getStatusFromTrafficMonitor() string {
	Logger.Debug().Str("url", trafficMonitors[0]+apiPath).Msg("connecting to TM")
	r, err := http.Get(trafficMonitors[0] + apiPath)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Str("url", trafficMonitors[0]+apiPath).Msg("could not connect to TM")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Fatal().Err(err).Stack().Caller().Msg("unable to close connection with TM")
		}
	}(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Msg("unable to read response from TM")
	}
	return string(body)
}

func parseTrafficMonitorStatus(response string) map[string]map[string]string {
	var data = map[string]map[string]string{}
	Logger.Trace().Str("response", response).Msg("unmarshalling response from TM")
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Str("response", response).Msg("unable to parse response from TM")
	}
	return data
}
