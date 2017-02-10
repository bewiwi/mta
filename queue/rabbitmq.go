package queue

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Queue
	conn            *amqp.Connection
	channel         *amqp.Channel
	requestConsume  <-chan amqp.Delivery
	responseConsume <-chan amqp.Delivery
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func (r *RabbitMQ) Init() {
	log.Info("Use RabbitMQ as Queue system")
	conn, err := amqp.Dial(viper.GetString("QUEUE.RABBITMQ.HOST"))
	r.failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	r.conn = conn

	ch, err := r.conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	//defer ch.Close()
	r.channel = ch
}

func (r *RabbitMQ) initRequestConsumer(queue_name string) <-chan amqp.Delivery {
	q, err := r.channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	r.failOnError(err, "Failed to declare a queue")

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	r.failOnError(err, "Failed to register a consumer")
	return msgs
}
func (r *RabbitMQ) InitRequestConsumer() {
	r.requestConsume = r.initRequestConsumer(viper.GetString("QUEUE.RABBITMQ.REQUEST_QUEUE"))
}

func (r *RabbitMQ) InitResponseConsumer() {
	r.requestConsume = r.initRequestConsumer(viper.GetString("QUEUE.RABBITMQ.RESPONSE_QUEUE"))
}

func (r *RabbitMQ) PushCheckRequest(check models.CheckV1) {
	value, err := json.Marshal(check)
	if err != nil {
		log.WithError(err).Error("error jsonify")
	}

	err = r.channel.Publish(
		"", // exchange
		viper.GetString("QUEUE.RABBITMQ.REQUEST_QUEUE"), // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(value),
		})
	r.failOnError(err, "Failed to publish a message")
}

func (r *RabbitMQ) PushCheckResponse(check *models.CheckResponse) error {
	value, err := json.Marshal(check)
	if err != nil {
		log.WithError(err).Error("error jsonify")
	}

	err = r.channel.Publish(
		"", // exchange
		viper.GetString("QUEUE.RABBITMQ.RESPONSE_QUEUE"), // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(value),
		})
	r.failOnError(err, "Failed to publish a message")
	return err
}

func (r *RabbitMQ) GetNextCheckRequest() (*models.CheckV1, func()) {
	var err error

	msg := <-r.requestConsume
	var checkRequest models.CheckV1
	err = json.Unmarshal(msg.Body, &checkRequest)
	if err != nil {
		log.WithError(err).Error("error unmarchal")
	}
	f := func() {
		msg.Ack(false)
	}
	return &checkRequest, f
}

func (r *RabbitMQ) GetNextCheckResponse() (*models.CheckResponse, func()) {
	var err error

	msg := <-r.requestConsume
	var checkResponse models.CheckResponse
	err = json.Unmarshal(msg.Body, &checkResponse)
	if err != nil {
		log.WithError(err).Error("error unmarchal")
	}
	f := func() {
		msg.Ack(false)
	}
	return &checkResponse, f
}

func init() {
	viper.SetDefault("QUEUE.RABBITMQ.REQUEST_QUEUE", "mta-request")
	viper.SetDefault("QUEUE.RABBITMQ.RESPONSE_QUEUE", "mta-response")
}
