package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"

	// Register mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/pkg/errors"
)

// NewConnection create new database connection
func NewConnection(config *config.Config) (*sqlx.DB, error) {

	logrus.WithFields(logrus.Fields{
		"dsn": config.SQL.Dsn,
	}).Debug("connection to database")

	db, err := sqlx.Connect("mysql", config.SQL.Dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "unable connect to database with dsn '%s'.", config.SQL.Dsn)
	}

	db.SetConnMaxLifetime(time.Duration(config.SQL.Conns.Max.Lifetime) * time.Second)
	db.SetMaxIdleConns(config.SQL.Conns.Max.Idle)
	db.SetMaxOpenConns(config.SQL.Conns.Max.Open)

	return db, nil
}
