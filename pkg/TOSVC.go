package atc_mid_health_check

import (
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
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
		Logger.Error().Msgf("unable to connect to %s", viper.GetString("TO_HOSTNAME"))
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
	hostList.Lock(viper.GetInt("TO_CHECK_INTERVAL") / 2)
	defer hostList.Unlock()
	for _, server := range response.Response {
		_, hostExists := hostList.Hosts[*server.HostName]
		if !hostExists {
			continue
		}
		updateCmd := ""
		if *server.Status == "ADMIN_DOWN" && hostList.Hosts[*server.HostName].Manual != "DOWN" {
			updateCmd = fmt.Sprintf("host down %s", server.FQDN)
		} else if *server.Status == "ALL" && hostList.Hosts[*server.HostName].Manual != "UP" {
			updateCmd = fmt.Sprintf("host up %s", server.FQDN)
		}
		if updateCmd != "" {
			updateCmds = append(updateCmds, updateCmd)
		}
	}
	return updateCmds
}

func toAuth() (*toclient.Session, error) {
	schema := "https"
	if viper.GetBool("TO_INSECURE") {
		schema = "http"
	}
	session, _, err := toclient.LoginWithAgent(
		fmt.Sprintf("%s://%s", schema, viper.GetString("TO_HOSTNAME")),
		viper.GetString("TO_USERNAME"),
		viper.GetString("TO_PASSWORD"),
		viper.GetBool("TO_INSECURE"),
		"MHC",
		false,
		time.Duration(viper.GetInt("TO_API_TIMEOUT")) * time.Second)
	return session, err
}
