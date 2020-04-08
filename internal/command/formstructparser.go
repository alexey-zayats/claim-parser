package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/formstruct"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// FormstructParser структура данных команды
type FormstructParser struct {
	config *config.Config
	wg     sync.WaitGroup
	out    chan interface{}
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
		wg:     sync.WaitGroup{},
		out:    make(chan interface{}),
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

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, params, cmd.out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	cmd.wg.Wait()

	return nil
}

// HandleParsed ...
func (cmd *FormstructParser) HandleParsed(ctx context.Context) {

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
			data, err := json.MarshalIndent(claim, "", "\t")
			if err != nil {
				logrus.WithFields(logrus.Fields{"reason": err, "claim": claim}).Error("unable marshal")
			}

			fmt.Printf("%s\n", string(data))
		}
	}
}
