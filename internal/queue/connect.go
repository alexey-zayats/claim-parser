package queue

import (
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

// Connection ...
type Connection struct {
	conf        *config.Config
	amqp        *amqp.Connection
	closeNotify chan *amqp.Error
	lock        sync.RWMutex
}

// NewConnection ...
func NewConnection(conf *config.Config) *Connection {

	c := &Connection{
		conf: conf,
	}

	return c
}

func (c *Connection) handleREConnect() {

	select {
	case err := <-c.closeNotify:
		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
			"dsn":    c.conf.Amqp.Dsn,
		}).Warn("Lost connection")

		//panic("------ Lost connection ----")
		c.Reconnect()
	}
}

func (c *Connection) establishConnection() *amqp.Connection {

	for {
		connection, err := amqp.Dial(c.conf.Amqp.Dsn)

		if err == nil {
			return connection
		}

		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
			"dns":    c.conf.Amqp.Dsn,
		}).Error("Error connect to AMQP")

		logrus.WithFields(logrus.Fields{
			"dns": c.conf.Amqp.Dsn,
		}).Info("Trying to reconnect to AMQP")

		time.Sleep(500 * time.Millisecond)
	}
}

// Reconnect ...
func (c *Connection) Reconnect() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.amqp = c.establishConnection()
	c.closeNotify = c.amqp.NotifyClose(make(chan *amqp.Error))

	go c.handleREConnect()

}

// Connection ...
func (c *Connection) Connection() *amqp.Connection {

	if c.amqp == nil {
		c.Reconnect()
	}

	return c.amqp
}

// Close ...
func (c *Connection) Close() {
	c.amqp = nil
	//c.closeNotify <- amqp.Close
}
