package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Param маппинг параметров конфигруации между Cobra и Viper
type Param struct {
	Name      string
	Type      string
	Value     interface{}
	Usage     string
	ViperBind string
}

// Apply связывает viper и cobra в части параметров командной строки
func Apply(cmd *cobra.Command, params []Param) {

	for _, p := range params {
		switch tv := p.Value.(type) {
		case string:
			cmd.PersistentFlags().String(p.Name, p.Value.(string), p.Usage)
		case int:
			cmd.PersistentFlags().Int(p.Name, p.Value.(int), p.Usage)
		case map[string]string:
			cmd.PersistentFlags().StringToString(p.Name, p.Value.(map[string]string), p.Usage)
		case []string:
			cmd.PersistentFlags().StringSlice(p.Name, p.Value.([]string), p.Usage)
		case bool:
			cmd.PersistentFlags().Bool(p.Name, p.Value.(bool), p.Usage)
		default:
			logrus.Panicf("nothing is known about the type %T!", tv)
		}

		if err := viper.BindPFlag(p.ViperBind, cmd.PersistentFlags().Lookup(p.Name)); err != nil {
			if err = cmd.Help(); err != nil {
				logrus.WithFields(logrus.Fields{"reason": err}).Errorf("unable call cmd.Help")
			}
		}

		viper.SetDefault(p.ViperBind, p.Value)
	}
}
