package model

import "time"

// VehicleClaim ...
type VehicleClaim struct {
	Code       string
	Created    time.Time
	DistrictID int64
	District   string

	Company Company

	Passes []VehiclePass

	Agreement   string
	Reliability string

	// исходные данные списка пропусков
	Source string

	Success bool
	Reason  []string
}
