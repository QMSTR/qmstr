package main

import (
	"testing"

	"github.com/QMSTR/qmstr/pkg/module/reporter/htmlreporter"
	"github.com/stretchr/testify/require"
)

func TestParseVersions(t *testing.T) {
	{
		releaseVersionString := "Hugo Static Site Generator v0.36.1 darwin/amd64 BuildDate:"
		version, err := htmlreporter.ParseVersion([]byte(releaseVersionString))
		require.NoError(t, err, "ParseAndCheckVersion should not fail on a release input version string")
		require.Equal(t, version, "0.36.1", "ParseAndCheckVersion should return 0.36.1")
	}
	{
		debugVersionString := "Hugo Static Site Generator v0.33-DEV darwin/amd64 BuildDate:"
		version, err := htmlreporter.ParseVersion([]byte(debugVersionString))
		require.NoError(t, err, "ParseAndCheckVersion should not fail on a debug input version string")
		require.Equal(t, version, "0.33", "ParseAndCheckVersion should return 0.33")
	}
}

func TestCheckMinimumRequiredVersion(t *testing.T) {
	require.NoError(t, htmlreporter.CheckMinimumRequiredVersion("0.33-DEV"), "Hugo version 0.33-DEV should be good enough")
	require.NoError(t, htmlreporter.CheckMinimumRequiredVersion("0.36.1"), "Hugo version 0.36.1 should be good enough")
	require.Error(t, htmlreporter.CheckMinimumRequiredVersion("0.31-DEV"), "Hugo version 0.31-DEV is too old")
	require.Error(t, htmlreporter.CheckMinimumRequiredVersion("0.31.32"), "Hugo version 0.32.32 is too old")
	require.Error(t, htmlreporter.CheckMinimumRequiredVersion("bring thee us a shrubbery"), "this is not a version")
}
