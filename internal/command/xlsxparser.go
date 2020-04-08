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
	return &XlsxParser{
		config: params.Config,
		wg:     sync.WaitGroup{},
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *XlsxParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(excel.Name)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	params := dict.New()
	params.Set("path", cmd.config.Parser.Path)

	out := make(chan interface{})

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx, out)

	if err := backend.Parse(ctx, params, out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	cmd.wg.Wait()

	return nil
}

// HandleParsed ...
func (cmd *XlsxParser) HandleParsed(ctx context.Context, out chan interface{}) {

	defer cmd.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case iface := <-out:
			claim := iface.(*model.Claim)
			data, err := json.MarshalIndent(claim, "", "\t")
			if err != nil {
				logrus.WithFields(logrus.Fields{"reason": err, "claim": claim}).Error("unable marshal")
			}

			fmt.Printf("%s\n", string(data))
			return
		}
	}
}
