package TrafficCtl

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os/exec"
	"sync"
)

var (
	mu sync.Mutex
	logger zerolog.Logger
)

func Init(svcLogger zerolog.Logger) {
	mu = sync.Mutex{}
	logger = svcLogger
}

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
