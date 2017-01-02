package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/bsm/sarama-cluster"
	"github.com/spf13/viper"
)

func GetConfig() *sarama.Config {
	config := sarama.NewConfig()
	//config.Net.TLS.Enable = true
	//config.Net.SASL.Enable = true
	//config.Net.SASL.User = Username
	//config.Net.SASL.Password = Password
	config.Version = sarama.V0_10_0_1
	config.ClientID = viper.GetString("KAFKA.CLIENTID")

	return config
}

func GetConsumer(topic []string) *cluster.Consumer {
	config := GetConfig()
	clusterConfig := cluster.NewConfig()
	clusterConfig.Config = *config
	clusterConfig.Consumer.Return.Errors = true
	clusterConfig.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(
		viper.GetStringSlice("KAFKA.HOSTS"),
		viper.GetString("KAFKA.GROUPID"),
		topic,
		clusterConfig)
	if err != nil {
		log.WithError(err).Fatal("Error during consumption")
	}
	return consumer
}

func GetSyncProducer() sarama.SyncProducer {
	config := GetConfig()

	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("KAFKA.HOSTS"), config)
	if err != nil {
		log.WithError(err).Fatal("Error during consumption")
	}
	return producer
}

type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewProducer() *Producer {
	return &Producer{
		topic:    viper.GetString("KAFKA.TOPIC_ANSWER"),
		producer: GetSyncProducer(),
	}
}

func (p *Producer) SendAnswer(answer *models.CheckAnswer) error {
	log.Debug("Sending answer: ", answer.Values)
	answer.Print()
	var err error
	value, err := json.Marshal(answer)
	if err != nil {
		log.WithError(err).Error("error jsonify")
		return err
	}
	msg := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.StringEncoder(value)}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		log.WithError(err).Error("error sendig")
		return err
	}
	return nil
}
