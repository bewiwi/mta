package worker

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/check"
	"github.com/bewiwi/mta/check/ping"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/check/http"
)

func Run() {
	consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_REQUEST"))
	producer := kafka.NewProducer()

	// Can be better
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
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

		var current_check check.CheckRun
		if checkRequest.Metadata.Type == "ping" {
			convert_check := ping_check.Ping{}
			err = json.Unmarshal(*checkRequest.Param, &convert_check)
			if err != nil {
				log.WithError(err).Error("error unmarchal param")
			}
			current_check = &convert_check
		} else if checkRequest.Metadata.Type == "http" {
			convert_check := http_check.HttpCheck{}
			err = json.Unmarshal(*checkRequest.Param, &convert_check)
			if err != nil {
				log.WithError(err).Error("error unmarchal param")
			}
			current_check = &convert_check
		} else {
			log.Warn("Message type unknow: ", checkRequest.Metadata.Type)
			continue

		}

		err = json.Unmarshal(*checkRequest.Param, &current_check)
		if err != nil {
			log.WithError(err).Error("error unmarchal param")
		}

		go func() {
			response, _ := current_check.Run()
			response.CheckMetadata = checkRequest.Metadata
			err := producer.SendResponse(response)
			if err != nil {
				log.Error("Error sending response")
			} else {
				consumer.MarkOffset(msg, "")
			}
		}()

	}

}
