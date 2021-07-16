package mhcsvc

import (
	"fmt"
	"regexp"
	"strings"
)

// HostMid Represents a single mid cache that is tracked by traffic_ctl. Acts as the authoritative state for that
// server to all three checking services. Assumes that the parent HostList handles locking and atomic transactions
// protections.
type HostMid struct {
	Hostname   string `json:"hostname"`
	Type       string `json:"type"`
	Available  bool   `json:"combined_available"`
	Manual     string `json:"MANUAL"`
	FQDN       string `json:"fqdn"`
	Status     string `json:"STATUS"`
	Active     string
	Local      string
	SelfDetect string
}

// buildHostStatusStruct Takes a line from the output of traffic_ctl and converts it to a HostMid object.
func buildHostStatusStruct(fqdn string, statusLine string) (HostMid, error) {
	if statusLine == "" {
		return HostMid{}, fmt.Errorf("statusLine is empty")
	}
	statusStruct := HostMid{}
	statusStruct.FQDN = fqdn
	statusStruct.Hostname = strings.Split(fqdn, ".")[0]
	re := regexp.MustCompile(`HOST_STATUS_(?P<STATUS>[a-zA-Z]+),ACTIVE:(?P<ACTIVE>[a-zA-Z]+):\d+:\d+,LOCAL:(?P<LOCAL>[a-zA-Z]+):\d+:\d+,MANUAL:(?P<MANUAL>[a-zA-Z]+):\d+:\d+,SELF_DETECT:(?P<SELF_DETECT>[a-zA-Z]+):\d+`)
	if !re.MatchString(statusLine) {
		Logger.Error().Str("line", statusLine).Str("fqdn", fqdn).Msgf("traffic_ctl does not match expected format")
		return HostMid{}, fmt.Errorf("traffic_ctl output is missing HostMid data")
	}
	matches := re.FindStringSubmatch(statusLine)
	statusStruct.Status = matches[re.SubexpIndex("STATUS")]
	statusStruct.Active = matches[re.SubexpIndex("ACTIVE")]
	statusStruct.Local = matches[re.SubexpIndex("LOCAL")]
	statusStruct.Manual = matches[re.SubexpIndex("MANUAL")]
	statusStruct.SelfDetect = matches[re.SubexpIndex("SELF_DETECT")]

	return statusStruct, nil
}
