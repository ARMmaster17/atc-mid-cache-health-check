package atc_mid_health_check

import (
	"io"
	"io/ioutil"
	"net/http"
)

func getStatusFromTrafficMonitor() string {
	Logger.Debug().Str("url", trafficMonitors[0]+apiPath).Msg("connecting to TM")
	r, err := http.Get(trafficMonitors[0] + apiPath)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Str("url", trafficMonitors[0]+apiPath).Msg("could not connect to TM")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Fatal().Err(err).Stack().Caller().Msg("unable to close connection with TM")
		}
	}(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Fatal().Err(err).Stack().Caller().Msg("unable to read response from TM")
	}
	return string(body)
}
