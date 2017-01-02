package consumer

import (
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/kafka"
	"github.com/spf13/viper"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

func consume(f func(models.CheckAnswer)){
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_ANSWER"))
	for msg := range consumer.Messages() {
		var checkAnswer models.CheckAnswer
		err := json.Unmarshal(msg.Value, &checkAnswer)
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}
		f(checkAnswer)
	}
}