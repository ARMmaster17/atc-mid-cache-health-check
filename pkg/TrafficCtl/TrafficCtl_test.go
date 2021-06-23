package TrafficCtl

import (
	"github.com/magiconair/properties/assert"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestInit(t *testing.T) {
	type args struct {
		svcLogger zerolog.Logger
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

			Init(tArgs.svcLogger)

		})
	}
}

func TestExecuteCommand(t *testing.T) {
	type args struct {
		subCommand  string
		printOutput bool
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "invalid subcommand",
			args: func(t *testing.T) args {
				return args{
					subCommand: "",
					printOutput: false,
				}
			},
			want1: "",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, err.Error(), "empty command, nothing to run")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := ExecuteCommand(tArgs.subCommand, tArgs.printOutput)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ExecuteCommand got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExecuteCommand error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
