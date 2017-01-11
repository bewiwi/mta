package worker

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/checks"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_REQUEST"))
	producer := kafka.NewProducer()

	// Can be better
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		<- ch
		log.Debug("Close consumer")
		consumer.Close()
	}()

	// Consume errors
	go func() {
		for err := range consumer.Errors() {
			log.WithError(err).Error("Error during consumption")
		}
	}()

	// Consume notifications
	go func() {
		for note := range consumer.Notifications() {
			log.WithField("note", note).Debug("Rebalanced kafka")
		}
	}()

	log.Info("Ready to consume messages...")
	// Consume messages
	for msg := range consumer.Messages() {

		var checkRequest models.CheckRequestV1
		err := json.Unmarshal(msg.Value, &checkRequest)
		if err != nil {
			log.WithError(err).Error("error unmarchal")
		}

		if checkRequest.Metadata.Type == "ping" {
			var param models.CheckPingParam
			err := json.Unmarshal(*checkRequest.Param, &param)
			if err != nil {
				log.WithError(err).Error("error unmarchal param")
			}

			ping := checks.NewPing(param.Host)
			go func() {
				answer, _ := ping.Run()
				answer.CheckID = checkRequest.Metadata.Id
				err := producer.SendAnswer(answer)
				if err != nil {
					log.Error("Error sending answer")
				} else {
					consumer.MarkOffset(msg, "")
				}
			}()
		} else {
			log.Warn("Message type unknow: ", checkRequest.Metadata.Type)
		}

	}

}
