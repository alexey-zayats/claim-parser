package gsheet

import (
	"github.com/alexey-zayats/claim-parser/internal/parser"
)

// Register ...
func Register() {
	parser.Instance().Add("gsheet.vehicle", NewVehicleParser)
	parser.Instance().Add("gsheet.people", NewPeopleParser)
}
