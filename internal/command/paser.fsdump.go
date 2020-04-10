package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/fsdump"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// FsdumpParser структура данных команды
type FsdumpParser struct {
	config *config.Config
	svc    *services.ClaimService
	wg     sync.WaitGroup
	out    chan interface{}
}

// FsdumpParserDI - DI параметры команды
type FsdumpParserDI struct {
	dig.In
	Config *config.Config
	Svc    *services.ClaimService
}

func init() {
	fsdump.Register()
}

// NewFsdumpParser - конструктор команды
func NewFsdumpParser(di FsdumpParserDI) Command {
	return &FsdumpParser{
		config: di.Config,
		svc:    di.Svc,
		wg:     sync.WaitGroup{},
		out:    make(chan interface{}, 1),
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *FsdumpParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(fsdump.Name)
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

	cmd.wg.Wait()

	return nil
}

// HandleParsed ...
func (cmd *FsdumpParser) HandleParsed(ctx context.Context) {
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

			if err := cmd.svc.SaveRecord(ctx, claim); err != nil {
				logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save claim")
			}
		}
	}

}
