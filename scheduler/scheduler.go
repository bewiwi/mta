package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/kafka"
	"github.com/spf13/viper"
	"sync"
	"time"
)

func Run() {
	p := kafka.GetSyncProducer()
	fmt.Println("scheduler called")
	checks := database.GetChecks()
	var wg sync.WaitGroup
	for _, check := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, err := json.Marshal(check)
			if err != nil {
				log.WithError(err).Error("error jsonify")
			}
			msg := &sarama.ProducerMessage{
				Topic: viper.GetString("KAFKA.TOPIC_REQUEST"),
				Value: sarama.StringEncoder(value),
			}
			for {
				_, offset, err := p.SendMessage(msg)
				if err != nil {
					log.WithError(err).Error("error sendig")
				}
				log.Debug(offset)
				time.Sleep(time.Duration(check.Metadata.Freq) * time.Second)
			}
		}()
	}
	wg.Wait()
	fmt.Println("scheduler called")
}
