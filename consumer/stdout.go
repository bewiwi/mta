package consumer

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
)

func RunStdout() {
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_ANSWER"))
	for msg := range consumer.Messages() {
		var checkAnswer models.CheckAnswer
		err := json.Unmarshal(msg.Value, &checkAnswer)
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}
		checkAnswer.Print()
	}
}
