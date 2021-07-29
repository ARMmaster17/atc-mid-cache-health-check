package mhcsvc

import (
	"github.com/stretchr/testify/assert"
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

func TestHostList_Merge(t *testing.T) {
	type args struct {
		Host HostMid
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *HostList
		inspect func(r *HostList, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args
	}{
		{
			name: "TestInsertHostInEmptyList",
			init: func(t *testing.T) *HostList {
				return &HostList{}
			},
			args: func(t *testing.T) args {
				return args{
					Host: HostMid{
						Hostname:   "test",
						Type:       "MID",
						Available:  true,
						Manual:     "UP",
						FQDN:       "test.example.com",
						Status:     "UP",
						Active:     "UP",
						Local:      "UP",
						SelfDetect: "UP",
						TOUp:       "UP",
						TMUp:       "UP",
					},
				}
			},
			inspect: func(r *HostList, t *testing.T) {
				assert.Equal(t, 1, len(r.Hosts))
			},
		},
		{
			name: "TestMergeExistingHost",
			init: func(t *testing.T) *HostList {
				hl := &HostList{}
				hl.Merge(HostMid{
					Hostname:   "test",
					Type:       "MID",
					Available:  true,
					Manual:     "DOWN",
					FQDN:       "test.example.com",
					Status:     "UP",
					Active:     "UP",
					Local:      "UP",
					SelfDetect: "UP",
					TOUp:       "UP",
					TMUp:       "UP",
				})
				return hl
			},
			args: func(t *testing.T) args {
				return args{
					Host: HostMid{
						Hostname:   "test",
						Type:       "MID",
						Available:  true,
						Manual:     "UP",
						FQDN:       "test.example.com",
						Status:     "UP",
						Active:     "UP",
						Local:      "UP",
						SelfDetect: "UP",
						TOUp:       "UP",
						TMUp:       "UP",
					},
				}
			},
			inspect: func(r *HostList, t *testing.T) {
				assert.Equal(t, 1, len(r.Hosts))
				assert.Equal(t, "UP", r.Hosts["test"].Manual)
			},
		},
		{
			name: "TestMergeExistingHostKeepsTOUp",
			init: func(t *testing.T) *HostList {
				hl := &HostList{}
				hl.Merge(HostMid{
					Hostname:   "test",
					Type:       "MID",
					Available:  true,
					Manual:     "DOWN",
					FQDN:       "test.example.com",
					Status:     "UP",
					Active:     "UP",
					Local:      "UP",
					SelfDetect: "UP",
					TOUp:       "UP",
					TMUp:       "UP",
				})
				return hl
			},
			args: func(t *testing.T) args {
				return args{
					Host: HostMid{
						Hostname:   "test",
						Type:       "MID",
						Available:  true,
						Manual:     "UP",
						FQDN:       "test.example.com",
						Status:     "UP",
						Active:     "UP",
						Local:      "UP",
						SelfDetect: "UP",
						TOUp:       "DOWN",
						TMUp:       "UP",
					},
				}
			},
			inspect: func(r *HostList, t *testing.T) {
				assert.Equal(t, 1, len(r.Hosts))
				assert.Equal(t, "UP", r.Hosts["test"].TOUp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			receiver.Merge(tArgs.Host)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

		})
	}
}
