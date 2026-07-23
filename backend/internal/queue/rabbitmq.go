package queue

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewRabbitmq(url string) (*Rabbitmq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		"pipeline_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{
		conn:    conn,
		channel: channel,
		queue:   &queue,
	}, nil
}

func (r *Rabbitmq) PublishPipeline(ctx context.Context, pipelineId int64) error {
	message := PipelineMessage{
		PipelineId: pipelineId,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Priority:     0,
		})
	if err != nil {
		return err
	}
	log.Printf("Pipeline %d published successfully", pipelineId)
	return nil
}

func (r *Rabbitmq) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
