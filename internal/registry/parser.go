package registry

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// Parser ...
type Parser struct {
}

// Name ...
const Name = "registry"

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

	logrus.WithFields(logrus.Fields{"name": Name, "path": path}).Debug("registry.Parse")

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
		return errors.Wrapf(err, "unable get rows for sheet %s", sheetName)
	}

	re := regexp.MustCompile(`\s`)

	for rows.Next() {
		for i, col := range rows.Columns() {
			if i != 4 {
				continue
			}

			col = re.ReplaceAllString(col, "")
			col = strings.ToUpper(col)

			if len(col) < 6 || len(col) > 12 {
				continue
			}

			out <- col
		}
	}

	return nil
}
