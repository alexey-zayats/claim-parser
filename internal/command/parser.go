package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/formstruct"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/xlsx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// Parser структура данных команды
type Parser struct {
	config *config.Config
}

// ParserParams - DI параметры команды
type ParserParams struct {
	dig.In
	Config *config.Config
}

func init() {
	xlsx.Register()
	formstruct.Register()
}

// NewParser - конструктор команды
func NewParser(params ParserParams) Command {
	return &Parser{
		config: params.Config,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *Parser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(cmd.config.Parser.Source)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	company, err := backend.Parse(ctx, cmd.config.Parser.Path)
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
