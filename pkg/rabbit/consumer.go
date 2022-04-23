package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/perfectogo/template-service-with-mq/models"
	"github.com/streadway/amqp"
)

// Consumer ...
type Consumer struct {
	consumerName string
	channel      *amqp.Channel
	exchangeName string
	routingKey   string
	queueName    string
	handler      func(amqp.Delivery) models.SocketResponse
	messages     <-chan amqp.Delivery
	errors       chan error
}

// AddConsumer ...
func (rmq *RMQ) AddConsumer(consumerName, exchangeName, queueName, routingKey string, handler func(amqp.Delivery) models.SocketResponse) {
	if rmq.consumers[consumerName] != nil {
		panic(errors.New("consumer with the same name already exists: " + consumerName))
	}

	ch, err := rmq.conn.Channel()

	if err != nil {
		panic(err)
	}

	err = declareExchange(ch, exchangeName)

	if err != nil {
		fmt.Printf("Exchange Declare: %s", err.Error())
		panic(err)
	}

	q, err := declareQueue(ch, queueName)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	messages, err := ch.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		true,
		nil,
	)

	if err != nil {
		panic(err)
	}

	rmq.consumers[consumerName] = &Consumer{
		consumerName: consumerName,
		channel:      ch,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		queueName:    queueName,
		handler:      handler,
		messages:     messages,
		errors:       rmq.consumerErrors,
	}

	return
}

// Start ...
func (c *Consumer) Start(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.messages:
			if !ok {
				panic(errors.New("error while reading consumer messages"))
			} else {
				resp := c.handler(msg)
				resp.CorrelationID = msg.MessageId

				b, e := json.Marshal(resp)

				var respToSocket amqp.Publishing

				if e != nil {
					c.errors <- e
				} else {
					respToSocket = amqp.Publishing{
						Headers:       msg.Headers,
						Type:          "JSON",
						ContentType:   "application/json",
						DeliveryMode:  amqp.Persistent,
						CorrelationId: msg.MessageId,
						Body:          b,
					}
				}

				if len(resp.Error.Code) < 0 {
					c.pushReplay(msg.ReplyTo, respToSocket)
					msg.Nack(false, false)
				} else {
					c.pushReplay(msg.ReplyTo, respToSocket)
					msg.Ack(false)
				}
			}
		case <-ctx.Done():
			{
				err := c.channel.Cancel("", true)

				if err != nil {
					c.errors <- err
				}

				return
			}
		}
	}
}

// PushReplay ...
func (c *Consumer) pushReplay(replyTo string, msg amqp.Publishing) {

	err := c.channel.Publish(
		c.exchangeName,
		replyTo,
		false,
		false,
		msg,
	)

	if err != nil {
		c.errors <- err
	}

	return
}
