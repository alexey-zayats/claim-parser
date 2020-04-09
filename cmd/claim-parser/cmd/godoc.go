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

var godocCmd = &cobra.Command{
	Use:   "godoc",
	Short: "godoc",
	Long:  "godoc",
	Run:   godocMain,
}

func init() {
	parserCmd.AddCommand(godocCmd)
}

func godocMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                         config.NewConfig,
			"database.NewConnection":         database.NewConnection,
			"repository.NewPassRepository":   repository.NewPassRepository,
			"repository.NewBidRepository":    repository.NewBidRepository,
			"repository.NewIssuedRepository": repository.NewIssuedRepository,
			"service.NewPassService":         services.NewPassService,
			"service.NewBidService":          services.NewBidService,
			"services.NewClaimService":       services.NewClaimService,
			"services.NewIssuedService":      services.NewIssuedService,
			"command.NewGodocParser":         command.NewGodocParser,
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
