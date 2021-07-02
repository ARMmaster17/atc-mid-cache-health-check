package mhc

import "fmt"

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
	statusStruct := HostMid{}
	statusStruct.FQDN = fqdn // TODO: Process FQDN in this function only.
	var status string
	var active string
	var local string
	var manual string
	var selfDetect string
	tokensFound, err := fmt.Sscanf(statusLine, "HOST_STATUS_%s,ACTIVE:%s:0:0,LOCAL:%s:0:0,MANUAL:%s:0:0,SELF_DETECT:%s:0:0", &status, &active, &local, &manual, &selfDetect)
	if err != nil {
		Logger.Error().Err(err).Str("line", statusLine).Str("fqdn", fqdn).Msg("unable to parse traffic_ctl output")
		return HostMid{}, err
	}
	if tokensFound < 5 {
		Logger.Error().Str("line", statusLine).Str("fqdn", fqdn).Msgf("expected 5 tokens in traffic_ctl output, found %d", tokensFound)
		return HostMid{}, fmt.Errorf("traffic_ctl output is missing HostMid data")
	}
	return statusStruct, nil
}
