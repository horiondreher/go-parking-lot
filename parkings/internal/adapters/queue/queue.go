package queue

import (
	"github.com/horiondreher/go-parking-lot/parkings/internal/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

const UserUpdatedChannel = "user_updates"

type QueueAdapter struct {
	config *utils.Config
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	queue  amqp091.Queue
}

func NewQueueAdapter() (*QueueAdapter, error) {
	config := utils.GetConfig()

	conn, err := amqp091.Dial(config.QueueServerAddress)
	if err != nil {
		log.Err(err).Msg("failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Err(err).Msg("failed to open a channel")
	}

	queue, err := ch.QueueDeclare(UserUpdatedChannel, true, false, false, false, nil)
	if err != nil {
		log.Err(err).Msg("failed to decalre a queue")
	}

	log.Info().Msg("connected to RabbitMQ and queue declared")

	queueAdapter := &QueueAdapter{
		config: config,
		conn:   conn,
		ch:     ch,
		queue:  queue,
	}

	return queueAdapter, nil
}

func (q *QueueAdapter) PublishOnUserUpdated(msg string) error {
	err := q.ch.Publish(
		"",
		q.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)

	return err
}
