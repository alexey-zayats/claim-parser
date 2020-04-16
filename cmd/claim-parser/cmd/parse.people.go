package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var peopleCmd = &cobra.Command{
	Use:   "people",
	Short: "people",
	Long:  "people",
	Run:   peopleMain,
}

func init() {
	parseCmd.AddCommand(peopleCmd)

}

func peopleMain(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable call cmd.Help")
	}
}
