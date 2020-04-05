package command

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/watcher"
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

// NewWatcher - конструктор команды
func NewWatcher(params WatcherParams) Command {
	return &Watcher{
		config:  params.Config,
		watcher: params.Watcher,
	}
}

// Run - имплементация метода Run интерфейса Command
func (cmd *Watcher) Run(ctx context.Context, args []string) error {
	cmd.watcher.Watch(ctx)
	return nil
}
