package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersions(t *testing.T) {
	{
		releaseVersionString := "Hugo Static Site Generator v0.36.1 darwin/amd64 BuildDate:"
		version, err := ParseAndCheckVersion([]byte(releaseVersionString))
		require.NoError(t, err, "ParseAndCheckVersion should not fail on a release input version string")
		require.Equal(t, version, "0.36.1", "ParseAndCheckVersion should return 0.36.1")
	}
	{
		debugVersionString := "Hugo Static Site Generator v0.33-DEV darwin/amd64 BuildDate:"
		version, err := ParseAndCheckVersion([]byte(debugVersionString))
		require.NoError(t, err, "ParseAndCheckVersion should not fail on a debug input version string")
		require.Equal(t, version, "0.33", "ParseAndCheckVersion should return 0.33")
	}
}
