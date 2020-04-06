package watcher

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"path/filepath"
	"sync"
)

// Watcher структура наблюдателя
type Watcher struct {
	config *config.Config
	db     *sqlx.DB
	wg     *sync.WaitGroup
	fs     *fsnotify.Watcher
}

// InputParams DI наблюдателя
type InputParams struct {
	dig.In
	Config *config.Config
	DB     *sqlx.DB
}

// New создат экземпляр наблюдателя
func New(params InputParams) (*Watcher, error) {

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "unable init fsnotify.Watcher")
	}

	return &Watcher{
		config: params.Config,
		db:     params.DB,
		wg:     &sync.WaitGroup{},
		fs:     fs,
	}, nil
}

// Watch ...
func (w *Watcher) Watch(ctx context.Context) error {

	for i := 1; i <= w.config.Watcher.Workers; i++ {
		w.wg.Add(1)
		go w.Worker(ctx, i)
	}

	logrus.Debug("start watching xlsx")

	if err := w.fs.Add(w.config.Watcher.Path.Excel); err != nil {
		return errors.Wrap(err, "unable add watch")
	}

	logrus.Debug("start watching formstruct")

	if err := w.fs.Add(w.config.Watcher.Path.FormStruct); err != nil {
		return errors.Wrap(err, "unable add watch")
	}

	w.wg.Wait()

	return nil
}

// Worker ...
func (w *Watcher) Worker(ctx context.Context, worker int) {
	defer w.wg.Done()

	logger := logrus.WithFields(logrus.Fields{"worker": worker})
	logger.Debug("start watcher.Worker")

	for {
		select {
		case <-ctx.Done():
			logger.Debug("stop watcher.Worker")
			return
		case event, ok := <-w.fs.Events:
			if !ok {
				continue
			}
			w.processCreateEvent(ctx, event)

		case err, ok := <-w.fs.Errors:
			if !ok {
				continue
			}
			logger.Error(err)
		}
	}
}

func (w *Watcher) processCreateEvent(ctx context.Context, event fsnotify.Event) {

	if event.Op != fsnotify.Create {
		return
	}

	spew.Dump(event)

	path := event.Name
	name := filepath.Base(filepath.Dir(path))

	b, err := parser.Instance().Backend(name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err, "name": name}).Error("unable get parser")
	}

	company, err := b.Parse(ctx, path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err, "path": path}).Error("unable parse")
	}

	spew.Dump(company)
}
