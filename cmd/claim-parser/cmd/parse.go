package cmd

import (
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse",
	Long:  "parse",
	Run:   parserMain,
}

func init() {
	rootCmd.AddCommand(parseCmd)

	cfgParams := []config.Param{
		{Name: "path", Value: "", Usage: "path to file for parse", ViperBind: "Parser.Path"},
	}

	config.Apply(parseCmd, cfgParams)

}

func parserMain(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable call cmd.Help")
	}
}
