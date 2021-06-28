package atc_mid_health_check

import (
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/mid-health-check/pkg/TrafficCtl"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var (
	trafficMonitors []string
	apiPath         string

	LogLevel    zerolog.Level
	LogLocation = "/var/log/mid-health-check/mhc.log"
	Logger      zerolog.Logger
)

// StartServiceBase Entry point for ServiceBase. Manages all three check services and hostList updates.
func StartServiceBase() {
	Logger.Debug().Msg("setting up")
	Logger.Trace().Msg("initializing program data")
	initVars()
	Logger.Trace().Msg("initializing TrafficCtl module")
	TrafficCtl.Init(Logger)
	s, err := registerCronJobs()
	if err != nil {
		Logger.Fatal().Err(err).Msg("unable to register interval checks with go-cron")
		return
	}

	// Seed the list of mids before starting.
	Logger.Debug().Msg("getting list of hosts")
	HostListUpdate()
	if hostStatus == nil {
		Logger.Fatal().Msg("unable to seed mid list")
	}
	Logger.Debug().Msg("starting checks")
	s.StartAsync()

	///////////////////////////////////
	// OLD IMPLEMENTATION --- IGNORE //
	///////////////////////////////////
	rawTMResponse := getStatusFromTrafficMonitor()
	tmStatus := parseTrafficMonitorStatus(rawTMResponse)
	cmds := getTrafficMonitorStatus(hostStatus, tmStatus)
	for i, cmd := range cmds {
		Logger.Trace().Msgf("updating host status (%d/%d)", i+1, len(cmds))
		_, err := TrafficCtl.ExecuteCommand(cmd, true)
		if err != nil {
			Logger.Error().Err(err).Msgf("unable to run command %s (%d/%d)", cmd, i+1, len(cmds))
			return
		}
	}
}

// registerCronJobs Sets up interval jobs so that checks can be performed on a specific schedule. Resource locking
// and job overrun protection is handled by gocron.
func registerCronJobs() (*gocron.Scheduler, error) {
	Logger.Trace().Msg("setting up scheduled API checks")
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Minutes().Do(HostListUpdate)  // Reload list of mids
	if err != nil {
		return nil, err
	}
	_, err = s.Every(viper.GetInt("TCP_CHECK_INTERVAL")).Seconds().Do(nil /*TCP Check*/) // Ignore for now
	if err != nil {
		return nil, err
	}
	_, err = s.Every(viper.GetInt("TM_CHECK_INTERVAL")).Seconds().Do(CheckTMService)
	if err != nil {
		return nil, err
	}
	_, err = s.Every(viper.GetInt("TO_CHECK_INTERVAL")).Seconds().Do(nil /*TO Check*/)   // TODO: X
	if err != nil {
		return nil, err
	}
	return s, nil
}

// initVars Loads variables from the environment. Currently performs no validation checks on variable contents.
func initVars() {
	LogLevel = zerolog.Level(viper.GetInt("LOG_LEVEL"))
	trafficMonitors = strings.Split(viper.GetString("TM_HOSTS"), ",")
	apiPath = viper.GetString("TM_API_PATH")
}

func buildHostStatusStruct(fqdn string, statusLine string) (HostMid, error) {
	statusStruct := HostMid{}
	statusStruct.FQDN = fqdn
	var status string
	var active string
	var local string
	var manual string
	var selfDetect string
	tokensFound, err := fmt.Sscanf(statusLine, "HOST_STATUS_%s,ACTIVE:%s:0:0,LOCAL:%s:0:0,MANUAL:%s:0:0,SELF_DETECT:%s:0:0", &status, &active, &local, &manual, &selfDetect)
	if err != nil {
		Logger.Error().Err(err).Str("line", statusLine).Str("fqdn", fqdn).Msg("unable to parse traffic_ctl output")
		return HostMid{}, err
	}
	if tokensFound < 5 {
		Logger.Error().Str("line", statusLine).Str("fqdn", fqdn).Msgf("expected 5 tokens in traffic_ctl output, found %d", tokensFound)
		return HostMid{}, fmt.Errorf("traffic_ctl output is missing HostMid data")
	}
	return statusStruct, nil
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

func parseTrafficMonitorStatus(response string) map[string]map[string]string {
	var data = map[string]map[string]string{}
	Logger.Trace().Str("response", response).Msg("unmarshalling response from TM")
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Str("response", response).Msg("unable to parse response from TM")
	}
	return data
}
