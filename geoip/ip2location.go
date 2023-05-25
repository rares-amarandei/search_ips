package geoip

import (
	"github.com/ip2location/ip2location-go/v9"
	"github.com/pkg/errors"
)

// providerIP2location implements `Provider`.
type providerIP2location struct {
	ip2locationDb *ip2location.DB
}

func (l *providerIP2location) Close() error {
	if l.ip2locationDb != nil {
		l.ip2locationDb.Close()
	}

	return nil
}

func (l *providerIP2location) Lookup(IP string) (Info, error) {
	record, err := l.ip2locationDb.Get_all(IP)
	if err != nil {
		return Info{}, err
	}

	return Info{
		IP:          IP,
		City:        record.City,
		Country:     record.Country_long,
		CountryCode: record.Country_short,
		Region:      record.Region,
		ISP:         record.Isp,
	}, nil
}

func (l *providerIP2location) LookupExtended(IP string) (ExtendedInfo, error) {
	rLoc, err := l.ip2locationDb.Get_all(IP)
	if err != nil {
		return ExtendedInfo{}, err
	}

	return ExtendedInfo{
		Ip:          IP,
		CountryCode: rLoc.Country_short,
		Country:     rLoc.Country_long,
		Region:      rLoc.Region,
		City:        rLoc.City,
		Isp:         rLoc.Isp,
		NetSpeed:    rLoc.Netspeed,
		Mobile:      rLoc.Mobilebrand,
		UsageType:   rLoc.Usagetype,
	}, nil
}

func newIP2Location(ip2locationPath string) (Provider, error) {
	ip2locationDb, err := ip2location.OpenDB(ip2locationPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open ip2locationDb")
	}

	return &providerIP2location{
		ip2locationDb: ip2locationDb,
	}, nil
}
