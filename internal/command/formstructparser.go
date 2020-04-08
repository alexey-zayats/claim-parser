package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/formstruct"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// FormstructParser структура данных команды
type FormstructParser struct {
	config *config.Config
}

// FormstructParserParams - DI параметры команды
type FormstructParserParams struct {
	dig.In
	Config *config.Config
}

func init() {
	formstruct.Register()
}

// NewFormstructParser - конструктор команды
func NewFormstructParser(params FormstructParserParams) Command {
	return &FormstructParser{
		config: params.Config,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *FormstructParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(formstruct.Name)
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
