package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/fs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"strings"
	"sync"
)

// VehicleFSParser структура данных команды
type VehicleFSParser struct {
	config *config.Config
	wg     sync.WaitGroup
	out    chan *model.Out
	parser string
}

// VehicleFSParserDI - DI параметры команды
type VehicleFSParserDI struct {
	dig.In
	Config *config.Config
}

func init() {
	fs.Register()
}

// NewVehicleFSParser - конструктор команды
func NewVehicleFSParser(di VehicleFSParserDI) Command {
	return &VehicleFSParser{
		config: di.Config,
		wg:     sync.WaitGroup{},
		out:    make(chan *model.Out, 1),
		parser: "form.struct",
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *VehicleFSParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(cmd.parser)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	backend.WithEvent(&model.Event{
		Filepath:  cmd.config.Parser.Path,
		CreatedBy: 1,
		PassType:  1,
	})

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, cmd.out); err != nil {
		return errors.Wrap(err, "unable call parser.Parse ")
	}

	cmd.wg.Wait()

	return nil
}

// HandleParsed ...
func (cmd *VehicleFSParser) HandleParsed(ctx context.Context) {

	defer cmd.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-cmd.out:

			switch out.Kind {
			case model.OutVehicleClaim:
				claim := out.Value.(*model.VehicleClaim)
				data, err := json.MarshalIndent(claim, "", "\t")
				if err != nil {
					logrus.WithFields(logrus.Fields{"reason": err, "claim": claim}).Error("unable marshal")
				}
				fmt.Printf("%s\n", string(data))
				if claim.Success == false {
					fmt.Printf("%s;%d;%s;parse: %s\n", claim.Created, claim.Company.TIN, claim.Company.Title, strings.Join(claim.Reason, ", "))
				}
				return
			default:
				return
			}

		}
	}
}
