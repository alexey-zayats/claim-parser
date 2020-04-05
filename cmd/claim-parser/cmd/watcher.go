package cmd

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/command"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/di"
	"github.com/alexey-zayats/claim-parser/internal/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watcherCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch",
	Long:  "watch",
	Run:   watcherMain,
}

func init() {

	rootCmd.AddCommand(watcherCmd)

	cfgParams := []config.Param{

		{Name: "watcher-path-excel", Value: "/tmp/source/excel", Usage: "path to watch for excel", ViperBind: "Watcher.Path.Excel"},
		{Name: "watcher-path-formstruct", Value: "/tmp/source/formstruct", Usage: "path to watch for formstruct", ViperBind: "Watcher.Path.FormStruct"},

		{Name: "sql-dsn", Value: "pass:pass@tcp(127.0.0.1:3306)/pass", Usage: "sql driver", ViperBind: "Sql.Dsn"},
		{Name: "sql-conns-max-idle", Value: 0, Usage: "Maximum number of connections in the idle", ViperBind: "Sql.Conns.Max.Idle"},
		{Name: "sql-conns-max-open", Value: 2, Usage: "Maximum number of open connections to the database", ViperBind: "Sql.Conns.Max.Open"},
		{Name: "sql-conns-max-lifetime", Value: 10, Usage: "Maximum amount of time a connection may be reused", ViperBind: "Sql.Conns.Max.Open"},

		{Name: "watcher-workers", Value: 2, Usage: "number of workers to parse", ViperBind: "Watcher.Workers"},
	}

	config.Apply(watcherCmd, cfgParams)

	viper.AutomaticEnv()

	cobra.OnInitialize(config.Init)

}

func watcherMain(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	di := &di.Runner{
		Provide: map[string]interface{}{
			"config":                 config.NewConfig,
			"database.NewConnection": database.NewConnection,
			"watcher.New":            watcher.New,
			"command.NewWatcher":     command.NewWatcher,
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