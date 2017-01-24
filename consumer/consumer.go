package consumer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/queue"
)

func Consume(f func(*models.CheckResponse)error) {
	var err error
	consumer := kafka.GetConsumer(viper.GetStringSlice("QUEUE.KAFKA.TOPIC_ANSWER"))
	defer consumer.Close()

	for {
		checkAnswer, ackFunction := queue.GetQueue().GetNextCheckResponse()
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}
		err = f(checkAnswer)
		if err != nil {
			log.Error("Error sending response")
		} else {
			ackFunction()
		}
	}
}
