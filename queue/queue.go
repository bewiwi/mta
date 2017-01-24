package queue

import (
	"github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
)

type QueueInterface interface {
	InitProducer()
	PushCheckRequest(models.CheckRequestV1)
	PushCheckResponse(*models.CheckResponse) error
	GetNextCheckRequest() (*models.CheckRequestV1, func())
	GetNextCheckResponse() (*models.CheckResponse, func())
}

type Queue struct {

}

func (q Queue) init(){

}


func GetQueue() QueueInterface {
	queueType := viper.GetString("QUEUE_TYPE")
	if (queueType == "KAFKA") {
		kafka := Kafka{}
		kafka.Init()
		return &kafka
	}
	log.Fatal("Invalid QUEUE_TYPE: ", queueType)
	return nil
}
