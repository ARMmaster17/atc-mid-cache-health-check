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
	tctlPath := viper.GetString("TRAFFIC_CTL_DIR")
	cmd := fmt.Sprintf("%s/bin/traffic_ctl %s", tctlPath, subCommand)
	logger.Debug().Str("cmd", cmd).Msg("invoking traffic_ctl")
	out, err := exec.Command(cmd).Output()
	if printOutput {
		fmt.Printf("%s %s", out, err)
	}
	return string(out), nil
}
