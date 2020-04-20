package queue

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Ack send Ack to queue with error logging
func Ack(msg amqp.Delivery) {
	if err := msg.Ack(false); err != nil {
		logrus.WithFields(logrus.Fields{}).Errorf("Error on Ack: %s; ", err.Error())
	}
}

// Nack send Nack to queue with error logging
func Nack(msg amqp.Delivery, requeue bool) {
	if err := msg.Nack(false, requeue); err != nil {
		logrus.WithFields(logrus.Fields{}).Errorf("Error on Nack: %s; ", err.Error())
	}
}
