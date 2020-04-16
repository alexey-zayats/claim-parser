package cmd

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/command"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/di"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var formstructCmd = &cobra.Command{
	Use:   "formstruct",
	Short: "formstruct",
	Long:  "formstruct",
	Run:   formstructMain,
}

func init() {
	vehicleCmd.AddCommand(formstructCmd)
}

func formstructMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":            config.NewConfig,
			"command.NewParser": command.NewVehicleFSParser,
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
