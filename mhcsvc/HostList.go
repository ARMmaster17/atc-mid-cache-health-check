package mhcsvc

import (
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

// HostList Singleton object that handles an updated list of all mids that are tracked by traffic_ctl. A lockable
// singleton class is used to ensure that all three checking services are using the same data, and do not attempt
// to perform the same operation more than once.
type HostList struct {
	Hosts map[string]*HostMid
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
	hl.Hosts = map[string]*HostMid{}
	for i, line := range strings.Split(trafficCtlOutput, "\n") {
		if line == "" {
			Logger.Trace().Msg("ignoring empty line from traffic_ctl output")
			continue
		}
		Logger.Trace().Str("line", line).Msgf("processing line %d from traffic_ctl output", i)
		tmpLine := strings.Split(line, " ")
		if len(tmpLine) < 2 {
			log.Error().Str("line", line).Msg("traffic_ctl returned an invalid result")
			continue
		}
		fqdn := strings.Replace(tmpLine[0], "proxy.process.host_status.", "", -1)
		Logger.Debug().Str("line", line).Str("fqdn", fqdn).Msg("got FQDN from traffic_ctl output")
		hostname := strings.Split(fqdn, ".")[0]
		Logger.Debug().Str("line", line).Str("hostname", hostname).Msg("got hostname from traffic_ctl output")
		newHostStruct, err := buildHostStatusStruct(fqdn, tmpLine[1])
		if err != nil {
			Logger.Error().Err(err).Str("line", line).Msg("unable to parse TrafficCtl output line")
			continue
		}
		hl.Merge(newHostStruct)
	}
}

// Merge adds a host into the host list, while maintaining the previous TO and TM status if the host already exists.
func (hl *HostList) Merge(host HostMid) {
	if hl.Hosts == nil {
		hl.Hosts = map[string]*HostMid{}
	}
	if existingHost, hostExists := hl.Hosts[host.Hostname]; hostExists {
		host.TOUp = existingHost.TOUp
		host.TMUp = existingHost.TMUp
	}
	hl.Hosts[host.Hostname] = &host
}
