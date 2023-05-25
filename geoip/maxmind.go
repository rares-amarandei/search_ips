package geoip

import (
	"net"

	mmdb "github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
)

// providerMaxmind implements `Provider`.
type providerMaxmind struct {
	db *mmdb.Reader
}

func (m *providerMaxmind) Close() error {
	if m.db != nil {
		return m.db.Close()
	}

	return nil
}

func (m *providerMaxmind) Lookup(IP string) (Info, error) {
	const defaultLocale = "en"

	netIP := net.ParseIP(IP)

	info := Info{}
	record, err := m.db.City(netIP)
	if err != nil {
		return info, errors.Wrapf(err, "lookup: %v", IP)
	}

	isp, err := m.db.ISP(netIP)
	if err == nil {
		info.ISP = isp.ISP
	}

	info.IP = IP
	info.City = record.City.Names[defaultLocale]
	info.Country = record.Country.Names[defaultLocale]
	info.CountryCode = record.Country.IsoCode
	// info.Region is not supported at the moment

	return info, nil
}

func (m *providerMaxmind) LookupExtended(IP string) (ExtendedInfo, error) {
	return ExtendedInfo{}, nil
}

func newMaxMind(dbname string) (Provider, error) {
	db, err := mmdb.Open(dbname)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open geodb")
	}

	return &providerMaxmind{db: db}, nil
}
