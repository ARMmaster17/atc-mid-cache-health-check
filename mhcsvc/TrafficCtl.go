package mhcsvc

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/exec"
	"sync"
)

var (
	trafficControlMU sync.Mutex
)

// Init Initializes the mutex object and copies the given logger for use by the module.
func Init(svcLogger zerolog.Logger) {
	trafficControlMU = sync.Mutex{}
}

// ExecuteTrafficCtlCommand Runs the given command using the user-provided path to traffic_ctl. Returns STDOUT upon success,
// or an error if an error occured. Function will block if it does not have
func ExecuteTrafficCtlCommand(subCommand string, printOutput bool) (string, error) {
	if subCommand == "" {
		return "", fmt.Errorf("empty command, nothing to run")
	}
	Logger.Debug().Msg("obtaining lock on traffic_ctl executable...")
	trafficControlMU.Lock()
	Logger.Debug().Msg("lock obtained on traffic_ctl executable")
	defer func() {
		Logger.Debug().Msg("removing lock on traffic_ctl executable")
		trafficControlMU.Unlock()
	}()
	tctlPath := os.Getenv("TRAFFIC_CTL_DIR")
	cmd := fmt.Sprintf("%s/bin/traffic_ctl %s", tctlPath, subCommand)
	Logger.Debug().Str("cmd", cmd).Msg("invoking traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if printOutput {
		fmt.Printf("%s %v", string(out), err)
	}
	return string(out), nil
}
