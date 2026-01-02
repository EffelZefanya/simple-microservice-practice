package platform

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare("orders_queue", true, false, false, false, nil)
	return &RabbitMQ{conn: conn, channel: ch}, err
}

func (r *RabbitMQ) PublishOrder(ctx context.Context, event interface{}) error {
	body, _ := json.Marshal(event)
	return r.channel.PublishWithContext(ctx,
		"",             // exchange
		"orders_queue", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}