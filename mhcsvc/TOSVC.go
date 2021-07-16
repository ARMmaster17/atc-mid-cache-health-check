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
	toc, err := toAuth()
	if err != nil {
		Logger.Error().Msgf("unable to connect to %s", os.Getenv("MHC_TO_HOSTNAME"))
		return tc.ServersV3Response{}, true
	}
	params := url.Values{}
	params.Add("type", "MID")
	response, info, err := toc.GetServersWithHdr(&params, nil)
	if err != nil || info.StatusCode != http.StatusOK {
		Logger.Error().Msg("unable to get list of MID servers from TO")
		return tc.ServersV3Response{}, true
	}
	return response, false
}

func getAdminDownMids(response tc.ServersV3Response) []string {
	var updateCmds []string
	toCheckInterval, _ := strconv.ParseInt(os.Getenv("MHC_TO_CHECK_INTERVAL"), 10, 64)
	hostList.Lock(int(toCheckInterval) / 2)
	defer hostList.Unlock()
	for _, server := range response.Response {
		_, hostExists := hostList.Hosts[*server.HostName]
		if !hostExists {
			continue
		}
		updateCmd := ""
		if *server.Status == "ADMIN_DOWN" && hostList.Hosts[*server.HostName].Manual != "DOWN" {
			Logger.Debug().Str("hostname", *server.HostName).Msg("manual is not DOWN, but TO reports server as ADMIN_DOWN")
			updateCmd = fmt.Sprintf("host down %s", *server.FQDN)
		} else if *server.Status == "ALL" && hostList.Hosts[*server.HostName].Manual != "UP" {
			Logger.Debug().Str("hostname", *server.HostName).Msg("manual is not UP, but TO reports server is not ADMIN_DOWN")
			updateCmd = fmt.Sprintf("host up %s", *server.FQDN)
		}
		if updateCmd != "" {
			updateCmds = append(updateCmds, updateCmd)
		}
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
		fmt.Sprintf("%s://%s", schema, os.Getenv("TO_HOSTNAME")),
		os.Getenv("TO_USERNAME"),
		os.Getenv("TO_PASSWORD"),
		os.Getenv("TO_INSECURE") == "TRUE",
		"MHC",
		false,
		time.Duration(toApiTimeout)*time.Second)
	return session, err
}
