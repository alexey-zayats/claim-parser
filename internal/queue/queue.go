package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

// Queue ...
type Queue struct {
	conn *Connection
}

// NewQueue ...
func NewQueue(conn *Connection) *Queue {
	return &Queue{
		conn: conn,
	}
}

// NumMessages ...
func (q *Queue) NumMessages(queName string) (int, error) {

	ch, err := q.conn.Connection().Channel()
	if err != nil {
		return 0, errors.Wrap(err, "Channel")
	}

	que, err := ch.QueueInspect(queName)
	if err == nil {
		ch.Close()
		return que.Messages, nil
	}

	return 0, err
}

// Publish ...
func (q *Queue) Publish(exchange string, routing string, data interface{}, headers amqp.Table, props amqp.Table) error {

	ch, err := q.conn.Connection().Channel()
	if err != nil {
		return errors.Wrap(err, "unable open channel")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "unable declare topic exchange %s", exchange)
	}

	msg, _ := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable marshal json for data %#v", data))
	}

	pub := amqp.Publishing{
		Headers:      headers,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(msg),
	}

	for key, value := range props {
		switch key {
		case "ContentType":
			pub.ContentType = value.(string)
		case "ContentEncoding":
			pub.ContentEncoding = value.(string)
		case "DeliveryMode":
			pub.DeliveryMode = value.(uint8)
		case "Priority":
			pub.Priority = value.(uint8)
		case "CorrelationID":
			pub.CorrelationId = value.(string)
		case "ReplyTo":
			pub.ReplyTo = value.(string)
		case "Expiration":
			pub.Expiration = value.(string)
		case "MessageId":
			pub.MessageId = value.(string)
		case "Timestamp":
			pub.Timestamp = value.(time.Time)
		case "Type":
			pub.Type = value.(string)
		case "UserId":
			pub.UserId = value.(string)
		case "AppId":
			pub.AppId = value.(string)
		}
	}

	err = ch.Publish(
		exchange, // exchange
		routing,  // routing key
		false,    // mandatory
		false,    // immediate
		pub,
	)

	if err != nil {
		return errors.Wrapf(err, "unable publish message %#v", msg)
	}

	return nil
}

// Consume ...
func (q *Queue) Consume(
	ctx context.Context,
	exch string,
	rkey string,
	que string,
	handlerFunc func(context.Context, amqp.Delivery),
	qosPrefetch int) error {

	for {
		c := q.conn.Connection()
		lostConnChan := c.NotifyClose(make(chan *amqp.Error))

		ch, err := c.Channel()
		if err != nil {
			return errors.Wrap(err, "unable get channel")
		}

		err = ch.ExchangeDeclare(exch, "topic", true, false, false, false, nil)
		if err != nil {
			return errors.Wrap(err, "enable declare exchange")
		}

		err = ch.Qos(qosPrefetch, 0, false)
		if err != nil {
			return errors.Wrap(err, "unable set qos")
		}

		queue, err := ch.QueueDeclare(que, true, false, false, false, nil)
		if err != nil {
			return errors.Wrap(err, "unable declare queue")
		}

		err = ch.QueueBind(queue.Name, rkey, exch, false, nil)
		if err != nil {
			return errors.Wrap(err, "unable bind queue to exchange")
		}

		deliveries, err := ch.Consume(
			queue.Name, // queue
			"",         // consumer
			false,      // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if err != nil {
			return errors.Wrap(err, "queue.Consume: Consume")
		}

		var delivery amqp.Delivery
		exitLoop := false

		for {
			select {
			case <-ctx.Done():
				ch.Close()
				return errors.New("consume canceled by ctx")
			case err = <-lostConnChan:
				exitLoop = true
				q.conn.Reconnect()
			case delivery = <-deliveries:
				go func(ctx context.Context, delivery amqp.Delivery, handlerFunc func(context.Context, amqp.Delivery)) {
					handlerFunc(ctx, delivery)
				}(ctx, delivery, handlerFunc)
			}
			if exitLoop {
				logrus.Debug("exitLoop")
				break
			}
		}
	}
}
