package model

// AppKind ...
type AppKind int

const (
	// KindVehicle ...
	KindVehicle AppKind = iota
	// KindPeople ...
	KindPeople
)

// Pass ...
type Pass struct {
	Car        string
	Lastname   string
	Firstname  string
	Middlename string
}

// Application ..
type Application struct {
	Dirty        bool
	Kind         AppKind
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
	Passes       []Pass
}
