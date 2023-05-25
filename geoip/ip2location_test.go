package geoip

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//
// Download databases from `https://lite.ip2location.com/ip2location-lite`
//

func TestIP2LocationProvider(t *testing.T) {
	const ip2locationPath = "test-data/IP2LOCATION-DB.BIN"
	const ip2proxyPath = "test-data/PX11.SAMPLE.BIN"

	var provider Provider

	ip2locationFile, err := os.OpenFile(ip2locationPath, os.O_RDONLY, 0)
	if err != nil {
		t.Skipf("SKIP: ip2location not found: %v", err)
	}
	ip2locationFile.Close()

	ip2proxyFile, err := os.OpenFile(ip2proxyPath, os.O_RDONLY, 0)
	if err != nil {
		t.Skipf("SKIP: ip2proxy not found: %v", err)
	}
	ip2proxyFile.Close()

	t.Run("Open", func(t *testing.T) {
		p, err := newIP2Location("someunknownpath", "someunknownpath")
		require.Error(t, err)
		require.Nil(t, p)

		provider, err = newIP2Location(ip2locationPath, ip2proxyPath)
		require.NoError(t, err)
		require.NotNil(t, provider)
	})

	t.Run("Lookup", func(t *testing.T) {
		IP := "81.2.69.160"

		inf, err := provider.Lookup(IP)
		require.NoError(t, err)
		require.Equal(t, IP, inf.IP)

		t.Logf("record = %+v", inf)

		require.Equal(t, "United Kingdom of Great Britain and Northern Ireland", inf.Country)
		require.Equal(t, "GB", inf.CountryCode)
		require.NotEmpty(t, inf.ISP)
	})

	t.Run("LookupExtended", func(t *testing.T) {
		IP := "81.2.69.160"

		inf, err := provider.LookupExtended(IP)
		require.NoError(t, err)
		require.Equal(t, IP, inf.Ip)

		t.Logf("record = %+v", inf)

		require.Equal(t, "United Kingdom of Great Britain and Northern Ireland", inf.Country)
		require.Equal(t, "GB", inf.CountryCode)
		require.Equal(t, inf.IsProxy, "0")
		require.NotEmpty(t, inf.Isp)
	})

	t.Run("Close", func(t *testing.T) {
		require.NoError(t, provider.Close())
		require.NoError(t, new(providerIP2location).Close())
	})
}
