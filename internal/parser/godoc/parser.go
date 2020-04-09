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

	i := 1
	rows.Next()

	numStyle, _ := f.NewStyle(`{"number_format":1}`)
	dateStyle, _ := f.NewStyle(`{"custom_number_format": "m.d.yyyy h:mm:ss"}`)

	for rows.Next() {

		i = i + 1

		//fmt.Printf("A%d\n", i)

		test := f.GetCellValue(sheetName, fmt.Sprintf("A%d", i))
		if len(test) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			dateAxis := fmt.Sprintf("A%d", i)
			f.SetCellStyle("Sheet1", dateAxis, dateAxis, dateStyle)

			claim := &model.Claim{
				Valid:   true,
				Event:   event,
				Created: excelDateToDate(f.GetCellValue(sheetName, fmt.Sprintf("A%d", i))),
			}

			claim.Company.Activity = f.GetCellValue(sheetName, fmt.Sprintf("B%d", i))
			claim.Company.Title = f.GetCellValue(sheetName, fmt.Sprintf("C%d", i))
			claim.Company.Address = f.GetCellValue(sheetName, fmt.Sprintf("D%d", i))

			innAxis := fmt.Sprintf("E%d", i)
			f.SetCellStyle("Sheet1", innAxis, innAxis, numStyle)
			claim.Company.INN = f.GetCellValue(sheetName, innAxis)

			fio := strings.Split(f.GetCellValue(sheetName, fmt.Sprintf("F%d", i)), " ")
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

			claim.Company.Head.Contact.Phone = f.GetCellValue(sheetName, fmt.Sprintf("G%d", i))
			claim.Company.Head.Contact.EMail = f.GetCellValue(sheetName, fmt.Sprintf("H%d", i))

			claim.Source = f.GetCellValue(sheetName, fmt.Sprintf("I%d", i))
			claim.Cars = parser.ParseCars(claim.Source)

			claim.Agreement = f.GetCellValue(sheetName, fmt.Sprintf("J%d", i))
			claim.Reliability = f.GetCellValue(sheetName, fmt.Sprintf("K%d", i))

			out <- claim
		}
	}

	out <- nil

	return nil
}
