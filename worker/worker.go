package worker

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/check"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/queue"
	"time"
)

func Run() {
	var err error
	queue := queue.GetQueue()
	queue.InitResponseProducer()
	queue.InitRequestConsumer()

	log.Info("Ready to consume messages...")
	// Consume messages
	for {
		checkRequest, ackFunc := queue.GetNextCheckRequest()

		//Check must be run ?
		if checkMustBeRun(checkRequest.Metadata) == false {
			log.Warn("Check request to old")
			ackFunc()
			continue
		}
		var current_check check.CheckRunInterface
		if checkRequest.Metadata.Type == "ping" {
			convert_check := check.Ping{}
			err = json.Unmarshal(*checkRequest.Param, &convert_check)
			if err != nil {
				log.WithError(err).Error("error unmarchal param")
			}
			current_check = &convert_check
		} else if checkRequest.Metadata.Type == "http" {
			convert_check := check.HttpCheck{}
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
			err := queue.PushCheckResponse(response)
			if err != nil {
				log.Error("Error sending response")
			} else {
				ackFunc()
			}
		}()
	}
}

func checkMustBeRun(metadata models.CheckMetadaV1) bool {
	toLate := metadata.Timestamp + int64(metadata.Freq)
	if time.Now().Unix() > toLate {
		log.Debug(toLate)
		return false
	}
	return true
}
