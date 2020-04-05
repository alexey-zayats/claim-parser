package model

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

// Head ...
type Head struct {
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
	Region      string
	Kind        string
	Name        string
	Address     string
	INN         string
	Head        Head
	Cars        []Car
	Agreement   string
	Reliability string
	Reason      *string
	Valid       bool
}
