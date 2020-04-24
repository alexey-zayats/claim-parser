package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/application"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/queue"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"os"
	"sync"
)

// Consumer ...
type Consumer struct {
	config  *config.Config
	queue   *queue.Queue
	wg      *sync.WaitGroup
	appChan chan interface{}
	csvChan chan interface{}

	vehicleSvc *services.VehicleApplicationService
	peopleSvc  *services.PeopleApplicationService
	singleSvc  *services.SingleApplicationService
}

// ConsumerDI ...
type ConsumerDI struct {
	dig.In
	Config *config.Config
	Queue  *queue.Queue

	VehicleSvc *services.VehicleApplicationService
	PeopleSvc  *services.PeopleApplicationService
	SingleSvc  *services.SingleApplicationService
}

// NewConsumer ...
func NewConsumer(di ConsumerDI) Command {
	return &Consumer{
		config:     di.Config,
		queue:      di.Queue,
		vehicleSvc: di.VehicleSvc,
		peopleSvc:  di.PeopleSvc,
		singleSvc: di.SingleSvc,
		wg:         &sync.WaitGroup{},
		appChan:    make(chan interface{}, 1),
		csvChan:    make(chan interface{}),
	}
}

// Run ...
func (c Consumer) Run(ctx context.Context, args []string) error {
	logrus.WithFields(logrus.Fields{}).Debug("start consumer")

	c.wg.Add(1)
	go c.vehicleCSV(ctx)

	c.wg.Add(1)
	go c.peopleCSV(ctx)

	c.wg.Add(1)
	go c.singleCSV(ctx)

	for i := 0; i < c.config.Amqp.Workers; i++ {
		c.wg.Add(1)
		go c.addWorker(ctx, i)
	}

	c.wg.Add(1)
	go c.consumeVehicle(ctx)

	c.wg.Add(1)
	go c.consumePeople(ctx)

	c.wg.Add(1)
	go c.consumeSingle(ctx)

	c.wg.Wait()

	return nil
}

func (c Consumer) addWorker(ctx context.Context, worker int) {
	defer c.wg.Done()

	logrus.WithFields(logrus.Fields{
		"worker": worker,
	}).Debug("application worker")

	for {
		select {
		case <-ctx.Done():
			return
		case face := <-c.appChan:

			c.csvChan <- face

			switch face.(type) {
			case *application.Vehicle:

				app := face.(*application.Vehicle)

				logrus.WithFields(logrus.Fields{
					"company": app.Title,
					"inn":     app.Inn,
					"ogrn":    app.Ogrn,
					"ceo":     app.CeoName,
				}).Debug("Vehicle.Claim")

				if err := c.vehicleSvc.SaveRecord(app); err != nil {
					logrus.WithFields(logrus.Fields{
						"reason": err,
					}).Error("unable save application")
				}

			case *application.People:

				app := face.(*application.People)

				logrus.WithFields(logrus.Fields{
					"company": app.Title,
					"inn":     app.Inn,
					"ogrn":    app.Ogrn,
					"ceo":     app.CeoName,
				}).Debug("People.Claim")

				if err := c.peopleSvc.SaveRecord(app); err != nil {
					logrus.WithFields(logrus.Fields{
						"reason": err,
					}).Error("unable save application")
				}

			case *application.Single:

				app := face.(*application.Single)

				logrus.WithFields(logrus.Fields{
					"company": app.Title,
					"inn":     app.Inn,
					"ogrn":    app.Ogrn,
					"ceo":     app.CeoName,
				}).Debug("People.Claim")

				if err := c.singleSvc.SaveRecord(app); err != nil {
					logrus.WithFields(logrus.Fields{
						"reason": err,
					}).Error("unable save application")
				}

			}
		}
	}
}

func (c *Consumer) consumeVehicle(ctx context.Context) {
	defer c.wg.Done()

	logrus.WithFields(logrus.Fields{
		"exchange": c.config.Amqp.Exchange,
		"key":      c.config.Amqp.Vehicle.Routing,
		"queue":    c.config.Amqp.Vehicle.Queue,
	}).Debug("start consume vehicle application messages")

	err := c.queue.Consume(ctx, c.config.Amqp.Exchange,
		c.config.Amqp.Vehicle.Routing,
		c.config.Amqp.Vehicle.Queue,
		func(ctx context.Context, delivery amqp.Delivery) {
			var app *application.Vehicle
			dec := json.NewDecoder(bytes.NewReader(delivery.Body))
			if err := dec.Decode(&app); err != nil {
				queue.Nack(delivery, false)
				logrus.WithFields(logrus.Fields{
					"msg":    delivery,
					"reason": err.Error(),
				}).Error("enable decode delivery.Body to model.Application")
				return
			}

			c.appChan <- app

			queue.Ack(delivery)
		},
		1)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"exchange": c.config.Amqp.Exchange,
			"key":      c.config.Amqp.Vehicle.Routing,
			"queue":    c.config.Amqp.Vehicle.Queue,
			"reason":   err.Error(),
		}).Errorf("unable consume application messages: %s", err.Error())
	}
}

func (c *Consumer) consumePeople(ctx context.Context) {
	defer c.wg.Done()

	logrus.WithFields(logrus.Fields{
		"exchange": c.config.Amqp.Exchange,
		"key":      c.config.Amqp.People.Routing,
		"queue":    c.config.Amqp.People.Queue,
	}).Debug("start consume people application messages")

	err := c.queue.Consume(ctx, c.config.Amqp.Exchange,
		c.config.Amqp.People.Routing,
		c.config.Amqp.People.Queue,
		func(ctx context.Context, delivery amqp.Delivery) {
			var app *application.People
			dec := json.NewDecoder(bytes.NewReader(delivery.Body))
			if err := dec.Decode(&app); err != nil {
				queue.Nack(delivery, false)
				logrus.WithFields(logrus.Fields{
					"msg":    delivery,
					"reason": err.Error(),
				}).Error("enable decode delivery.Body to model.Application")
				return
			}

			c.appChan <- app

			queue.Ack(delivery)
		},
		1)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"exchange": c.config.Amqp.Exchange,
			"key":      c.config.Amqp.People.Routing,
			"queue":    c.config.Amqp.People.Queue,
			"reason":   err.Error(),
		}).Errorf("unable consume application messages: %s", err.Error())
	}
}

func (c *Consumer) consumeSingle(ctx context.Context) {
	defer c.wg.Done()

	logrus.WithFields(logrus.Fields{
		"exchange": c.config.Amqp.Exchange,
		"key":      c.config.Amqp.Single.Routing,
		"queue":    c.config.Amqp.Single.Queue,
	}).Debug("start consume single application messages")

	err := c.queue.Consume(ctx, c.config.Amqp.Exchange,
		c.config.Amqp.Single.Routing,
		c.config.Amqp.Single.Queue, func(ctx context.Context, delivery amqp.Delivery) {
			var app *application.Single
			dec := json.NewDecoder(bytes.NewReader(delivery.Body))
			if err := dec.Decode(&app); err != nil {
				queue.Nack(delivery, false)
				logrus.WithFields(logrus.Fields{
					"msg":    delivery,
					"reason": err.Error(),
				}).Error("enable decode delivery.Body to model.Application")
				return
			}

			c.appChan <- app

			queue.Ack(delivery)
		},
		1)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"exchange": c.config.Amqp.Exchange,
			"key":      c.config.Amqp.Single.Routing,
			"queue":    c.config.Amqp.Single.Queue,
			"reason":   err.Error(),
		}).Errorf("unable consume application messages: %s", err.Error())
	}
}

func (c *Consumer) getFile(path string) (*os.File, error) {

	var file *os.File
	var err error

	_, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, errors.Wrap(err, "unable create file")
		}
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "unable open file")
		}
	}

	return file, nil
}

func (c *Consumer) vehicleCSV(ctx context.Context) {
	defer c.wg.Done()

	path := c.config.CSV.Vehicle

	logrus.WithFields(logrus.Fields{"path": path}).Debug("csv.writer")

	file, err := c.getFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"path": path, "reason": err}).Error("unable get file")
		return
	}

	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case face := <-c.csvChan:

			switch face.(type) {
			case *application.Vehicle:
				app := face.(*application.Vehicle)

				for _, p := range app.Passes {
					line := fmt.Sprintf("%v;%d;%d;%s;%d;%d;%s;%s;%s;%s;%s;%s;%s;%d;%d;%d\n",
						app.Dirty,
						app.DistrictID,
						app.PassType,
						app.Title,
						app.Inn,
						app.Ogrn,
						app.CeoName,
						app.CeoPhone,
						app.CeoEmail,
						p.Car,
						p.Lastname,
						p.Firstname,
						p.Middlename,
						app.ActivityKind,
						app.Agreement,
						app.Reliability)

					file.WriteString(line)
				}
			}

		}
	}
}

func (c *Consumer) peopleCSV(ctx context.Context) {
	defer c.wg.Done()

	path := c.config.CSV.People

	logrus.WithFields(logrus.Fields{"path": path}).Debug("csv.writer")

	file, err := c.getFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"path": path, "reason": err}).Error("unable get file")
		return
	}

	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case face := <-c.csvChan:

			switch face.(type) {
			case *application.Vehicle:
				app := face.(*application.People)

				for _, p := range app.Passes {
					line := fmt.Sprintf("%d;%d;%s;%d;%d;%s;%s;%s;%s;%s;%s;%d;%d;%d\n",
						app.DistrictID,
						app.PassType,
						app.Title,
						app.Inn,
						app.Ogrn,
						app.CeoName,
						app.CeoPhone,
						app.CeoEmail,
						p.Lastname,
						p.Firstname,
						p.Middlename,
						app.ActivityKind,
						app.Agreement,
						app.Reliability)

					file.WriteString(line)
				}
			}

		}
	}
}

func (c *Consumer) singleCSV(ctx context.Context) {
	defer c.wg.Done()

	path := c.config.CSV.Single

	logrus.WithFields(logrus.Fields{"path": path}).Debug("csv.writer")

	file, err := c.getFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{"path": path, "reason": err}).Error("unable get file")
		return
	}

	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case face := <-c.csvChan:

			switch face.(type) {
			case *application.Vehicle:
				app := face.(*application.Single)

				for _, p := range app.Passes {
					line := fmt.Sprintf("%d;%d;%s;%d;%d;%s;%s;%s;%s;%s;%s;%s;%d;%d;%d\n",
						app.DistrictID,
						app.PassType,
						app.Title,
						app.Inn,
						app.Ogrn,
						app.CeoName,
						app.CeoPhone,
						app.CeoEmail,
						p.Car,
						p.Lastname,
						p.Firstname,
						p.Middlename,
						app.ActivityKind,
						app.Agreement,
						app.Reliability)

					file.WriteString(line)
				}
			}

		}
	}
}
