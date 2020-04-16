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

var vehicleIssuedCmd = &cobra.Command{
	Use:   "issued",
	Short: "issued",
	Long:  "issued",
	Run:   vehicleIssuedMain,
}

func init() {
	vehicleCmd.AddCommand(vehicleIssuedCmd)
}

func vehicleIssuedMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                           config.NewConfig,
			"database.NewConnection":           database.NewConnection,
			"repository.NewVehiclePassRepo":    repository.NewVehiclePassRepo,
			"repository.NewVehicleBidRepo":     repository.NewVehicleBidRepo,
			"repository.NewVehicleIssuedRepo":  repository.NewVehicleIssuedRepo,
			"repository.FileRepository":        repository.NewFileRepository,
			"service.NewVehiclePassService":    services.NewVehiclePassService,
			"service.NewVehicleBidService":     services.NewVehicleBidService,
			"services.NewVehicleIssuedService": services.NewVehicleIssuedService,
			"services.FileService":             services.NewFileService,
			"command.NewVehicleIssuedParser":   command.NewVehicleIssuedParser,
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
