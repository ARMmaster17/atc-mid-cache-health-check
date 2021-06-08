package main

import (
	atc_mid_health_check "github.com/ARMmaster17/mid-health-check/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func main() {
	initLogger()
	viper.SetEnvPrefix("MHC_")
	viper.AutomaticEnv()
	atc_mid_health_check.StartServiceBase()
}

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(atc_mid_health_check.LogLevel)
	logFile, err := os.OpenFile(atc_mid_health_check.LogLocation, os.O_RDWR, 0644)
	if err != nil {
		log.Warn().Msgf("unable to open '%s':\n%w", atc_mid_health_check.LogLocation, err)
		atc_mid_health_check.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		return
	}
	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	atc_mid_health_check.Logger = zerolog.New(multi).With().Timestamp().Logger()
}