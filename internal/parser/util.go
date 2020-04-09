package parser

import (
	"strconv"
	"time"
)

var excelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

// ExcelDateToDate ...
func ExcelDateToDate(excelDate string) time.Time {
	var days, _ = strconv.ParseFloat(excelDate, 64)
	return excelEpoch.Add(time.Second * time.Duration(days*86400))
}
