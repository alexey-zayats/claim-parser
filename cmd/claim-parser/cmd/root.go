package cmd

import (
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "claim-parser",
	Short: "parse claims",
	Long:  "parse claims for pass",
	Run:   rootMain,
}

func rootMain(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable call cmd.Help")
	}
}

// ----------

// Execute entry point to the app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable execute rootCmd")
	}
}

func init() {

	viper.SetEnvPrefix(config.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfgParams := []config.Param{
		{Name: "log-level", Value: "info", Usage: "log level", ViperBind: "Log.Level"},
		{Name: "log-caller", Value: false, Usage: "log caller", ViperBind: "Log.Caller"},
	}

	config.Apply(rootCmd, cfgParams)

	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&config.File, "config", config.FilePath, "Config file")

	cobra.OnInitialize(config.Init)
}
