package watcher

import (
	"context"
	"encoding/json"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/parser"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"io/ioutil"
	"sync"
)

// Watcher структура наблюдателя
type Watcher struct {
	config *config.Config
	db     *sqlx.DB
	wg     *sync.WaitGroup
	fs     *fsnotify.Watcher
	es     *services.EventService
}

// InputParams DI наблюдателя
type InputParams struct {
	dig.In
	Config *config.Config
	DB     *sqlx.DB
	ES     *services.EventService
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
		es:     params.ES,
	}, nil
}

// Watch ...
func (w *Watcher) Watch(ctx context.Context) error {

	for i := 1; i <= w.config.Watcher.Workers; i++ {
		w.wg.Add(1)
		go w.Worker(ctx, i)
	}

	logrus.Debug("start watching events")

	if err := w.fs.Add(w.config.Watcher.Events); err != nil {
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
			w.processEvent(ctx, event)

		case err, ok := <-w.fs.Errors:
			if !ok {
				continue
			}
			logger.Error(err)
		}
	}
}

func (w *Watcher) processEvent(ctx context.Context, event fsnotify.Event) {

	// we only process Create events
	if event.Op != fsnotify.Create {
		return
	}

	path := event.Name

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err, "path": path}).Error("unable read event data")
		return
	}

	var e *model.Event
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err, "data": string(data)}).Error("unable parse json")
		return
	}

	var sourceType string
	switch e.Source {
	case 1:
		sourceType = "excel"
	case 2:
		sourceType = "formstruct"
	default:
		sourceType = "unknown"
	}

	b, err := parser.Instance().Backend(sourceType)
	if err != nil {
		w.es.UpdateState(&model.State{
			ID:     e.FileID,
			Status: 1,
			Error:  err,
		})
		logrus.WithFields(logrus.Fields{"reason": err, "sourceType": sourceType}).Error("unable get parser")
		return
	}

	e.Company, err = b.Parse(ctx, e.Filepath)
	if err != nil {
		w.es.UpdateState(&model.State{
			ID:     e.FileID,
			Status: 2,
			Error:  err,
		})
		logrus.WithFields(logrus.Fields{"reason": err, "path": path}).Error("unable parse")
		return
	}

	w.es.StoreEvent(e)
}
