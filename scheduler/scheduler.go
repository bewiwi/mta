package scheduler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/kafka"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
)

type scheduler struct {
	Producer sarama.SyncProducer
	Wait     sync.WaitGroup
}

func (s scheduler) schedule(check models.CheckRequestV1) {
	go func() {
		defer s.Wait.Done()

		for {
			log.Debug("Schedule check id : ", check.Metadata.Id)

			check.Metadata.Timestamp = time.Now().Unix()
			value, err := json.Marshal(check)
			if err != nil {
				log.WithError(err).Error("error jsonify")
			}

			msg := &sarama.ProducerMessage{
				Topic: viper.GetString("KAFKA.TOPIC_REQUEST"),
				Value: sarama.StringEncoder(value),
			}

			_, _, err = s.Producer.SendMessage(msg)
			if err != nil {
				log.WithError(err).Error("error sendig")
			}
			time.Sleep(time.Duration(check.Metadata.Freq) * time.Second)
		}
	}()
}

func (s scheduler) RunLoopSchedule() {
	checks := database.GetChecks()
	for _, check := range checks {
		s.Wait.Add(1)
		s.schedule(check)
	}
	s.Wait.Wait()
	log.Debug("Scheduler quit")
}

func Run() {
	scheduler := scheduler{}
	scheduler.Producer = kafka.GetSyncProducer()
	scheduler.RunLoopSchedule()

}
