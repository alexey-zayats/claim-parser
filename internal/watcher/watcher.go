package watcher

import (
	"context"
	"encoding/json"
	"fmt"
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
	es     *services.EventService
	wg     sync.WaitGroup
	pg     sync.WaitGroup

	out chan interface{}

	events *dict.Dict
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
		fs:     fs,
		es:     params.ES,
		wg:     sync.WaitGroup{},
		pg:     sync.WaitGroup{},
		out:    make(chan interface{}),
		events: dict.New(),
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
		case event, ok := <-w.fs.Events:
			if !ok {
				continue
			}

			w.processEvent(ctx, worker, event)

		case err, ok := <-w.fs.Errors:
			if !ok {
				continue
			}
			logger.Error(err)
		}
	}
}

func (w *Watcher) processEvent(ctx context.Context, worker int, event fsnotify.Event) {

	// we only process Create events
	if event.Op != fsnotify.Create {
		return
	}

	// FIXME: притормозим чутка
	time.Sleep(100 * time.Millisecond)

	path := event.Name

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err, "path": path}).Error("unable read event data")
		return
	}

	logrus.WithFields(logrus.Fields{"data": string(data)}).Debug("event data")

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
	case 3:
		sourceType = "fsdump"
	case 4:
		sourceType = "godoc"
	case 5:
		sourceType = "registry"
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

	params := dict.New()
	params.Set("path", e.Filepath)

	w.events.Set(fmt.Sprintf("worker-%d", worker), e)

	if err := b.Parse(ctx, params, w.out); err != nil {

		state := &model.State{
			ID:     e.FileID,
			Status: 2,
			Error:  err,
		}

		w.es.UpdateState(state)

		logrus.WithFields(logrus.Fields{"reason": err, "path": path}).Error("unable parse")
		return
	}
}

// HandleParsed ...
func (w *Watcher) HandleParsed(ctx context.Context, worker int) {

	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case claim := <-w.out:

			key := fmt.Sprintf("worker-%d", worker)
			if event, ok := w.events.Get(key); ok {

				w.events.Delete(key)

				e := event.(*model.Event)
				e.Claim = claim.(*model.Claim)

				logrus.WithFields(logrus.Fields{"company": e.Claim.Company.Title}).Debug("claim")

				w.es.StoreEvent(e)
			}
		}
	}
}
