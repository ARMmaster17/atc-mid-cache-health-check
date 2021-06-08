package atc_mid_health_check

import (
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/mid-health-check/pkg/TrafficCtl"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	trafficMonitors []string
	apiPath         string

	LogLevel    zerolog.Level
	LogLocation = "/var/log/mid-health-check/mhc.log"
	Logger      zerolog.Logger

	tcpCheckInterval = 2
	tmCheckInterval = 10
	toCheckInterval = 30
)

func StartServiceBase() {
	initVars()
	TrafficCtl.Init(Logger)
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(HostListUpdate) // Reload list of mids
	s.Every(tcpCheckInterval).Seconds().Do(nil /*TCP Check*/)
	s.Every(tmCheckInterval).Seconds().Do(nil /*TM Check*/)
	s.Every(toCheckInterval).Seconds().Do(nil /*TO Check*/)

	// Seed the list of mids before starting.
	HostListUpdate()
	if hostStatus == nil {
		Logger.Fatal().Msg("unable to seed mid list")
	}
	s.StartAsync()

	///////////// Old main
	trafficCtlStatus, err := pollTrafficCtlStatus()
	if err != nil {
		// FATAL: Not handled here.
		return
	}
	rawTMResponse := getStatusFromTrafficMonitor()
	tmStatus := parseTrafficMonitorStatus(rawTMResponse)
	cmds := getTrafficMonitorStatus(hostStatus, tmStatus)
	for i, cmd := range cmds {
		Logger.Trace().Msgf("updating host status (%d/%d)", i + 1, len(cmds))
		TrafficCtl.ExecuteCommand(cmd, true)
	}
}

func initVars() {
	LogLevel = zerolog.Level(viper.GetInt("LOG_LEVEL"))
	trafficMonitors = strings.Split(viper.GetString("TM_HOSTS"), ",")
	apiPath = viper.GetString("TM_API_PATH")
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

func pollTrafficCtlStatus() (string, error) {
	return TrafficCtl.ExecuteCommand("metric match host_status", false)
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
					updateCmd = fmt.Sprintf("host up %s", hostStatus[hostname]["fqdn"])
				}
			} else {
				status := "DOWN"
				Logger.Trace().Str("hostname", hostname).Msg("host is not available")
				if hostStatus[hostname]["MANUAL"] != status {
					log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					updateCmd = fmt.Sprintf("host down %s", hostStatus[hostname]["fqdn"])
				}
			}
		}
		if updateCmd != "" {
			updateCmds = append(updateCmds, updateCmd)
		}
	}
	return updateCmds
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
