package consumer

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
)

func Consume(f func(models.CheckResponse)) {
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_ANSWER"))
	defer consumer.Close()

	for msg := range consumer.Messages() {
		var checkAnswer models.CheckResponse
		err := json.Unmarshal(msg.Value, &checkAnswer)
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}
		f(checkAnswer)
	}
}
