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

var resgistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "registry",
	Long:  "registry",
	Run:   registryMain,
}

func init() {

	parserCmd.AddCommand(resgistryCmd)

	cfgParams := []config.Param{
		{Name: "sheet", Value: "", Usage: "sheet name for parse", ViperBind: "Parser.Sheet"},
	}

	config.Apply(resgistryCmd, cfgParams)
}

func registryMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                       config.NewConfig,
			"database.NewConnection":       database.NewConnection,
			"repository.NewPassRepository": repository.NewPassRepository,
			"repository.NewBidRepository":  repository.NewBidRepository,
			"service.NewPassService":       services.NewPassService,
			"service.NewBidService":        services.NewBidService,
			"services.NewRegistryService":  services.NewRegistryService,
			"command.NewRegistryParser":    command.NewRegistryParser,
		},
		Invoke: func(ctx context.Context, args []string) interface{} {
			return func(i command.Command) {
				if err := i.Run(ctx, args); err != nil {
					logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable run command")
				}
			}
		},
	}

	di.Run(ctx, di, args)
}
