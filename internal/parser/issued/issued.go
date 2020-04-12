package issued

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Parser ...
type Parser struct {
}

// Name ...
const Name = "issued"

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

	i := 2
	rows.Next()
	rows.Next()

	numStyle, _ := f.NewStyle(`{"number_format":1}`)
	dateStyle, _ := f.NewStyle(`{"custom_number_format": "m.d.yyyy h:mm:ss"}`)

	for rows.Next() {

		i++

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			axis := map[string]string{
				"inn":             fmt.Sprintf("A%d", i),
				"ogrn":            fmt.Sprintf("B%d", i),
				"name":            fmt.Sprintf("C%d", i),
				"fio":             fmt.Sprintf("D%d", i),
				"car":             fmt.Sprintf("E%d", i),
				"basement-pass":   fmt.Sprintf("F%d", i),
				"district":        fmt.Sprintf("G%d", i),
				"pass-type":       fmt.Sprintf("H%d", i),
				"issued-at":       fmt.Sprintf("I%d", i),
				"registry-number": fmt.Sprintf("J%d", i),
				"shipping":        fmt.Sprintf("K%d", i),
			}

			f.SetCellStyle("Sheet1", axis["issued-at"], axis["issued-at"], dateStyle)
			f.SetCellStyle("Sheet1", axis["inn"], axis["inn"], numStyle)
			f.SetCellStyle("Sheet1", axis["ogrn"], axis["ogrn"], numStyle)
			f.SetCellStyle("Sheet1", axis["ogrn"], axis["ogrn"], numStyle)

			var legalBasement string
			var passNumber string

			basementPass := f.GetCellValue(sheetName, axis["basement-pass"])

			semicolon := strings.LastIndex(basementPass, ";")
			colonIndex := strings.LastIndex(basementPass, ",")

			if semicolon == -1 && colonIndex == -1 {
				passNumber = basementPass
			} else {
				splitIndex := -1
				if semicolon > colonIndex {
					splitIndex = semicolon
				} else {
					splitIndex = colonIndex
				}
				legalBasement = basementPass[0:splitIndex]
				passNumber = basementPass[splitIndex+1:]
			}

			passNumber = strings.ReplaceAll(passNumber, "№", "")

			legalBasement = util.TrimSpaces(passNumber)
			passNumber = util.TrimSpaces(passNumber)

			passType := 0
			passTypeStr := f.GetCellValue(sheetName, axis["pass-type"])
			if strings.Compare(passTypeStr, "Краснодар") == 0 {
				passType = 1
			} else if strings.Compare(passTypeStr, "Краснодарский край") == 0 {
				passType = 2
			}

			shipping := 0
			shippingStr := f.GetCellValue(sheetName, axis["shipping"])
			if strings.Compare(shippingStr, "электронно") == 0 {
				shipping = 1
			} else if strings.Compare(shippingStr, "нарочно") == 0 {
				shipping = 2
			}

			var issued *time.Time
			issuedStr := f.GetCellValue(sheetName, axis["issued-at"])
			if len(issuedStr) > 0 {
				t := parser.ExcelDateToDate(issuedStr)
				issued = &t
			}

			var car string
			carCell := f.GetCellValue(sheetName, axis["car"])
			carCell = strings.ToUpper(strings.ReplaceAll(carCell, " ", ""))
			cars := strings.Split(carCell, ",")
			if len(cars) > 0 {
				car = cars[0]
			} else {
				car = carCell
			}

			record := &model.Registry{
				CompanyInn:     f.GetCellValue(sheetName, axis["inn"]),
				CompanyOgrn:    f.GetCellValue(sheetName, axis["ogrn"]),
				CompanyName:    f.GetCellValue(sheetName, axis["name"]),
				CompanyFio:     f.GetCellValue(sheetName, axis["fio"]),
				CompanyCar:     parser.NormalizeCarNumber(car),
				LegalBasement:  legalBasement,
				PassNumber:     passNumber,
				District:       f.GetCellValue(sheetName, axis["district"]),
				PassType:       passType,
				IssuedAt:       issued,
				RegistryNumber: f.GetCellValue(sheetName, axis["registry-number"]),
				Shipping:       shipping,
				Event:          event,
			}

			out <- record
		}
	}

	out <- nil

	return nil
}
