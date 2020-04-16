package gsheet

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

// VehicleParser ...
type VehicleParser struct {
	name  string
	event *model.Event
	path  string
}

// NewVehicleParser ...
func NewVehicleParser(name string) (parser.Backend, error) {
	return &VehicleParser{
		name: name,
	}, nil
}

// WithEvent ...
func (p *VehicleParser) WithEvent(event *model.Event) {
	p.event = event
	p.path = event.Filepath
}

// Parse ...
func (p *VehicleParser) Parse(ctx context.Context, out chan *model.Out) error {

	logrus.WithFields(logrus.Fields{"name": p.name, "path": p.path, "event": p.event}).Debug("gsheet.Vehicle")

	f, err := excelize.OpenFile(p.path)
	if err != nil {
		return errors.Wrapf(err, "unable open xlsx file %s", p.path)
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

		test := f.GetCellValue(sheetName, fmt.Sprintf("A%d", i))
		if len(test) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			var ok bool

			axis := map[string]string{
				"created":         fmt.Sprintf("A%d", i),
				"activity":        fmt.Sprintf("B%d", i),
				"company-title":   fmt.Sprintf("C%d", i),
				"company-address": fmt.Sprintf("D%d", i),
				"company-inn":     fmt.Sprintf("E%d", i),
				"company-fio":     fmt.Sprintf("F%d", i),
				"company-phone":   fmt.Sprintf("G%d", i),
				"company-email":   fmt.Sprintf("H%d", i),
				"company-cars":    fmt.Sprintf("I%d", i),
				"agreement":       fmt.Sprintf("J%d", i),
				"reliability":     fmt.Sprintf("K%d", i),
				"ogrn":            fmt.Sprintf("L%d", i),
			}

			f.SetCellStyle(sheetName, axis["created"], axis["created"], dateStyle)
			f.SetCellStyle(sheetName, axis["company-inn"], axis["company-inn"], numStyle)
			f.SetCellStyle(sheetName, axis["ogrn"], axis["ogrn"], numStyle)

			claim := &model.VehicleClaim{
				Created: parser.ExcelDateToDate(f.GetCellValue(sheetName, axis["created"])),
				Success: true,
			}

			claim.Company.Activity = f.GetCellValue(sheetName, axis["activity"])
			claim.Company.Title = f.GetCellValue(sheetName, axis["company-title"])
			claim.Company.Address = f.GetCellValue(sheetName, axis["company-address"])

			claim.Company.HeadName = strings.TrimSpace(f.GetCellValue(sheetName, axis["company-fio"]))

			claim.Company.HeadPhone = f.GetCellValue(sheetName, axis["company-phone"])
			claim.Company.HeadEmail = f.GetCellValue(sheetName, axis["company-email"])

			claim.Source = f.GetCellValue(sheetName, axis["company-cars"])
			if claim.Passes, ok = parser.ParseVehicles(claim.Source); ok == false {
				claim.Reason = append(claim.Reason, "не удалось разобрать список номер/фио")
				claim.Success = false
			}

			claim.Agreement = f.GetCellValue(sheetName, axis["agreement"])
			claim.Reliability = f.GetCellValue(sheetName, axis["reliability"])

			if claim.Company.TIN, ok = parser.ParseInt64(f.GetCellValue(sheetName, axis["company-inn"])); ok == false {
				claim.Reason = append(claim.Reason, "ИНН не является числом")
				claim.Success = false
			}

			d := util.DigitsCount(claim.Company.TIN)
			if d < 10 || d > 12 {
				claim.Reason = append(claim.Reason, fmt.Sprintf("ИНН содержит %d цифр", d))
				claim.Success = false
			}

			if claim.Company.PSRN, ok = parser.ParseInt64(f.GetCellValue(sheetName, axis["ogrn"])); ok == false {
				claim.Reason = append(claim.Reason, "поле ОРГ не является числом")
				claim.Success = false
			}

			d = util.DigitsCount(claim.Company.PSRN)
			if d < 13 || d > 15 {
				claim.Reason = append(claim.Reason, fmt.Sprintf("ОРГН содержит %d цифр", d))
				claim.Success = false
			}

			out <- &model.Out{
				Kind:  model.OutVehicleClaim,
				Event: p.event,
				Value: claim,
			}
		}
	}

	out <- &model.Out{
		Kind: model.OutUnknown,
	}

	return nil
}
