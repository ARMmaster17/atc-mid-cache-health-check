package mhcsvc

import (
	"reflect"
	"testing"
)

func Test_buildHostStatusStruct(t *testing.T) {
	type args struct {
		fqdn       string
		statusLine string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      HostMid
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "testIgnoresEmptyLines",
			args: func(t *testing.T) args {
				return args{
					fqdn: "test.example.com",
					statusLine: "",
				}
			},
			want1: HostMid{},
			wantErr: true,
			inspectErr: nil,
		},
		{
			name: "testIgnoresInvalidLines",
			args: func(t *testing.T) args {
				return args{
					fqdn: "test.example.com",
					statusLine: "invalid",
				}
			},
			want1: HostMid{},
			wantErr: true,
			inspectErr: nil,
		},
		{
			name: "testBuildsValidUpLine",
			args: func(t *testing.T) args {
				return args{
					fqdn: "test.example.com",
					statusLine: "HOST_STATUS_UP,ACTIVE:UP:0:0,LOCAL:UP:0:0,MANUAL:UP:0:0,SELF_DETECT:UP:0:0",
				}
			},
			want1: HostMid{
				Hostname:   "test",
				Type:       "",
				Available:  false,
				Manual:     "UP",
				FQDN:       "test.example.com",
				Status:     "UP",
				Active:     "UP",
				Local:      "UP",
				SelfDetect: "UP",
			},
			wantErr: false,
			inspectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := buildHostStatusStruct(tArgs.fqdn, tArgs.statusLine)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildHostStatusStruct got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("buildHostStatusStruct error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Benchmark_buildHostStatusStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = buildHostStatusStruct("test.example.com", "HOST_STATUS_UP,ACTIVE:UP:0:0,LOCAL:UP:0:0,MANUAL:UP:0:0,SELF_DETECT:UP:0:0")
	}
}
