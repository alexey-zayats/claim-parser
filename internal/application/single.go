package application

import "time"

// Single ..
type Single struct {
	Dirty             bool
	DistrictID        int64
	PassType          int
	Title             string
	Address           string
	Inn               string
	Ogrn              string
	CeoName           string
	CeoPhone          string
	CeoEmail          string
	ActivityKind      int64
	Agreement         int
	Reliability       int
	Passes            []Pass
	CityFrom          string
	CityTo            string
	AddressDest       string
	DateFrom          time.Time
	DateTo            time.Time
	OtherReason       string
	WhoNeedsHelp      string
	WhoNeedsHelpPhone string
	DocLinks          string
}
