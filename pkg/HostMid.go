package atc_mid_health_check

type HostMid struct {
	Hostname string `json:"hostname"`
	Type string `json:"type"`
	Available bool `json:"combined_available"`
	Manual string `json:"MANUAL"`
	FQDN string `json:"fqdn"`
}
