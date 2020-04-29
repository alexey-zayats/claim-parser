package application

// Vehicle ..
type Vehicle struct {
	Dirty        bool
	DistrictID   int64
	PassType     int
	Title        string
	Address      string
	Inn          string
	Ogrn         string
	CeoName      string
	CeoPhone     string
	CeoEmail     string
	ActivityKind int64
	Agreement    int
	Reliability  int
	Passes       []Pass
}
