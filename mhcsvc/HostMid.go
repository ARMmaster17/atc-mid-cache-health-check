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
	TOUp       string
	TMUp       string
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

// CalculateCommand evaluates the current state of the HostMid and generates a command to set the correct state in
// traffic_ctl. If no change is necessary, an empty string is returned. This method is NOT thread-safe, and assumes
// the parent caller has locked the global host list.
func (h *HostMid) CalculateCommand() string {
	if h.TOUp == "UP" && h.TMUp == "UP" && h.Manual != "UP" {
		Logger.Info().Str("fqdn", h.FQDN).Msg("host is reported online in TO and TM, bringing up in traffic_ctl")
		return fmt.Sprintf("host up %s", h.FQDN)
	} else if (h.TOUp == "DOWN" || h.TMUp == "DOWN") && h.Manual != "DOWN" {
		Logger.Info().Str("fqdn", h.FQDN).Msg("host is reported down by at least one source, bringing down in traffic_ctl")
		return fmt.Sprintf("host down %s", h.FQDN)
	} else {
		Logger.Trace().Str("fqdn", h.FQDN).Msg("host status has not changed state, skipping...")
		return ""
	}
}
