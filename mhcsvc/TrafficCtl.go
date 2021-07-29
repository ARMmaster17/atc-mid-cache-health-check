package mhcsvc

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/exec"
	"strings"
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
// or an error if an error occurred. Function will block if it does not have
func ExecuteTrafficCtlCommand(subCommand string, printOutput bool) (string, error) {
	if subCommand == "" {
		return "", nil
	}
	Logger.Debug().Msg("obtaining lock on traffic_ctl executable...")
	trafficControlMU.Lock()
	Logger.Debug().Msg("lock obtained on traffic_ctl executable")
	defer func() {
		Logger.Debug().Msg("removing lock on traffic_ctl executable")
		trafficControlMU.Unlock()
	}()
	tctlPath := os.Getenv("MHC_TRAFFIC_CTL_DIR")
	Logger.Debug().Str("cmd", fmt.Sprintf("%s/bin/traffic_ctl %s", tctlPath, subCommand)).Msg("invoking traffic_ctl")
	splitCmd := strings.Split(subCommand, " ")
	out, err := exec.Command("sudo", append([]string{fmt.Sprintf("%s/bin/traffic_ctl", tctlPath)}, splitCmd...)...).CombinedOutput()
	if printOutput {
		fmt.Printf("%s %v", string(out), err)
	}
	return string(out), nil
}
