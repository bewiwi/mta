package consumer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/queue"
)

func Consume(f func(*models.CheckResponse)error) {
	var err error
	consume_queue := queue.GetQueue()
	consume_queue.InitResponseConsumer()

	for {
		checkAnswer, ackFunction := consume_queue.GetNextCheckResponse()
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
