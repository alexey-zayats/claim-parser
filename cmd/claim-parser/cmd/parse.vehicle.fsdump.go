package cmd

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/command"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/di"
	"github.com/alexey-zayats/claim-parser/internal/repository"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fsdumpCmd = &cobra.Command{
	Use:   "fsdump",
	Short: "fsdump",
	Long:  "fsdump",
	Run:   fsdumpMain,
}

func init() {
	vehicleCmd.AddCommand(fsdumpCmd)
}

func fsdumpMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                         config.NewConfig,
			"database.Connection":            database.NewConnection,
			"repository.VehiclePassRepo":     repository.NewVehiclePassRepo,
			"repository.VehicleBidRepo":      repository.NewVehicleBidRepo,
			"repository.VehicleIssuedRepo":   repository.NewVehicleIssuedRepo,
			"repository.VehicleCompanyRepo":  repository.NewVehicleCompanyRepo,
			"repository.BranchRepository":    repository.NewBranchRepository,
			"repository.FileRepository":      repository.NewFileRepository,
			"service.VehiclePassService":     services.NewVehiclePassService,
			"service.VehicleBidService":      services.NewVehicleBidService,
			"services.VehicleClaimService":   services.NewVehicleClaimService,
			"services.VehicleIssuedService":  services.NewVehicleIssuedService,
			"services.VehicleCompanyService": services.NewVehicleCompanyService,
			"services.FileService":           services.NewFileService,
			"services.BranchService":         services.NewBranchService,
			"command.Parser":                 command.NewVehicleFSDumpParser,
		},
		Invoke: func(ctx context.Context, args []string) interface{} {
			return func(i command.Command) {
				if err := i.Run(ctx, args); err != nil {
					logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable run command")
				}
			}
		},
	}

	di.Run(ctx, args)
}
