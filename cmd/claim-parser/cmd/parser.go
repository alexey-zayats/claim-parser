package cmd

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/command"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/di"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var parserCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse",
	Long:  "parse",
	Run:   parserMain,
}

func init() {

	rootCmd.AddCommand(parserCmd)

	viper.SetEnvPrefix(config.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfgParams := []config.Param{
		{Name: "source", Value: "", Usage: "source type for parse: excel, formstruct", ViperBind: "Parser.Source"},
		{Name: "path", Value: "", Usage: "path to file for parse", ViperBind: "Parser.Path"},
	}

	config.Apply(parserCmd, cfgParams)
}

func parserMain(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":            config.NewConfig,
			"command.NewParser": command.NewParser,
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
