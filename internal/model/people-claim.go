package model

import "time"

// PeopleClaim ...
type PeopleClaim struct {
	Created    time.Time
	DistrictID int64
	District   string

	Company Company

	Passes []Pass

	Agreement   string
	Reliability string

	// исходные данные списка пропусков
	Source string

	Success bool
	Reason  []string
}
