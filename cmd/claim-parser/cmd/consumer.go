package cmd

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/command"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/di"
	"github.com/alexey-zayats/claim-parser/internal/queue"
	"github.com/alexey-zayats/claim-parser/internal/repository"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consume",
	Short: "consume",
	Long:  "consume",
	Run:   consumerMain,
}

func init() {

	rootCmd.AddCommand(consumerCmd)

	cfgParams := []config.Param{
		{Name: "sql-dsn", Value: "pass:pass@tcp(127.0.0.1:3306)/pass", Usage: "sql driver", ViperBind: "Sql.Dsn"},
		{Name: "sql-conns-max-idle", Value: 0, Usage: "Maximum number of connections in the idle", ViperBind: "Sql.Conns.Max.Idle"},
		{Name: "sql-conns-max-open", Value: 2, Usage: "Maximum number of open connections to the database", ViperBind: "Sql.Conns.Max.Open"},
		{Name: "sql-conns-max-lifetime", Value: 10, Usage: "Maximum amount of time a connection may be reused", ViperBind: "Sql.Conns.Max.Open"},

		{Name: "amqp-dsn", Value: "amqp://pass:pass@127.0.0.1:5672/", Usage: "AMQP datasource", ViperBind: "Amqp.Dsn"},
		{Name: "amqp-exchange", Value: "collector", Usage: "AMQP Exchange name publish to", ViperBind: "Amqp.Exchange"},

		{Name: "amqp-workers", Value: 4, Usage: "form workers number", ViperBind: "Amqp.Workers"},

		{Name: "amqp-vehicle-routing", Value: "form.vehicle", Usage: "vehicle form routing key", ViperBind: "Amqp.Vehicle.Routing"},
		{Name: "amqp-vehicle-queue", Value: "form.vehicle", Usage: "vehicle form queue name", ViperBind: "Amqp.Vehicle.Queue"},

		{Name: "amqp-people-routing", Value: "form.people", Usage: "people form routing key", ViperBind: "Amqp.People.Routing"},
		{Name: "amqp-people-queue", Value: "form.people", Usage: "people form queue name", ViperBind: "Amqp.People.Queue"},

		{Name: "pass-source", Value: 0, Usage: "pass source", ViperBind: "Pass.Source"},
		{Name: "pass-creator", Value: 0, Usage: "pass creator", ViperBind: "Pass.Creator"},
		{Name: "pass-clean", Value: 0, Usage: "clean pass operator", ViperBind: "Pass.Clean"},
		{Name: "pass-dirty", Value: 0, Usage: "dirty pass operator", ViperBind: "Pass.Dirty"},

		{Name: "csv", Value: "claims.csv", Usage: "csv file", ViperBind: "CSV"},
	}

	config.Apply(consumerCmd, cfgParams)
}

func consumerMain(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                             config.NewConfig,
			"database.Connection":                database.NewConnection,
			"queue.Connection":                   queue.NewConnection,
			"queue.Queue":                        queue.NewQueue,
			"repository.FileRepository":          repository.NewFileRepository,
			"repository.VehiclePassRepo":         repository.NewVehiclePassRepo,
			"repository.VehicleBidRepo":          repository.NewVehicleBidRepo,
			"repository.VehicleIssuedRepo":       repository.NewVehicleIssuedRepo,
			"repository.PeoplePassRepo":          repository.NewPeoplePassRepo,
			"repository.PeopleBidRepo":           repository.NewPeopleBidRepo,
			"repository.PeopleCompanyRepo":       repository.NewPeopleCompanyRepo,
			"repository.BranchRepository":        repository.NewBranchRepository,
			"repository.VehicleCompanyRepo":      repository.NewVehicleCompanyRepo,
			"service.VehiclePassService":         services.NewVehiclePassService,
			"service.VehicleBidService":          services.NewVehicleBidService,
			"services.VehicleIssuedService":      services.NewVehicleIssuedService,
			"services.VehicleClaimService":       services.NewVehicleClaimService,
			"service.PeoplePassService":          services.NewPeoplePassService,
			"service.PeopleBidService":           services.NewPeopleBidService,
			"services.PeopleClaimService":        services.NewPeopleClaimService,
			"services.PeopleCompanyService":      services.NewPeopleCompanyService,
			"services.VehicleCompanyService":     services.NewVehicleCompanyService,
			"services.BranchService":             services.NewBranchService,
			"service.FileService":                services.NewFileService,
			"services.EventService":              services.NewEventService,
			"services.VehicleApplicationService": services.NewVehicleApplicationService,
			"services.PeopleApplicationService":  services.NewPeopleApplicationService,
			"command.Watcher":                    command.NewConsumer,
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
