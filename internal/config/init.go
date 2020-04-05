package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var (
	// File ...
	File string

	// FilePath ...
	FilePath = "config/claim-parser.yaml"

	// EnvPrefix ...
	EnvPrefix = "CP"
)

// Init ...
func Init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stderr)
	viper.SetConfigFile(File)
}
