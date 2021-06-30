package atc_mid_health_check

import (
	"github.com/ARMmaster17/mid-health-check/pkg/TrafficCtl"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
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

	hostList HostList
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
	hostList = HostList{}
	hostList.Refresh()
	if hostList.Hosts == nil {
		Logger.Fatal().Msg("unable to seed mid list")
	}
	Logger.Debug().Msg("starting checks")
	s.StartBlocking()
	Logger.Fatal().Msg("go-cron terminated unexpectedly")
}

// registerCronJobs Sets up interval jobs so that checks can be performed on a specific schedule. Resource locking
// and job overrun protection is handled by gocron.
func registerCronJobs() (*gocron.Scheduler, error) {
	Logger.Trace().Msg("setting up scheduled API checks")
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Minutes().Do(hostList.Refresh) // Reload list of mids
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
	_, err = s.Every(viper.GetInt("TO_CHECK_INTERVAL")).Seconds().Do(nil /*TO Check*/) // TODO: X
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

// pollTrafficCtlStatus Updates HostList with the latest mid cache data using traffic_ctl.
func pollTrafficCtlStatus() (string, error) {
	return TrafficCtl.ExecuteCommand("metric match host_status", false)
}
