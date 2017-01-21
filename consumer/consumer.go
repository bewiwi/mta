package consumer

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
)

func Consume(f func(models.CheckResponse)error) {
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_ANSWER"))
	defer consumer.Close()

	for msg := range consumer.Messages() {
		var checkAnswer models.CheckResponse
		err := json.Unmarshal(msg.Value, &checkAnswer)
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}
		err = f(checkAnswer)
		if err != nil {
			log.Error("Error sending response")
		} else {
			consumer.MarkOffset(msg, "")
		}
	}
}
