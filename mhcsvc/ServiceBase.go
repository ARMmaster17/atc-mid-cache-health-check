package mhcsvc

import (
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const mutexLocked = 1

var (
	trafficMonitors []string
	apiPath         string

	LogLevel    zerolog.Level
	LogLocation = "/var/log/mhc.log"
	Logger      zerolog.Logger

	hostList HostList
)

// StartServiceBase Entry point for ServiceBase. Manages all three check services and hostList updates.
func StartServiceBase() {
	Logger.Debug().Msg("setting up")
	Logger.Trace().Msg("initializing program data")
	initVars()
	Logger.Trace().Msg("initializing TrafficCtl module")
	Init(Logger)
	s, err := registerCronJobs()
	if err != nil {
		Logger.Fatal().Err(err).Stack().Msg("unable to register interval checks with go-cron")
		return
	}

	// Seed the list of mids before starting.
	Logger.Debug().Msg("getting list of hosts")
	hostList = HostList{}
	hostList.Refresh(1)
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
	_, err := s.Every(1).Minutes().Do(func() {
		hostList.Refresh(5)
	}) // Reload list of mids
	if err != nil {
		return nil, err
	}
	//_, err = s.Every(viper.GetInt("TCP_CHECK_INTERVAL")).Seconds().Do(nil /*TCP Check*/) // Ignore for now
	//if err != nil {
	//	return nil, err
	//}
	tmCheckInterval, err := strconv.ParseInt(os.Getenv("MHC_TM_CHECK_INTERVAL"), 10, 64)
	if err != nil {
		return nil, err
	}
	_, err = s.Every(int(tmCheckInterval)).Seconds().Do(CheckTMService)
	if err != nil {
		return nil, err
	}
	toCheckInterval, err := strconv.ParseInt(os.Getenv("MHC_TO_CHECK_INTERVAL"), 10, 64)
	if err != nil {
		return nil, err
	}
	_, err = s.Every(int(toCheckInterval)).Seconds().Do(CheckTOService)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// initVars Loads variables from the environment. Currently performs no validation checks on variable contents.
func initVars() {
	trafficMonitors = strings.Split(os.Getenv("MHC_TM_HOSTS"), ",")
	apiPath = os.Getenv("MHC_TM_API_PATH")
}

// pollTrafficCtlStatus Updates HostList with the latest mid cache data using traffic_ctl.
func pollTrafficCtlStatus() (string, error) {
	return ExecuteTrafficCtlCommand("metric match host_status", false)
}

// lockMutex Common function for locking a mutex. Will display a trace message if the mutex is currently locked,
// and will show an error if the wait exceeded the specified timeout.
func lockMutex(mu *sync.Mutex, timeout int) {
	state := reflect.ValueOf(mu).Elem().FieldByName("state")
	if state.Int()&mutexLocked == mutexLocked {
		Logger.Trace().Msg("mutex is locked, waiting for exclusive access")
	}
	start := time.Now()
	mu.Lock()
	elapsed := time.Since(start)
	if int(elapsed.Seconds()) >= timeout {
		Logger.Warn().Msgf("lock acquire method exceeded threshold. Expected %ds, got %ds", timeout, int(elapsed.Seconds()))
	}
}

func updateMidsInTrafficCtl(updateCmds []string) {
	for i, cmd := range updateCmds {
		Logger.Trace().Msgf("updating host status (%d/%d)", i+1, len(updateCmds))
		_, err := ExecuteTrafficCtlCommand(cmd, true)
		if err != nil {
			Logger.Error().Err(err).Msgf("unable to run command %s (%d/%d)", cmd, i+1, len(updateCmds))
		}
	}
}
