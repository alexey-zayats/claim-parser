package entity

import "time"

// IssuedPeople ...
type IssuedPeople struct {
	ID             int64      `db:"id"`
	DistrictID     int64      `db:"district_id"`
	CompanyID      int64      `db:"company_id"`
	Lastname       string     `db:"lastname"`
	Firstname      string     `db:"firstname"`
	Patrname       string     `db:"patrname"`
	LegalBasement  string     `db:"legal_basement"`
	PassNumber     string     `db:"pass_number"`
	CreatedAt      *time.Time `db:"created_at"`
	CreatedBy      int64      `db:"created_by"`
	IssuedAt       *time.Time `db:"issued_at"`
	RegistryNumber int64      `db:"registry_number"`
	Shiping        int        `db:"shiping"`
	ArmNumber      string     `db:"arm_number"`
	ArmNumberBy    int        `db:"arm_number_by"`
	ArmNumberAt    *time.Time `db:"arm_number_at"`
}
