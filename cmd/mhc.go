package main

import (
	atc_mid_health_check "atc-mid-health-check/mhcsvc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

// main Entry point of the application. Handles core services before handing off to ServiceBase.
func main() {
	if initLogger() != nil {
		log.Fatal().Msg("unable to start logging")
	}
	atc_mid_health_check.StartServiceBase()
}

// initLogger Initializes the logging platform using ZeroLog. Outputs to both the command line using a JSON schema
// and to a syslog file located at LogLocation. If the specified log file location, initialization will fail silently
// and return a console-only logger.
func initLogger() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevelConfig, _ := strconv.ParseInt(os.Getenv("MHC_LOG_LEVEL"), 10, 64)
	atc_mid_health_check.LogLevel = zerolog.Level(logLevelConfig)
	zerolog.SetGlobalLevel(atc_mid_health_check.LogLevel)
	var multi zerolog.LevelWriter
	if os.Getenv("MHC_USE_LOGFILE") == "TRUE" {
		logFile, err := os.OpenFile(atc_mid_health_check.LogLocation, os.O_RDWR, 0644)
		if err != nil {
			log.Warn().Msgf("unable to open '%s':\n%v", atc_mid_health_check.LogLocation, err)
			atc_mid_health_check.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
			return err
		}
		multi = zerolog.MultiLevelWriter(logFile, os.Stdout)
	} else {
		multi = zerolog.MultiLevelWriter(os.Stdout)
	}
	atc_mid_health_check.Logger = zerolog.New(multi).With().Timestamp().Logger()
	return nil
}
