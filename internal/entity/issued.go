package entity

import "time"

// Issued ...
type Issued struct {
	ID             int64     `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	CreatedBy      int64     `db:"created_by"`
	CompanyInn     string    `db:"company_inn"`
	CompanyOgrn    string    `db:"company_ogrn"`
	CompanyName    string    `db:"company_name"`
	CompanyFio     string    `db:"company_fio"`
	CompanyCar     string    `db:"company_car"`
	LegalBasement  string    `db:"legal_basement"`
	PassNumber     string    `db:"pass_number"`
	PassType       int       `db:"pass_type"`
	District       string    `db:"district"`
	IssuedAt       time.Time `db:"issued_at"`
	RegistryNumber string    `db:"registry_number"`
	Shipping       int       `db:"shipping"`
	FileID         int64     `db:"file_id"`
}
