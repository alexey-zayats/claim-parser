package model

import (
	"time"
)

// FIO ...
type FIO struct {
	Surname    string
	Name       string
	Patronymic string
}

// Contact ...
type Contact struct {
	EMail string
	Phone string
}

// Person ...
type Person struct {
	FIO     FIO
	Contact Contact
}

// Car ...
type Car struct {
	Number string
	FIO    FIO
	Reason *string
	Valid  bool
}

// Company ...
type Company struct {
	Activity string
	Title    string
	Address  string
	INN      string
	Head     Person
}

// Claim ...
type Claim struct {
	Code        string
	Created     time.Time
	DistrictID  int64
	District    string
	Company     Company
	Cars        []Car
	Agreement   string
	Reliability string
	Reason      *string
	Valid       bool
	Source      string
	Event       *Event
}
