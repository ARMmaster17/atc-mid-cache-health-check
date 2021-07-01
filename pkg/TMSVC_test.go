package atc_mid_health_check

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
		response string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]map[string]string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := parseTrafficMonitorStatus(tArgs.response)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseTrafficMonitorStatus got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_filterCachesByMidType(t *testing.T) {
	type args struct {
		tmStatus map[string]map[string]string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]map[string]string
	}{
		{
			name: "testParsesNilMap",
			args: func(t *testing.T) args {
				return args{
					tmStatus: nil,
				}
			},
			want1: make(map[string]map[string]string),
		},
		{
			name: "testParsesEmptyMap",
			args: func(t *testing.T) args {
				return args{
					tmStatus: make(map[string]map[string]string),
				}
			},
			want1: make(map[string]map[string]string),
		},
		{
			name: "testParsesNonMidMap",
			args: func(t *testing.T) args {
				result := make(map[string]map[string]string)
				var midData = make(map[string]string)
				midData["type"] = "OTHER"
				result["testhost1"] = midData
				return args{
					tmStatus: result,
				}
			},
			want1: make(map[string]map[string]string),
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
					Hostname:   "testhost1",
					Type:       "MID",
					Available:  true,
					Manual:     "UP",
					FQDN:       "testhost1.example.com",
					Status:     "UP",
				}

				return args{
					tmStatus: map[string]map[string]string{
						"testhost1": map[string]string{
							"type": "MID"},
						},
				}
			},
			want1: map[string]map[string]string{
				"testhost1": map[string]string{
					"type": "MID"},
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
		tmStatus map[string]map[string]string
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
