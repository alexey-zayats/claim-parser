package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/registry"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"sync"
)

// RegistryParser структура данных команды
type RegistryParser struct {
	config *config.Config
	svc    *services.RegistryService
	wg     sync.WaitGroup
}

// RegistryParserDI - DI параметры команды
type RegistryParserDI struct {
	dig.In
	Config *config.Config
	Svc    *services.RegistryService
}

func init() {
	registry.Register()
}

// NewRegistryParser - конструктор команды
func NewRegistryParser(di RegistryParserDI) Command {
	return &RegistryParser{
		config: di.Config,
		svc:    di.Svc,
		wg:     sync.WaitGroup{},
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

	out := make(chan interface{})

	cmd.wg.Add(1)
	go cmd.svc.HandleParsed(ctx, cmd.wg, out)

	if err := backend.Parse(ctx, params, out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	return nil
}
