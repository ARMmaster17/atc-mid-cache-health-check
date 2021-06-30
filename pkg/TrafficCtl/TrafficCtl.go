package TrafficCtl

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os/exec"
	"sync"
)

var (
	mu     sync.Mutex
	logger zerolog.Logger
)

// Init Initializes the mutex object and copies the given logger for use by the module.
func Init(svcLogger zerolog.Logger) {
	mu = sync.Mutex{}
	logger = svcLogger
}

// ExecuteCommand Runs the given command using the user-provided path to traffic_ctl. Returns STDOUT upon success,
// or an error if an error occured. Function will block if it does not have
func ExecuteCommand(subCommand string, printOutput bool) (string, error) {
	if subCommand == "" {
		return "", fmt.Errorf("empty command, nothing to run")
	}
	logger.Debug().Msg("obtaining lock on traffic_ctl executable...")
	mu.Lock()
	logger.Debug().Msg("lock obtained on traffic_ctl executable")
	defer func() {
		logger.Debug().Msg("removing lock on traffic_ctl executable")
		mu.Unlock()
	}()
	tctlPath := viper.GetString("TRAFFIC_CTL_DIR")
	cmd := fmt.Sprintf("%s/bin/traffic_ctl %s", tctlPath, subCommand)
	logger.Debug().Str("cmd", cmd).Msg("invoking traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if printOutput {
		fmt.Printf("%s %v", string(out), err)
	}
	return string(out), nil
}
