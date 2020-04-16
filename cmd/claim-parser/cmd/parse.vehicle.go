package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var vehicleCmd = &cobra.Command{
	Use:   "vehicle",
	Short: "vehicle",
	Long:  "vehicle",
	Run:   vehicleMain,
}

func init() {
	parseCmd.AddCommand(vehicleCmd)

}

func vehicleMain(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Fatal("unable call cmd.Help")
	}
}
