package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/excel"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// XlsxParser структура данных команды
type XlsxParser struct {
	config *config.Config
	wg     sync.WaitGroup
	out    chan interface{}
}

// XlsxParserParams - DI параметры команды
type XlsxParserParams struct {
	dig.In
	Config *config.Config
}

func init() {
	excel.Register()
}

// NewXlsxParser - конструктор команды
func NewXlsxParser(params XlsxParserParams) Command {

	cmd := &XlsxParser{
		config: params.Config,
		wg:     sync.WaitGroup{},
		out:    make(chan interface{}, 1),
	}

	return cmd
}

// Run - имплементация метода Run интерфейса Command
func (cmd *XlsxParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(excel.Name)
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
func (cmd *XlsxParser) HandleParsed(ctx context.Context) {

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
