package geoip

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaxMindProvider(t *testing.T) {
	var provider Provider

	t.Run("Open", func(t *testing.T) {
		p, err := newMaxMind("someunknownpath")
		require.Error(t, err)
		require.Nil(t, p)

		provider, err = newMaxMind("test-data/GeoIP2-City-Test.mmdb")
		require.NoError(t, err)
		require.NotNil(t, provider)
	})

	t.Run("Lookup", func(t *testing.T) {
		IP := "81.2.69.160"

		inf, err := provider.Lookup(IP)
		require.NoError(t, err)
		require.Equal(t, IP, inf.IP)

		t.Logf("record = %+v", inf)

		require.Equal(t, "London", inf.City)
		require.Equal(t, "United Kingdom", inf.Country)
		require.Equal(t, "GB", inf.CountryCode)
	})

	t.Run("Close", func(t *testing.T) {
		require.NoError(t, provider.Close())
		require.NoError(t, new(providerMaxmind).Close())
	})
}
