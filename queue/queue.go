package queue

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
)

type QueueInterface interface {
	Init()
	InitRequestProducer()
	InitResponseProducer()
	InitRequestConsumer()
	InitResponseConsumer()
	PushCheckRequest(models.CheckV1)
	PushCheckResponse(*models.CheckResponse) error
	GetNextCheckRequest() (*models.CheckV1, func())
	GetNextCheckResponse() (*models.CheckResponse, func())
}

type Queue struct{}

func (q Queue) Init()                 {}
func (q Queue) InitRequestProducer()  {}
func (q Queue) InitRequestConsumer()  {}
func (q Queue) InitResponseConsumer() {}
func (q Queue) InitResponseProducer() {}

func GetQueue() QueueInterface {
	queueType := viper.GetString("QUEUE_TYPE")
	if queueType == "KAFKA" {
		kafka := Kafka{}
		kafka.Init()
		return &kafka
	} else if queueType == "RABBITMQ" {
		rabbit := RabbitMQ{}
		rabbit.Init()
		return &rabbit
	}
	log.Fatal("Invalid QUEUE_TYPE: ", queueType)
	return nil
}
