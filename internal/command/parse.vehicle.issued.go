package command

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/issued"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"strings"
	"sync"
	"time"
)

// VehicleIssuedParser структура данных команды
type VehicleIssuedParser struct {
	config    *config.Config
	issuedSvc *services.VehicleIssuedService
	fileSvc   *services.FileService
	wg        sync.WaitGroup
	out       chan *model.Out
	parser    string
	file      *entity.File
}

// VehicleIssuedParserDI - DI параметры команды
type VehicleIssuedParserDI struct {
	dig.In
	Config    *config.Config
	IssuedSvc *services.VehicleIssuedService
	FileSvc   *services.FileService
}

func init() {
	issued.Register()
}

// NewVehicleIssuedParser - конструктор команды
func NewVehicleIssuedParser(di VehicleIssuedParserDI) Command {
	return &VehicleIssuedParser{
		config:    di.Config,
		issuedSvc: di.IssuedSvc,
		fileSvc:   di.FileSvc,
		wg:        sync.WaitGroup{},
		out:       make(chan *model.Out),
		parser:    "issued",
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *VehicleIssuedParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(cmd.parser)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	cmd.file = &entity.File{
		Filepath:  cmd.config.Parser.Path,
		Status:    0,
		Log:       "",
		Source:    "issued",
		CreatedAt: time.Now(),
	}

	if err := cmd.fileSvc.Create(cmd.file); err != nil {
		return errors.Wrap(err, "unable create file record")
	}

	backend.WithEvent(&model.Event{
		FileID:     cmd.file.ID,
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

	return nil
}

// HandleParsed ...
func (cmd *VehicleIssuedParser) HandleParsed(ctx context.Context) {
	defer cmd.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-cmd.out:

			switch out.Kind {
			case model.OutVehicleRegistry:

				registry := out.Value.(*model.VehicleRegistry)

				rec := fmt.Sprintf("%s;%d;%s", registry.IssuedAt, registry.CompanyInn, registry.CompanyName)

				if registry.Success {

					logrus.WithFields(logrus.Fields{
						"company": registry.CompanyName,
						"car":     registry.CompanyCar,
						"fio":     registry.CompanyFio,
					}).Debug("Vehicle.Issued")

					if err := cmd.issuedSvc.SaveRecord(out.Event, registry); err != nil {

						logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save claim")

						cmd.file.Log = rec + ";sql: " + err.Error() + "\n"
						cmd.file.Status = 3

						if err := cmd.fileSvc.UpdateState(cmd.file); err != nil {
							logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update file state")
						}
					}
				} else {

					cmd.file.Log = rec + ";parse: " + strings.Join(registry.Reason, ", ") + "\n"
					cmd.file.Status = 2

					if err := cmd.fileSvc.UpdateState(cmd.file); err != nil {
						logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update file state")
					}
				}

			default:
				return
			}

		}
	}
}
