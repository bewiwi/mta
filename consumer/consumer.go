package consumer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/queue"
	"github.com/spf13/viper"
)

type ConsumerInterface interface {
	Init()
	Push(*models.CheckResponse, func())
	GetValues(check models.CheckV1, startDate int, endDate int) (map[int]float64, error)
}

func GetConsumer() ConsumerInterface {
	consumerType := viper.GetString("CONSUMER_TYPE")
	if consumerType == "INFLUX" {
		influx := InfluxDB{}
		influx.Init()
		return &influx
	} else if consumerType == "STDOUT" {
		stdout := Stdout{}
		stdout.Init()
		return &stdout
	}
	log.Fatal("Invalid CONSUMER_TYPE: ", consumerType)
	return nil
}

func Consume() {
	consume_queue := queue.GetQueue()
	consume_queue.InitResponseConsumer()
	consumer := GetConsumer()

	for {
		checkAnswer, ackFunction := consume_queue.GetNextCheckResponse()
		consumer.Push(checkAnswer, ackFunction)
	}
}
