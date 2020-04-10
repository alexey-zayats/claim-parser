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
	"strings"
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

// Parse ...
func (p *Parser) Parse(ctx context.Context, param *dict.Dict, out chan interface{}) error {

	var path string

	if iface, ok := param.Get("path"); ok {
		path = iface.(string)
	} else {
		return fmt.Errorf("not found 'path' in param dict")
	}

	var event = &model.Event{
		CreatedBy: 1,
	}
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

			axis := map[string]string{
				"created":         fmt.Sprintf("A%d", i),
				"activity":        fmt.Sprintf("B%d", i),
				"compaty-title":   fmt.Sprintf("C%d", i),
				"compaty-address": fmt.Sprintf("D%d", i),
				"compaty-inn":     fmt.Sprintf("E%d", i),
				"compaty-fio":     fmt.Sprintf("F%d", i),
				"compaty-phone":   fmt.Sprintf("G%d", i),
				"compaty-email":   fmt.Sprintf("H%d", i),
				"compaty-cars":    fmt.Sprintf("I%d", i),
				"agreement":       fmt.Sprintf("J%d", i),
				"reliability":     fmt.Sprintf("K%d", i),
				"ogrn":            fmt.Sprintf("L%d", i),
			}

			f.SetCellStyle("Sheet1", axis["created"], axis["created"], dateStyle)

			claim := &model.Claim{
				Valid:   true,
				Event:   event,
				Created: parser.ExcelDateToDate(f.GetCellValue(sheetName, axis["created"])),
			}

			claim.Company.Activity = f.GetCellValue(sheetName, axis["activity"])
			claim.Company.Title = f.GetCellValue(sheetName, axis["compaty-title"])
			claim.Company.Address = f.GetCellValue(sheetName, axis["compaty-address"])

			f.SetCellStyle("Sheet1", axis["compaty-inn"], axis["compaty-inn"], numStyle)
			claim.Company.INN = f.GetCellValue(sheetName, axis["compaty-inn"])

			fio := strings.Split(f.GetCellValue(sheetName, axis["compaty-fio"]), " ")
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

			claim.Company.Head.Contact.Phone = f.GetCellValue(sheetName, axis["compaty-phone"])
			claim.Company.Head.Contact.EMail = f.GetCellValue(sheetName, axis["compaty-email"])

			claim.Source = f.GetCellValue(sheetName, axis["compaty-cars"])
			claim.Cars = parser.ParseCars(claim.Source)

			claim.Agreement = f.GetCellValue(sheetName, axis["agreement"])
			claim.Reliability = f.GetCellValue(sheetName, axis["reliability"])
			claim.Ogrn = f.GetCellValue(sheetName, axis["ogrn"])

			out <- claim
		}
	}

	out <- nil

	return nil
}
