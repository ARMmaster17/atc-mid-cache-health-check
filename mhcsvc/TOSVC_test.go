package mhcsvc

import (
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

func TestCheckTOService(t *testing.T) {
	tests := []struct {
		name string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckTOService()

		})
	}
}

func Test_getMidsFromTO(t *testing.T) {
	tests := []struct {
		name string

		want1 tc.ServersV3Response
		want2 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := getMidsFromTO()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getMidsFromTO got1 = %v, want1: %v", got1, tt.want1)
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getMidsFromTO got2 = %v, want2: %v", got2, tt.want2)
			}
		})
	}
}

func Test_getAdminDownMids(t *testing.T) {
	type args struct {
		response tc.ServersV3Response
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

			got1 := getAdminDownMids(tArgs.response)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getAdminDownMids got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_toAuth(t *testing.T) {
	tests := []struct {
		name string

		want1      *toclient.Session
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := toAuth()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("toAuth got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("toAuth error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
