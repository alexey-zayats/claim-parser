package watcher

import (
	"context"
	"encoding/json"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/dict"
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
	"time"
)

// Watcher структура наблюдателя
type Watcher struct {
	config *config.Config
	db     *sqlx.DB
	fs     *fsnotify.Watcher
	wg     sync.WaitGroup

	eventService *services.EventService

	out chan interface{}
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
		config:       params.Config,
		db:           params.DB,
		fs:           fs,
		eventService: params.ES,
		wg:           sync.WaitGroup{},
		out:          make(chan interface{}, 1),
	}, nil
}

// Watch ...
func (w *Watcher) Watch(ctx context.Context) error {

	for i := 1; i <= w.config.Watcher.Workers; i++ {
		w.wg.Add(1)
		go w.Worker(ctx, i)

		w.wg.Add(1)
		go w.HandleParsed(ctx, i)
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
		case event := <-w.fs.Events:

			if err := w.processEvent(ctx, event); err != nil {
				logger.WithFields(logrus.Fields{"reason": err}).Error("unable porcess event")
			}

		case err := <-w.fs.Errors:
			logger.Error(err)
		}
	}
}

func (w *Watcher) processEvent(ctx context.Context, e fsnotify.Event) error {

	// we only process Create events
	if e.Op != fsnotify.Create {
		return nil
	}

	// FIXME: притормозим чутка
	time.Sleep(100 * time.Millisecond)

	data, err := ioutil.ReadFile(e.Name)
	if err != nil {
		return errors.Wrapf(err, "unable read event data in path %s", e.Name)
	}

	logrus.WithFields(logrus.Fields{"data": string(data)}).Debug("event data")

	event := &model.Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return errors.Wrap(err, "unable parse json")
	}

	var sourceType string
	switch event.Source {
	case 1:
		sourceType = "excel"
	case 2:
		sourceType = "formstruct"
	case 3:
		sourceType = "fsdump"
	case 4:
		sourceType = "godoc"
	case 5:
		sourceType = "issued"
	default:
		sourceType = "unknown"
	}

	b, err := parser.Instance().Backend(sourceType)
	if err != nil {
		w.eventService.UpdateFile(event.FileID, 1, err.Error(), sourceType)
		return errors.Wrap(err, "unable get parser")
	}

	params := dict.New()
	params.Set("path", event.Filepath)
	params.Set("event", event)

	if err := b.Parse(ctx, params, w.out); err != nil {
		w.eventService.UpdateFile(event.FileID, 2, err.Error(), sourceType)
		logrus.WithFields(logrus.Fields{"reason": err, "path": e.Name}).Error("unable parse")
	}

	return nil
}

// HandleParsed ...
func (w *Watcher) HandleParsed(ctx context.Context, worker int) {

	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case face := <-w.out:

			switch face.(type) {
			case *model.Claim:

				record := face.(*model.Claim)
				w.eventService.StoreClaim(record)

			case *model.Registry:

				record := face.(*model.Registry)
				w.eventService.StoreRegistry(record)

			case nil:
				continue
			}

		}
	}
}
