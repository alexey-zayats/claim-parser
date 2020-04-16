package model

import "time"

// VehicleRegistry ...
type VehicleRegistry struct {
	CompanyInn     string
	CompanyOgrn    string
	CompanyName    string
	CompanyFio     string
	CompanyCar     string
	LegalBasement  string
	PassNumber     string
	District       string
	PassType       int
	IssuedAt       time.Time
	RegistryNumber string
	Shipping       int
	Success        bool
	Reason         []string
}
