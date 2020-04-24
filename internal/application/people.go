package application

// People ..
type People struct {
	DistrictID   int64
	PassType     int
	Title        string
	Address      string
	Inn          int64
	Ogrn         int64
	CeoName      string
	CeoPhone     string
	CeoEmail     string
	ActivityKind int64
	Agreement    int
	Reliability  int
	Passes       []FIO
}
