package command

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/queue"
	"github.com/alexey-zayats/claim-parser/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"sync"
)

// Consumer ...
type Consumer struct {
	config  *config.Config
	queue   *queue.Queue
	wg      *sync.WaitGroup
	appChan chan *model.Application

	vehicleSvc *services.VehicleApplicationService
	peopleSvc  *services.PeopleApplicationService
}

// ConsumerDI ...
type ConsumerDI struct {
	dig.In
	Config *config.Config
	Queue  *queue.Queue

	VehicleSvc *services.VehicleApplicationService
	PeopleSvc  *services.PeopleApplicationService
}

// NewConsumer ...
func NewConsumer(di ConsumerDI) Command {
	return &Consumer{
		config:     di.Config,
		queue:      di.Queue,
		vehicleSvc: di.VehicleSvc,
		peopleSvc:  di.PeopleSvc,
		wg:         &sync.WaitGroup{},
		appChan:    make(chan *model.Application, 1),
	}
}

// Run ...
func (c Consumer) Run(ctx context.Context, args []string) error {
	logrus.WithFields(logrus.Fields{}).Debug("start consumer")

	for i := 0; i < c.config.Amqp.Workers; i++ {
		c.wg.Add(1)
		go c.addWorker(ctx, i)
	}

	c.wg.Add(1)
	go c.consumeVehicle(ctx)

	c.wg.Add(1)
	go c.consumePeople(ctx)

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
		case app := <-c.appChan:
			switch app.Kind {
			case model.KindVehicle:

				logrus.WithFields(logrus.Fields{
					"company":  app.Title,
					"inn": app.Inn,
					"ogrn": app.Ogrn,
					"ceo":    app.CeoName,
				}).Debug("Vehicle.Claim")

				if err := c.vehicleSvc.SaveRecord(app); err != nil {
					logrus.WithFields(logrus.Fields{
						"reason": err,
					}).Error("unable save application")
				}
			case model.KindPeople:

				logrus.WithFields(logrus.Fields{
					"company":  app.Title,
					"inn": app.Inn,
					"ogrn": app.Ogrn,
					"ceo":    app.CeoName,
				}).Debug("People.Claim")

				if err := c.peopleSvc.SaveRecord(app); err != nil {
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
		c.applicationDelivery,
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
		c.applicationDelivery,
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

func (c *Consumer) applicationDelivery(ctx context.Context, delivery amqp.Delivery) {
	var application *model.Application
	dec := json.NewDecoder(bytes.NewReader(delivery.Body))
	if err := dec.Decode(&application); err != nil {
		queue.Nack(delivery, false)
		logrus.WithFields(logrus.Fields{
			"msg":    delivery,
			"reason": err.Error(),
		}).Error("enable decode delivery.Body to model.Application")
		return
	}

	c.appChan <- application

	queue.Ack(delivery)
}
