package main

import (
	atc_mid_health_check "github.com/ARMmaster17/atc-mid-health-check/mhc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

// main Entry point of the application. Handles core services before handing off to ServiceBase.
func main() {
	initLogger()
	viper.SetEnvPrefix("MHC")
	viper.AutomaticEnv()
	atc_mid_health_check.StartServiceBase()
}

// initLogger Initializes the logging platform using ZeroLog. Outputs to both the command line using a JSON schema
// and to a syslog file located at LogLocation. If the specified log file location, initialization will fail silently
// and return a console-only logger.
func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(atc_mid_health_check.LogLevel)
	logFile, err := os.OpenFile(atc_mid_health_check.LogLocation, os.O_RDWR, 0644)
	if err != nil {
		log.Warn().Msgf("unable to open '%s':\n%v", atc_mid_health_check.LogLocation, err)
		atc_mid_health_check.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		return
	}
	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	atc_mid_health_check.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
