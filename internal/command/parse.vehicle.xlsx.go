package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/excel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"strings"
	"sync"
)

// VehicleXlsxParser структура данных команды
type VehicleXlsxParser struct {
	config *config.Config
	wg     sync.WaitGroup
	out    chan *model.Out
	parser string
}

// VehicleXlsxParserDI - DI параметры команды
type VehicleXlsxParserDI struct {
	dig.In
	Config *config.Config
}

func init() {
	excel.Register()
}

// NewVehicdleXlsxParser - конструктор команды
func NewVehicdleXlsxParser(params VehicleXlsxParserDI) Command {

	cmd := &VehicleXlsxParser{
		config: params.Config,
		wg:     sync.WaitGroup{},
		out:    make(chan *model.Out, 1),
		parser: "excel",
	}

	return cmd
}

// Run - имплементация метода Run интерфейса Command
func (cmd *VehicleXlsxParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(cmd.parser)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	backend.WithEvent(&model.Event{
		Filepath:   cmd.config.Parser.Path,
		CreatedBy:  1,
		PassType:   1,
		DistrictID: 1,
	})

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, cmd.out); err != nil {
		return errors.Wrap(err, "unable call backend.Parse ")
	}

	cmd.wg.Wait()

	return nil
}

// HandleParsed ...
func (cmd *VehicleXlsxParser) HandleParsed(ctx context.Context) {

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
					fmt.Printf("%s;%s;%s;parse: %s\n", claim.Created, claim.Company.INN, claim.Company.Title, strings.Join(claim.Reason, ", "))
				}
				return
			default:
				return
			}
		}
	}
}
