package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/registry"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// RegistryParser структура данных команды
type RegistryParser struct {
	config *config.Config
}

// RegistryParserDI - DI параметры команды
type RegistryParserDI struct {
	dig.In
	Config *config.Config
}

func init() {
	registry.Register()
}

// NewRegistryParser - конструктор команды
func NewRegistryParser(di RegistryParserDI) Command {
	return &RegistryParser{
		config: di.Config,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *RegistryParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(registry.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)
	params.Set("sheet", "Реестр красный")

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
