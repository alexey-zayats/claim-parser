package watcher

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

// Watcher структура наблюдателя
type Watcher struct {
	config *config.Config
	db     *sqlx.DB
}

// InputParams DI наблюдателя
type InputParams struct {
	dig.In
	Config *config.Config
	DB     *sqlx.DB
}

// New создат экземпляр наблюдателя
func New(params InputParams) *Watcher {
	return &Watcher{
		config: params.Config,
		db:     params.DB,
	}
}

// Watch ...
func (w *Watcher) Watch(ctx context.Context) error {

	logrus.Debug("start watching")

	for {
		select {
		case <-ctx.Done():
			logrus.Debug("watch canceled")
			return nil
		}
	}
}
