package command

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/parser/gsheet"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
	"time"
)

// VehicleGSheetParser структура данных команды
type VehicleGSheetParser struct {
	config   *config.Config
	claimSvc *services.VehicleClaimService
	fileSvc  *services.FileService
	wg       sync.WaitGroup
	out      chan *model.Out
	parser   string
	file     *entity.File
}

// VehicleGSheetDI - DI параметры команды
type VehicleGSheetDI struct {
	dig.In
	Config   *config.Config
	ClaimSvc *services.VehicleClaimService
	FileSvc  *services.FileService
}

func init() {
	gsheet.Register()
}

// NewVehicleGSheetParser - конструктор команды
func NewVehicleGSheetParser(di VehicleGSheetDI) Command {
	return &VehicleGSheetParser{
		config:   di.Config,
		claimSvc: di.ClaimSvc,
		fileSvc:  di.FileSvc,
		wg:       sync.WaitGroup{},
		out:      make(chan *model.Out),
		parser:   "gsheet.vehicle",
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *VehicleGSheetParser) Run(ctx context.Context, args []string) error {

	backend, err := parser.Instance().Backend(cmd.parser)
	if err != nil {
		return errors.Wrap(err, "unable find parser for")
	}

	cmd.file = &entity.File{
		Filepath:  cmd.config.Parser.Path,
		Status:    0,
		Log:       "",
		Source:    "gsheet.vehicle",
		CreatedAt: time.Now(),
	}

	if err := cmd.fileSvc.Create(cmd.file); err != nil {
		return errors.Wrap(err, "unable create file record")
	}

	backend.WithEvent(&model.Event{
		FileID:    cmd.file.ID,
		Filepath:  cmd.config.Parser.Path,
		CreatedBy: 1,
		PassType:  1,
	})

	cmd.wg.Add(1)
	go cmd.HandleParsed(ctx)

	if err := backend.Parse(ctx, cmd.out); err != nil {
		return errors.Wrap(err, "unable call backend.Parse ")
	}

	return nil
}

// HandleParsed ...
func (cmd *VehicleGSheetParser) HandleParsed(ctx context.Context) {
	defer cmd.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-cmd.out:

			switch out.Kind {
			case model.OutVehicleClaim:

				claim := out.Value.(*model.VehicleClaim)
				rec := fmt.Sprintf("%s;%d;%s", claim.Created, claim.Company.TIN, claim.Company.Title)

				logrus.WithFields(logrus.Fields{"company": claim.Company.Title}).Debug("claim")

				if err := cmd.claimSvc.SaveRecord(out.Event, claim); err != nil {

					logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save claim")

					cmd.file.Log = rec + ";sql: " + err.Error() + "\n"
					cmd.file.Status = 3

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
