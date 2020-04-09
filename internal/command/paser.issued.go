package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/issued"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// IssuedParser структура данных команды
type IssuedParser struct {
	config *config.Config
	svc    *services.IssuedService
	wg     sync.WaitGroup
	out    chan interface{}
}

// IssuedParserDI - DI параметры команды
type IssuedParserDI struct {
	dig.In
	Config *config.Config
	Svc    *services.IssuedService
}

func init() {
	issued.Register()
}

// NewIssuedParser - конструктор команды
func NewIssuedParser(di IssuedParserDI) Command {
	return &IssuedParser{
		config: di.Config,
		svc:    di.Svc,
		wg:     sync.WaitGroup{},
		out:    make(chan interface{}),
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *IssuedParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(issued.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)
	params.Set("sheet", "Реестр красный")

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, params, cmd.out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	return nil
}

// HandleParsed ...
func (cmd *IssuedParser) HandleParsed(ctx context.Context) {
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

			record := iface.(*model.Registry)

			logrus.WithFields(logrus.Fields{
				"company": record.CompanyName,
			}).Debug("Registry")

			if err := cmd.svc.SaveRecord(ctx, record); err != nil {
				logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save registry")
			}
		}
	}
}
