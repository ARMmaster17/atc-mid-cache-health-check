package mid_health_check

import (
	"reflect"
	"testing"
)

func Test_getHostStatus(t *testing.T) {
	type args struct {
		trafficCtlStatus string
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

			got1 := getHostStatus(tArgs.trafficCtlStatus)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getHostStatus got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_buildHostStatusStruct(t *testing.T) {
	type args struct {
		fqdn       string
		statusLine string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := buildHostStatusStruct(tArgs.fqdn, tArgs.statusLine)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildHostStatusStruct got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_pollTrafficCtlStatus(t *testing.T) {
	tests := []struct {
		name string

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := pollTrafficCtlStatus()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("pollTrafficCtlStatus got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_getTrafficMonitorStatus(t *testing.T) {
	type args struct {
		hostStatus map[string]map[string]string
		tmStatus   map[string]map[string]string
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

			got1 := getTrafficMonitorStatus(tArgs.hostStatus, tArgs.tmStatus)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getTrafficMonitorStatus got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_executeUpdateCommands(t *testing.T) {
	type args struct {
		cmds []string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			executeUpdateCommands(tArgs.cmds)

		})
	}
}

func Test_getStatusFromTrafficMonitor(t *testing.T) {
	tests := []struct {
		name string

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := getStatusFromTrafficMonitor()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getStatusFromTrafficMonitor got1 = %v, want1: %v", got1, tt.want1)
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

func Test_initLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "initializes logger with no panic()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initLogger()

		})
	}
}
