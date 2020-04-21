package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Log struct {
		Level  string
		Caller bool
	}
	SQL struct {
		Dsn   string
		Conns struct {
			Max struct {
				Idle     int
				Open     int
				Lifetime int
			}
		}
	}
	Amqp struct {
		Dsn      string
		Exchange string
		Workers  int
		Vehicle  struct {
			Routing string
			Queue   string
		}
		People struct {
			Routing string
			Queue   string
		}
	}
	Watcher struct {
		Workers int
		Events  string
	}
	Parser struct {
		Sheet  string
		Source string
		Path   string
	}
	Pass struct {
		Source  int
		Creator int64
		Clean   int64
		Dirty   int64
	}
	CSV string
}

// NewConfig ...
func NewConfig() (*Config, error) {

	var level logrus.Level
	var err error
	config := &Config{}

	if err = viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal config")
	}

	level, err = logrus.ParseLevel(config.Log.Level)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal config")
	}

	logrus.SetLevel(level)
	logrus.SetReportCaller(config.Log.Caller)

	return config, nil
}
