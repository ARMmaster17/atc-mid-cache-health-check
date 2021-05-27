package mid_health_check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

var (
	cmd = "/opt/trafficserver/bin/traffic_ctl metric match host_status"
	traffic_ctl = "/opt/trafficserver"
	traffic_monitors = [...]string{"tm.example.com"}
	api_path = "/api/cache-statuses"
)

func main() {
	hostStatus, err := getHostStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = getTrafficMonitorStatus(hostStatus)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getHostStatus() (map[string]map[string]string, error) {
	var host_status = map[string]map[string]string{}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return nil, fmt.Errorf("unable to execute '%s':\n%w", cmd, err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		tmp_line := strings.Split(line, " ")
		fqdn := strings.Replace(tmp_line[0], "proxy.process.host_status.", "", -1)
		hostname := strings.Split(fqdn, ".")[0]
		host_status[hostname] = make(map[string]string)
		host_status[hostname]["fqdn"] = fqdn
		for _, s := range strings.Split(tmp_line[1], ",") {
			s_split := strings.Split(s, ":")
			if len(s_split) > 1 {
				host_status[hostname][s_split[0]] = s_split[1]
			} else {
				host_status[hostname]["STATUS"] = strings.Split(s_split[0], "_")[2]
			}
		}
	}
	return host_status, nil
}

func getTrafficMonitorStatus(hostStatus map[string]map[string]string) error {
	r, err := http.Get(traffic_monitors[0] + api_path)
	if err != nil {
		return fmt.Errorf("unable to GET '%s%s':\n%w", traffic_monitors[0], api_path, err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read response from Traffic Monitor:\n%w", err)
	}
	var data = map[string]map[string]string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("unable to parse response from Traffic Montior:\n%w", err)
	}
	for hostname, hostdata := range data {
		cmd = ""
		hType := hostdata["type"]
		available := hostdata["combined_available"]
		_, hostExists := hostStatus[hostname]
		if hType == "MID" && hostExists {
			if available == "true" {
				status := "UP"
				if hostStatus[hostname]["MANUAL"] != status {
					fmt.Printf("%s: Traffic Monitor reports UP, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					cmd = fmt.Sprintf("%s host up %s", traffic_ctl, hostStatus[hostname]["fqdn"])
				}
			} else {
				status := "DOWN"
				if hostStatus[hostname]["MANUAL"] != status {
					fmt.Printf("%s: Traffic Monitor reports DOWN, Manual override is %s, Host Status is %s\n", hostname, hostStatus[hostname]["MANUAL"], hostStatus[hostname]["STATUS"])
					cmd = fmt.Sprintf("%s host down %s", traffic_ctl, hostStatus[hostname]["fqdn"])
				}
			}
			if cmd != "" {
				fmt.Println(cmd)
				out, err := exec.Command(cmd).Output()
				fmt.Printf("%s %s", out, err)
			}
		}
	}
	return nil
}