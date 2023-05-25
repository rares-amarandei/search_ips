package geoip

// New creates new GEOIP provider from given database.
func New(filename string) (Provider, error) {
	return newIP2Location(filename)

}
