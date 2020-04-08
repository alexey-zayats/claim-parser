package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/excel"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// XlsxParser структура данных команды
type XlsxParser struct {
	config *config.Config
}

// XlsxParserParams - DI параметры команды
type XlsxParserParams struct {
	dig.In
	Config *config.Config
}

func init() {
	excel.Register()
}

// NewXlsxParser - конструктор команды
func NewXlsxParser(params XlsxParserParams) Command {
	return &XlsxParser{
		config: params.Config,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *XlsxParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(excel.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)

	company, err := backend.Parse(ctx, params)
	if err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	data, err := json.MarshalIndent(company, "", "\t")
	if err != nil {
		return errors.Wrap(err, "unable unmarshal ")
	}

	fmt.Printf("%s\n", string(data))

	return nil
}
