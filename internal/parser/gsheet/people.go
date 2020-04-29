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

// PeopleParser ...
type PeopleParser struct {
	name  string
	event *model.Event
	path  string
}

// NewPeopleParser ...
func NewPeopleParser(name string) (parser.Backend, error) {
	return &PeopleParser{
		name: name,
	}, nil
}

// WithEvent ...
func (p *PeopleParser) WithEvent(event *model.Event) {
	p.event = event
	p.path = event.Filepath
}

// Parse ...
func (p *PeopleParser) Parse(ctx context.Context, out chan *model.Out) error {

	logrus.WithFields(logrus.Fields{"name": p.name, "path": p.path, "event": p.event}).Debug("gsheet.People")

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

		//fmt.Printf("A%d\n", i)

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
				"employees":       fmt.Sprintf("I%d", i),
				"agreement":       fmt.Sprintf("J%d", i),
				"reliability":     fmt.Sprintf("K%d", i),
				"ogrn":            fmt.Sprintf("L%d", i),
			}

			f.SetCellStyle(sheetName, axis["created"], axis["created"], dateStyle)

			claim := &model.PeopleClaim{
				Created: parser.ExcelDateToDate(f.GetCellValue(sheetName, axis["created"])),
				Success: true,
			}

			claim.Company.Activity = f.GetCellValue(sheetName, axis["activity"])
			claim.Company.Title = f.GetCellValue(sheetName, axis["company-title"])
			claim.Company.Address = f.GetCellValue(sheetName, axis["company-address"])

			f.SetCellStyle(sheetName, axis["company-inn"], axis["company-inn"], numStyle)
			f.SetCellStyle(sheetName, axis["ogrn"], axis["ogrn"], numStyle)

			claim.Company.HeadName = strings.TrimSpace(f.GetCellValue(sheetName, axis["company-fio"]))
			claim.Company.HeadPhone = util.TrimSpaces(f.GetCellValue(sheetName, axis["company-phone"]))
			claim.Company.HeadEmail = util.TrimSpaces(f.GetCellValue(sheetName, axis["company-email"]))

			claim.Source = f.GetCellValue(sheetName, axis["employees"])
			if claim.Passes, ok = parser.ParsePeoples(claim.Source); ok == false {
				claim.Reason = append(claim.Reason, "не удалось разобрать список фио")
				claim.Success = false
			}

			claim.Agreement = util.TrimSpaces(f.GetCellValue(sheetName, axis["agreement"]))
			claim.Reliability = util.TrimSpaces(f.GetCellValue(sheetName, axis["reliability"]))

			claim.Company.INN = util.TrimSpaces(f.GetCellValue(sheetName, axis["company-inn"]))
			d1 := len(claim.Company.INN)
			if d1 < 10 || d1 > 12 {
				claim.Reason = append(claim.Reason, fmt.Sprintf("ИНН содержит %d цифр", d1))
				claim.Success = false
			}

			claim.Company.OGRN = util.TrimSpaces(f.GetCellValue(sheetName, axis["ogrn"]))
			d2 := len(claim.Company.OGRN)
			if d2 < 13 || d2 > 15 {
				claim.Reason = append(claim.Reason, fmt.Sprintf("ОРГН содержит %d цифр", d2))
				claim.Success = false
			}

			out <- &model.Out{
				Kind:  model.OutPeopleClaim,
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
