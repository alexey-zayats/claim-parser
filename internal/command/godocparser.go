package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/godoc"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"sync"
)

// GodocParser структура данных команды
type GodocParser struct {
	config *config.Config
	svc    *services.GodocService
	wg     sync.WaitGroup
}

// GodocParserDI - DI параметры команды
type GodocParserDI struct {
	dig.In
	Config *config.Config
	Svc    *services.GodocService
}

func init() {
	godoc.Register()
}

// NewGodocParser - конструктор команды
func NewGodocParser(di GodocParserDI) Command {
	return &GodocParser{
		config: di.Config,
		svc:    di.Svc,
		wg:     sync.WaitGroup{},
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *GodocParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(godoc.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)

	out := make(chan interface{})

	cmd.wg.Add(1)
	go cmd.svc.HandleParsed(ctx, cmd.wg, out)

	if err := backend.Parse(ctx, params, out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	return nil
}
