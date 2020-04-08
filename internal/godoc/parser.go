package godoc

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// Parser ...
type Parser struct {
}

// Name ...
const Name = "godoc"

// Register ...
func Register() {
	parser.Instance().Add(Name, NewParser)
}

// NewParser ...
func NewParser() (parser.Backend, error) {
	return &Parser{}, nil
}

var excelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

func excelDateToDate(excelDate string) time.Time {
	var days, _ = strconv.ParseFloat(excelDate, 64)
	return excelEpoch.Add(time.Second * time.Duration(days*86400))
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, param *dict.Dict, out chan interface{}) error {

	var path string

	if iface, ok := param.Get("path"); ok {
		path = iface.(string)
	} else {
		return fmt.Errorf("not found 'path' in param dict")
	}

	var event = &model.Event{}
	if iface, ok := param.Get("event"); ok {
		event = iface.(*model.Event)
	}

	claim := &model.Claim{
		Valid: true,
		Event: event,
	}

	logrus.WithFields(logrus.Fields{"name": Name, "path": path, "event": event}).Debug("godoc.Parse")

	f, err := excelize.OpenFile(path)
	if err != nil {
		return errors.Wrapf(err, "unable open xlsx file %s", path)
	}

	var sheetName string
	for _, sheet := range f.GetSheetMap() {
		sheetName = sheet
		break
	}

	rows, err := f.Rows(sheetName)
	if err != nil {
		return errors.Wrapf(err, "unable get rows by sheet %s", sheetName)
	}

	rows.Next()

	i := 1

	for rows.Next() {

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			for i, value := range rows.Columns() {

				if value == "" {
					break
				}

				switch i {
				case 0:
					claim.Created = excelDateToDate(value)
				case 1:
					claim.Company.Activity = value
				case 2:
					claim.Company.Title = value
				case 4:
					claim.Company.Address = value
				case 5:
					claim.Company.INN = value
				case 6:
					fio := strings.Split(value, "")
					if len(fio) == 3 {
						claim.Company.Head.FIO.Surname = fio[0]
						claim.Company.Head.FIO.Name = fio[1]
						claim.Company.Head.FIO.Patronymic = fio[2]
					} else if len(fio) == 2 {
						claim.Company.Head.FIO.Surname = fio[0]
						claim.Company.Head.FIO.Name = fio[1]
					} else if len(fio) == 1 {
						claim.Company.Head.FIO.Surname = fio[0]
					} else {
						claim.Valid = false
					}
				case 7:
					claim.Company.Head.Contact.Phone = value
				case 8:
					claim.Company.Head.Contact.EMail = value
				case 9:
					claim.Cars = parser.ParseCars(value)
				case 10:
					claim.Agreement = value
				case 11:
					claim.Reliability = value
				}
			}

			i++
			out <- claim
		}
	}

	out <- nil

	return nil
}
