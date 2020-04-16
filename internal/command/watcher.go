package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/parser/excel"
	"github.com/alexey-zayats/claim-parser/internal/parser/fs"
	godoc_vehicle2 "github.com/alexey-zayats/claim-parser/internal/parser/gsheet"
	"github.com/alexey-zayats/claim-parser/internal/parser/issued"
	"github.com/alexey-zayats/claim-parser/internal/watcher"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// Watcher структура данных команды
type Watcher struct {
	config  *config.Config
	watcher *watcher.Watcher
}

// WatcherParams - DI параметры команды
type WatcherParams struct {
	dig.In
	Config  *config.Config
	Watcher *watcher.Watcher
}

func init() {
	excel.Register()
	fs.Register()
	godoc_vehicle2.Register()
	issued.Register()
}

// NewWatcher - конструктор команды
func NewWatcher(params WatcherParams) Command {
	return &Watcher{
		config:  params.Config,
		watcher: params.Watcher,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *Watcher) Run(ctx context.Context, args []string) error {

	if err := cmd.watcher.Watch(ctx); err != nil {
		return errors.Wrap(err, "unable start watch")
	}

	return nil
}
