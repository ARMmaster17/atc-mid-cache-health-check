package atc_mid_health_check

import "strings"

var (
	hostStatus map[string]map[string]string
)

func HostListUpdate() {
	trafficCtlOutput, err := pollTrafficCtlStatus()
	if err != nil {
		Logger.Error().Err(err).Msg("unable to poll TrafficCtl")
		return
	}
	for i, line := range strings.Split(trafficCtlOutput, "\n") {
		Logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmpLine := strings.Split(line, " ")
		fqdn := strings.Replace(tmpLine[0], "proxy.process.host_status.", "", -1)
		Logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		Logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got FQDN from traffic_ctl output")
		hostStatus[hostname] = buildHostStatusStruct(fqdn, tmpLine[1])
	}
}
