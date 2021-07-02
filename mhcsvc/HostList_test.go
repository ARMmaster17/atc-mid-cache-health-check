package mhcsvc

import (
	"testing"
)

func TestHostList_Lock(t *testing.T) {
	type args struct {
		timeout int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *HostList
		inspect func(r *HostList, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args
	}{
		{
			name: "testLockMutex",
			init: func(t *testing.T) *HostList {
				return &HostList{}
			},
			inspect: func(r *HostList, t *testing.T) {
				
			},
			
			args: func(t *testing.T) args {
				return args{
					timeout: 5,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			receiver.Lock(tArgs.timeout)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}
			
			receiver.Unlock()

		})
	}
}

func TestHostList_Unlock(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *HostList
		inspect func(r *HostList, t *testing.T) //inspects receiver after test run

	}{
		{
			name: "testUnlockMutex",
			init: func(t *testing.T) *HostList {
				return &HostList{}
			},
			inspect: func(r *HostList, t *testing.T) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			receiver.Lock(5)
			receiver.Unlock()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}

func TestHostList_Refresh(t *testing.T) {
	type args struct {
		mutexTimeout int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *HostList
		inspect func(r *HostList, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			receiver.Refresh(tArgs.mutexTimeout)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}
