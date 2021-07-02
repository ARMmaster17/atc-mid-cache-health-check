package mhc_svc

import (
	"strings"
	"sync"
)

// HostList Singleton object that handles an updated list of all mids that are tracked by traffic_ctl. A lockable
// singleton class is used to ensure that all three checking services are using the same data, and do not attempt
// to perform the same operation more than once.
type HostList struct {
	Hosts map[string]HostMid
	mu    *sync.Mutex
}

// Lock Blocks current thread until exclusive access can be given to the referenced
// HostList object. Only needed for manual writes to the Hosts field.
func (hl *HostList) Lock(timeout int) {
	if hl.mu == nil {
		hl.mu = &sync.Mutex{}
	}
	lockMutex(hl.mu, timeout)
}

// Unlock Returns exclusive access to the application for another thread to claim. Does
// not check if caller owns the lock.
func (hl *HostList) Unlock() {
	hl.mu.Unlock()
}

// Refresh Gets a list of mid caches from the local instance of traffic_ctl. Stores the parsed
// information in the referenced HostList object. Suppresses all errors to ensure that go-cron does not
// stop sending scheduled jobs.
func (hl *HostList) Refresh(mutexTimeout int) {
	trafficCtlOutput, err := pollTrafficCtlStatus()
	if err != nil {
		Logger.Error().Err(err).Msg("unable to poll TrafficCtl")
		return
	}
	// Lock the array after we are sure we have data to work with from traffic_ctl.
	hl.Lock(mutexTimeout)
	defer hl.Unlock()
	// Reset the array so no old servers get left behind.
	hl.Hosts = map[string]HostMid{}
	for i, line := range strings.Split(trafficCtlOutput, "\n") {
		Logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmpLine := strings.Split(line, " ")
		fqdn := strings.Replace(tmpLine[0], "proxy.process.host_status.", "", -1)
		Logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		Logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got FQDN from traffic_ctl output")
		hl.Hosts[hostname], err = buildHostStatusStruct(fqdn, tmpLine[1])
		if err != nil {
			Logger.Error().Err(err).Msg("unable to parse TrafficCtl output")
			return
		}
	}
}
