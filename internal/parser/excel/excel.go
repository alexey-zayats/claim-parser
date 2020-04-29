package excel

import (
	"context"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

// Parser ...
type Parser struct {
	path  string
	event *model.Event
	name  string
}

// Register ...
func Register() {
	parser.Instance().Add("excel", NewParser)
}

// NewParser ...
func NewParser(name string) (parser.Backend, error) {
	return &Parser{
		name: name,
	}, nil
}

// WithEvent ...
func (p *Parser) WithEvent(event *model.Event) {
	p.event = event
	p.path = event.Filepath
}

// Parse ...
func (p *Parser) Parse(ctx context.Context, out chan *model.Out) error {

	logrus.WithFields(logrus.Fields{"name": p.name, "path": p.path, "event": p.event}).Debug("excel.Parse")

	f, err := excelize.OpenFile(p.path)
	if err != nil {
		return errors.Wrapf(err, "unable open xlsx file %s", p.path)
	}

	var sheetName string
	for _, sheet := range f.GetSheetMap() {
		sheetName = sheet
		break
	}

	ok := false

	kind := f.GetCellValue(sheetName, "B5")
	address := f.GetCellValue(sheetName, "B7")

	source := f.GetCellValue(sheetName, "B12")

	claim := &model.VehicleClaim{
		District: f.GetCellValue(sheetName, "A1"),

		Company: model.Company{
			Activity: strings.ReplaceAll(kind, "\n", ", "),
			Title:    f.GetCellValue(sheetName, "B6"),
			Address:  strings.ReplaceAll(address, "\n", ", "),
		},
		Agreement:   f.GetCellValue(sheetName, "B13"),
		Reliability: f.GetCellValue(sheetName, "B14"),
		Source:      source,
		Success:     true,
	}

	claim.Company.INN = util.TrimSpaces(f.GetCellValue(sheetName, "B8"))
	tdig := len(claim.Company.INN)
	if tdig < 10 || tdig > 12 {
		claim.Reason = append(claim.Reason, "ИНН меньше 10 или больше 12 знаков")
		claim.Success = false
	}

	claim.Company.HeadName = strings.TrimSpace(f.GetCellValue(sheetName, "B9"))
	claim.Company.HeadPhone = util.TrimSpaces(f.GetCellValue(sheetName, "B10"))
	claim.Company.HeadEmail = util.TrimSpaces(f.GetCellValue(sheetName, "B11"))

	if claim.Passes, ok = parser.ParseVehicles(source); ok == false {
		claim.Reason = append(claim.Reason, "не удалось разобрать список номер/фио")
		claim.Success = false
	}

	out <- &model.Out{
		Kind:  model.OutVehicleClaim,
		Event: p.event,
		Value: claim,
	}

	out <- &model.Out{
		Kind: model.OutUnknown,
	}

	return nil
}
