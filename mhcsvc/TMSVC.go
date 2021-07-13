package mhcsvc

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type TMServerStatus struct {
	Type string `json:"type"`
	Available bool `json:"combined_available"`
}

// CheckTMService Entry point for Traffic Monitor checks. Checks the current state of mid caches in HostList using
// Traffic Monitor. The local state and traffic_ctl will be updated to match the state returned by TM.
func CheckTMService() {
	rawTMResponse, err := getStatusFromTrafficMonitor()
	if err != nil {
		Logger.Error().Err(err).Stack().Caller().Msg("unable to sync with Traffic Monitor")
		// Do not return an error. Doing so may break go-cron.
		return
	}
	tmStatus, err := parseTrafficMonitorStatus(rawTMResponse)

	filteredTmStatus := filterCachesByMidType(tmStatus)
	cmds := checkForCacheStateChanges(filteredTmStatus)
	updateMidsInTrafficCtl(cmds)
}

// getStatusFromTrafficMonitor Connects to Traffic Monitor and returns the JSON-formatted response body with the mid
// cache status data.
func getStatusFromTrafficMonitor() ([]byte, error) {
	r, err := tryAllTrafficMonitors()
	if err != nil {
		return nil, fmt.Errorf("unable to get status from Traffic Monitor")
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Error().Err(err).Stack().Caller().Msg("unable to read response from TM")
		return nil, fmt.Errorf("the TM was response unreadable")
	}
	return body, nil
}

// tryAllTrafficMonitors Attempts to connect to each configured Traffic Monitor instance. On the first successful
// connection, the raw response object is returned.
func tryAllTrafficMonitors() (*http.Response, error) {
	for i, tm := range trafficMonitors {
		Logger.Debug().Str("url", tm+apiPath).Msgf("(%d/%d) connecting to TM instance %s", i+1, len(trafficMonitors), tm)
		r, err := http.Get(tm + apiPath)
		if err == nil {
			return r, nil
		}
	}
	Logger.Error().Msg("could not connect to any Traffic Monitor instance")
	return nil, fmt.Errorf("all traffic monitor instances are offline")
}

// parseTrafficMonitorStatus Converts the JSON payload from Traffic Monitor into a searchable struct map indexed by
// hostname.
func parseTrafficMonitorStatus(response []byte) (map[string]TMServerStatus, error) {
	var result map[string]TMServerStatus
	err := json.Unmarshal([]byte(response), &result)
	return result, err
}

// filterCachesByMidType Filters the output of Traffic Monitor and only returns items where the "type" key equals "MID".
func filterCachesByMidType(tmStatus map[string]TMServerStatus) map[string]TMServerStatus {
	Logger.Debug().Msg("filtering TM payload by MID status")
	var filteredTmStatus = make(map[string]TMServerStatus)
	Logger.Trace().Msg("obtaining lock on hostList")
	tmCheckInterval, _ := strconv.ParseInt(os.Getenv("MHC_TM_CHECK_INTERVAL"), 10, 64)
	hostList.Lock(int(tmCheckInterval) / 2)
	for hostname, hostdata := range tmStatus {
		_, hostExists := hostList.Hosts[hostname]
		if hostdata.Type == "MID" && hostExists {
			Logger.Trace().Str("hostname", hostname).Msg("host is of type MID and exists")
			filteredTmStatus[hostname] = hostdata
		}
	}
	Logger.Trace().Msg("releasing lock on hostList")
	hostList.Unlock()
	return filteredTmStatus
}

// checkForCacheStateChanges Compares each mid in the TM response map with the data in HostList. An array of subcommands
// are prepared for all mids that have mis-matched data to be run with traffic_ctl.
func checkForCacheStateChanges(tmStatus map[string]TMServerStatus) []string {
	var updateCmds []string
	tmCheckInterval, _ := strconv.ParseInt(os.Getenv("MHC_TM_CHECK_INTERVAL"), 10, 64)
	hostList.Lock(int(tmCheckInterval) / 2)
	// Locking at the method level because in Golang defer cannot be used inside a for loop.
	defer hostList.Unlock()
	for hostname, hostdata := range tmStatus {
		Logger.Trace().Str("hostname", hostname).Str("hostdata", fmt.Sprint(hostdata)).Msg("processing host")
		updateCmd := ""
		Logger.Debug().Str("hostname", hostname).Str("type", hostdata.Type).Msg("checking host type")
		Logger.Debug().Str("hostname", hostname).Bool("available", hostdata.Available).Msg("checking host availability")
		if hostdata.Available {
			Logger.Trace().Str("hostname", hostname).Msg("host is available")
			if hostList.Hosts[hostname].Manual != "UP" {
				log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports UP, Manual override is %s, Host Status is %s\n", hostname, hostList.Hosts[hostname].Manual, hostList.Hosts[hostname].Status)
				updateCmd = fmt.Sprintf("host up %s", hostList.Hosts[hostname].FQDN)
			}
		} else {
			Logger.Trace().Str("hostname", hostname).Msg("host is not available")
			if hostList.Hosts[hostname].Manual != "DOWN" {
				log.Info().Str("hostname", hostname).Msgf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostList.Hosts[hostname].Manual, hostList.Hosts[hostname].Status)
				updateCmd = fmt.Sprintf("host down %s", hostList.Hosts[hostname].FQDN)
			}
		}
		if updateCmd != "" {
			updateCmds = append(updateCmds, updateCmd)
		}
	}
	return updateCmds
}
