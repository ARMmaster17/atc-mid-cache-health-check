package main

import (
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_initLogger(t *testing.T) {
	require.NoError(t, os.Setenv("MHC_USE_LOGFILE", "FALSE"))
	assert.NoError(t, initLogger())
}

func Test_GoCronShouldActuallyWork(t *testing.T) {
	// This makes sure nothing is broken with go-cron, because sometimes it acts weird on certain systems.
	s := gocron.NewScheduler(time.UTC)
	testInt := 0
	_, err := s.Every(1).Seconds().Do(func(){
		testInt++
	})
	require.NoError(t, err)
	s.StartAsync()
	time.Sleep(1500 * time.Millisecond)
	assert.NotEqual(t, 0, testInt)
	s.Stop()
}

func Test_GoCronShouldHandleZeros(t *testing.T) {
	// This makes sure nothing is broken with go-cron, because sometimes it acts weird on certain systems.
	s := gocron.NewScheduler(time.UTC)
	testInt := 0
	_, err := s.Every(0).Seconds().Do(func(){
		testInt++
	})
	require.Error(t, err)
}

func Test_GoCronShouldHandleImplOfJobFun(t *testing.T) {
	// This makes sure nothing is broken with go-cron, because sometimes it acts weird on certain systems.
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Seconds().Do(HelperMethod)
	require.NoError(t, err)
}

func HelperMethod() {
	time.Sleep(1 * time.Second)
}
