package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/fsdump"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"sync"
)

// FSdumpParser структура данных команды
type FSdumpParser struct {
	config *config.Config
	svc    *services.FSDumpService
	wg     sync.WaitGroup
}

// FSdumpParserInput - DI параметры команды
type FSdumpParserInput struct {
	dig.In
	Config *config.Config
	Svc    *services.FSDumpService
}

func init() {
	fsdump.Register()
}

// NewFSdumpParser - конструктор команды
func NewFSdumpParser(params FSdumpParserInput) Command {
	return &FSdumpParser{
		config: params.Config,
		svc:    params.Svc,
		wg:     sync.WaitGroup{},
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *FSdumpParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(fsdump.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)

	out := make(chan interface{})

	params.Set("out", out)

	cmd.wg.Add(1)
	go cmd.svc.HandleParsed(ctx, cmd.wg, out)

	if err := backend.Parse(ctx, params, out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	return nil
}