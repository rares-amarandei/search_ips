package geoip

import (
	"io"
)

// Extended Infor about Ip.
// Full Documentation avaialble at: https://www.ip2location.com/database
type ExtendedInfo struct {
	// Original IP Address(v6 or v4).
	Ip string
	// Country.
	Country string
	// Two-character country code based on ISO 3166.
	CountryCode string
	// Region name.
	Region string
	// City name.
	City string
	// ISP name.
	Isp string
	// Net Speed(e.g. DSL, T1, etc.).
	NetSpeed string
	// Mobile carrier.
	Mobile string
	// Organization/business unit(e.g. EDU, GOV, etc.)
	UsageType string
}

// Info holds geo information about IP.
type Info struct {
	// Original IP address.
	IP string
	// City.
	City string
	// Country.
	Country string
	// ISO country code.
	CountryCode string
	// Region name.
	Region string
	// ISP name.
	ISP string
}

// Provider defines generic geoip provider.
type Provider interface {
	// Get geoip information by IP.
	Lookup(IP string) (Info, error)

	// Get extendend information from geoip databases(ip2location, ip2proxy).
	LookupExtended(IP string) (ExtendedInfo, error)

	io.Closer
}
