package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_initLogger(t *testing.T) {
	require.NoError(t, os.Setenv("MHC_USE_LOGFILE", "FALSE"))
	assert.NoError(t, initLogger())
}
