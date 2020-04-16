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

var peopleGsheetCmd = &cobra.Command{
	Use:   "gsheet",
	Short: "gsheet",
	Long:  "gsheet",
	Run:   peopleGSheetMain,
}

func init() {
	peopleCmd.AddCommand(peopleGsheetCmd)
}

func peopleGSheetMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                        config.NewConfig,
			"database.Connection":           database.NewConnection,
			"repository.PeoplePassRepo":     repository.NewPeoplePassRepo,
			"repository.PeopleBidRepo":      repository.NewPeopleBidRepo,
			"repository.PeopleCompanyRepo":  repository.NewPeopleCompanyRepo,
			"repository.BranchRepository":   repository.NewBranchRepository,
			"repository.FileRepository":     repository.NewFileRepository,
			"service.PeoplePassService":     services.NewPeoplePassService,
			"service.PeopleBidService":      services.NewPeopleBidService,
			"services.PeopleClaimService":   services.NewPeopleClaimService,
			"services.PeopleCompanyService": services.NewPeopleCompanyService,
			"services.BranchService":        services.NewBranchService,
			"services.FileService":          services.NewFileService,
			"command.PeopleGodocParser":     command.NewPeopleGSheetParser,
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
