package command

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
	"github.com/alexey-zayats/claim-parser/internal/fsdump"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// FSdumpParser структура данных команды
type FSdumpParser struct {
	config *config.Config
	svc    *services.FSDumpService
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

	iface, err := backend.Parse(ctx, params)
	if err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	for _, claim := range iface.([]*model.Claim) {

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:
			if err := cmd.svc.SaveClaim(claim); err != nil {
				return errors.Wrapf(err, "unable save claim %#v", claim)
			}
		}

	}

	return nil
}
