package geoip

// dummyProvider implements `Provider`.
type dummyProvider struct{}

func (*dummyProvider) Close() error {
	return nil
}

func (*dummyProvider) Lookup(IP string) (Info, error) {
	return Info{
		IP: IP,
	}, nil
}

func (*dummyProvider) LookupExtended(IP string) (ExtendedInfo, error) {
	return ExtendedInfo{
		Ip: IP,
	}, nil
}

func newDummy() (Provider, error) {
	return new(dummyProvider), nil
}
