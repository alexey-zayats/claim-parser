package excel

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// Parser ...
type Parser struct {
	reNumber *regexp.Regexp
}

// Name ...
const Name = "excel"

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

	logrus.WithFields(logrus.Fields{"name": Name, "path": path}).Debug("excel.Parse")

	f, err := excelize.OpenFile(path)
	if err != nil {
		return errors.Wrapf(err, "unable open xlsx file %s", path)
	}

	var sheetName string
	for _, sheet := range f.GetSheetMap() {
		sheetName = sheet
		break
	}

	kind := f.GetCellValue(sheetName, "B5")
	address := f.GetCellValue(sheetName, "B7")

	source := f.GetCellValue(sheetName, "B12")

	claim := &model.Claim{
		District: f.GetCellValue(sheetName, "A1"),
		Company: model.Company{
			Activity: strings.ReplaceAll(kind, "\n", ", "),
			Title:    f.GetCellValue(sheetName, "B6"),
			Address:  strings.ReplaceAll(address, "\n", ", "),
			INN:      strings.ReplaceAll(f.GetCellValue(sheetName, "B8"), " ", ""),
		},
		Cars:        parser.ParseCars(source),
		Agreement:   f.GetCellValue(sheetName, "B13"),
		Reliability: f.GetCellValue(sheetName, "B14"),
		Reason:      nil,
		Valid:       true,
		Source:      source,
	}

	claim.Company.Head.Contact = model.Contact{
		Phone: f.GetCellValue(sheetName, "B10"),
		EMail: f.GetCellValue(sheetName, "B11"),
	}

	fio := strings.Split(f.GetCellValue(sheetName, "B9"), " ")

	if len(fio) < 3 {
		claim.Valid = false
		reason := "Нет данных по ФИО руководителя"
		claim.Reason = &reason
	} else {
		claim.Company.Head.FIO = model.FIO{
			Surname:    fio[0],
			Name:       fio[1],
			Patronymic: fio[2],
		}
	}

	out <- claim

	return nil
}
