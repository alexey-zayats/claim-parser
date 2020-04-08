package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/godoc"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// GodocParser структура данных команды
type GodocParser struct {
	config *config.Config
	svc    *services.GodocService
	wg     sync.WaitGroup
	out    chan interface{}
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
		out:    make(chan interface{}),
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

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, params, cmd.out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	return nil
}

// HandleParsed ...
func (cmd *GodocParser) HandleParsed(ctx context.Context) {
	defer cmd.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case iface := <-cmd.out:

			switch iface.(type) {
			case nil:
				return
			}

			claim := iface.(*model.Claim)

			logrus.WithFields(logrus.Fields{"company": claim.Company.Title}).Debug("claim")

			if err := cmd.svc.SaveClaim(ctx, claim); err != nil {
				logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save claim")
			}
		}
	}
}
