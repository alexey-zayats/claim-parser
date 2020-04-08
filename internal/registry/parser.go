package registry

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
func (p *Parser) Parse(ctx context.Context, param *dict.Dict) (interface{}, error) {

	var path string
	var sheet string

	if iface, ok := param.Get("path"); ok {
		path = iface.(string)
	} else {
		return nil, fmt.Errorf("not found 'path' in param dict")
	}

	if iface, ok := param.Get("sheet"); ok {
		sheet = iface.(string)
	} else {
		return nil, fmt.Errorf("not found 'sheet' in param dict")
	}

	logrus.WithFields(logrus.Fields{"name": Name, "path": path}).Debug("Parser.Parse")

	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable open xlsx file %s", path)
	}

	rows, err := f.Rows("Реестр красный")
	if err != nil {
		return nil, errors.Wrapf(err, "unable get rows for sheet %s", sheet)
	}
	for rows.Next() {
		for i, col := range rows.Columns() {
			if i != 4 {
				fmt.Printf("%s\n", col)
			}
		}
	}

	return nil, nil
}
