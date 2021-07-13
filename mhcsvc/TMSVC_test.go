package mhcsvc

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCheckTMService(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckTMService()

		})
	}
}

func Test_getStatusFromTrafficMonitor(t *testing.T) {
	tests := []struct {
		name string

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := getStatusFromTrafficMonitor()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getStatusFromTrafficMonitor got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("getStatusFromTrafficMonitor error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_tryAllTrafficMonitors(t *testing.T) {
	tests := []struct {
		name string

		want1      *http.Response
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := tryAllTrafficMonitors()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("tryAllTrafficMonitors got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("tryAllTrafficMonitors error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_parseTrafficMonitorStatus(t *testing.T) {
	type args struct {
		response []byte
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]TMServerStatus

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "failureCaseReal1",
			args: func(t *testing.T) args {
				return args{
					response: []byte("{\"cmid1\": {\"type\": \"MID\",\"load_average\": 0,\"query_time_ms\": 1000,\"health_time_ms\": 1,\"stat_time_ms\": 1,\"stat_span_ms\": 1,\"health_span_ms\": 1000,\"status\": \"REPORTED - available\",\"status_poller\": \"health\",\"bandwidth_kbps\": 100,\"bandwidth_capacity_kbps\": 1000,\"connection_count\": 1,\"ipv4_available\": true,\"ipv6_available\": true,\"combined_available\": true,\"interfaces\": {\"eth0\": {\"status\": \"available\",\"status_poller\": \"health\",\"bandwidth_kbps\": 100,\"available\": true}}}}"),
				}
			},
			want1: map[string]TMServerStatus{
				"cmid1": {
					Type:      "MID",
					Available: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := parseTrafficMonitorStatus(tArgs.response)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseTrafficMonitorStatus got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("tryAllTrafficMonitors error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_filterCachesByMidType(t *testing.T) {
	type args struct {
		tmStatus map[string]TMServerStatus
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]TMServerStatus
	}{
		{
			name: "testParsesNilMap",
			args: func(t *testing.T) args {
				return args{
					tmStatus: nil,
				}
			},
			want1: make(map[string]TMServerStatus),
		},
		{
			name: "testParsesEmptyMap",
			args: func(t *testing.T) args {
				return args{
					tmStatus: make(map[string]TMServerStatus),
				}
			},
			want1: make(map[string]TMServerStatus),
		},
		{
			name: "testParsesNonMidMap",
			args: func(t *testing.T) args {
				return args{
					map[string]TMServerStatus{
						"testhost1": {
							Type:      "OTHER",
							Available: true,
						},
					},
				}
			},
			want1: make(map[string]TMServerStatus),
		},
		{
			name: "testParsesMidMap",
			args: func(t *testing.T) args {
				result := make(map[string]map[string]string)
				var midData = make(map[string]string)
				midData["type"] = "MID"
				result["testhost1"] = midData

				hostList = HostList{}
				hostList.Hosts = map[string]HostMid{}
				hostList.Hosts["testhost1"] = HostMid{
					Hostname:  "testhost1",
					Type:      "MID",
					Available: true,
					Manual:    "UP",
					FQDN:      "testhost1.example.com",
					Status:    "UP",
				}

				return args{
					tmStatus: map[string]TMServerStatus{
						"testhost1": {
							Type: "MID",
							Available: true,
						},
					},
				}
			},
			want1: map[string]TMServerStatus{
				"testhost1": {
					Type: "MID",
					Available: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := filterCachesByMidType(tArgs.tmStatus)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("filterCachesByMidType got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_checkForCacheStateChanges(t *testing.T) {
	type args struct {
		tmStatus map[string]TMServerStatus
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := checkForCacheStateChanges(tArgs.tmStatus)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("checkForCacheStateChanges got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
