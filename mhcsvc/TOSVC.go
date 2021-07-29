package mhcsvc

import (
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func CheckTOService() {
	response, err := getMidsFromTO()
	if err {
		return
	}
	updateCmds := getAdminDownMids(response)
	updateMidsInTrafficCtl(updateCmds)
}

func getMidsFromTO() (tc.ServersV3Response, bool) {
	Logger.Debug().Str("svc", "TOService").Msg("connecting to TO")
	toc, err := toAuth()
	if err != nil {
		Logger.Error().Err(err).Str("svc", "TOService").Msgf("unable to connect to %s", os.Getenv("MHC_TO_HOSTNAME"))
		return tc.ServersV3Response{}, true
	}
	params := url.Values{}
	params.Add("type", "MID")
	response, info, err := toc.GetServersWithHdr(&params, nil)
	if err != nil || info.StatusCode != http.StatusOK {
		Logger.Error().Err(err).Str("svc", "TOService").Msg("unable to get list of MID servers from TO")
		return tc.ServersV3Response{}, true
	}
	Logger.Trace().Str("svc", "TOService").Int("count", len(response.Response)).Msg("connected to TO")
	return response, false
}

func getAdminDownMids(response tc.ServersV3Response) []string {
	var updateCmds []string
	toCheckInterval, _ := strconv.ParseInt(os.Getenv("MHC_TO_CHECK_INTERVAL"), 10, 64)
	hostList.Lock(int(toCheckInterval) / 2)
	defer hostList.Unlock()
	for i, server := range response.Response {
		_, hostExists := hostList.Hosts[*server.HostName]
		Logger.Trace().Str("svc", "TOService").Str("to_hostname", *server.HostName).Msgf("(%d/%d) checking server from TO", i, len(response.Response))
		if !hostExists {
			Logger.Trace().Str("to_hostname", *server.HostName).Msg("ignoring, not in hostList")
			continue
		}
		Logger.Trace().Str("svc", "TOService").Str("hostname", *server.HostName).Str("status", *server.Status).Str("MANUAL", hostList.Hosts[*server.HostName].Manual).Msg("comparing server status")
		if *server.Status == "ADMIN_DOWN" {
			Logger.Debug().Str("svc", "TOService").Str("hostname", *server.HostName).Msg("manual is not DOWN, but TO reports server as ADMIN_DOWN")
			hostList.Hosts[*server.HostName].TOUp = "DOWN"
		} else if *server.Status == "REPORTED" {
			Logger.Debug().Str("svc", "TOService").Str("hostname", *server.HostName).Msg("manual is not UP, but TO reports server is not ADMIN_DOWN")
			hostList.Hosts[*server.HostName].TOUp = "UP"
		}
		updateCmds = append(updateCmds, hostList.Hosts[*server.HostName].CalculateCommand())
	}
	return updateCmds
}

func toAuth() (*toclient.Session, error) {
	schema := "https"
	if os.Getenv("MHC_TM_INSECURE") == "TRUE" {
		schema = "http"
	}
	toApiTimeout, _ := strconv.ParseInt(os.Getenv("MHC_TO_API_TIMEOUT"), 10, 64)
	session, _, err := toclient.LoginWithAgent(
		fmt.Sprintf("%s://%s", schema, os.Getenv("MHC_TO_HOSTNAME")),
		os.Getenv("MHC_TO_USERNAME"),
		os.Getenv("MHC_TO_PASSWORD"),
		os.Getenv("MHC_TO_INSECURE") == "TRUE",
		"MHC",
		false,
		time.Duration(toApiTimeout)*time.Second)
	return session, err
}
