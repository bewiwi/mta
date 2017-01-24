package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"gopkg.in/bsm/sarama-cluster.v2"
)

func GetConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Net.TLS.Enable = viper.GetBool("QUEUE.KAFKA.TLS")
	if viper.GetString("QUEUE.KAFKA.SASL_USER") != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = viper.GetString("QUEUE.KAFKA.SASL_USER")
		config.Net.SASL.Password = viper.GetString("QUEUE.KAFKA.SASL_PASSWORD")
	}
	config.Version = sarama.V0_10_0_1
	if viper.GetString("QUEUE.KAFKA.CLIENTID") != "" {
		config.ClientID = viper.GetString("QUEUE.KAFKA.CLIENTID")
	}

	return config
}

func GetConsumer(topic []string) *cluster.Consumer {
	config := GetConfig()
	clusterConfig := cluster.NewConfig()
	clusterConfig.Config = *config
	clusterConfig.Consumer.Return.Errors = true
	clusterConfig.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(
		viper.GetStringSlice("QUEUE.KAFKA.HOSTS"),
		viper.GetString("QUEUE.KAFKA.GROUPID"),
		topic,
		clusterConfig)
	if err != nil {
		log.WithError(err).Fatal("Error during consumption")
	}
	return consumer
}

func GetSyncProducer() sarama.SyncProducer {
	config := GetConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("QUEUE.KAFKA.HOSTS"), config)
	if err != nil {
		log.WithError(err).Fatal("Error during connecting producer")
	}
	return producer
}

type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewProducer() *Producer {
	return &Producer{
		topic:    viper.GetString("QUEUE.KAFKA.TOPIC_ANSWER"),
		producer: GetSyncProducer(),
	}
}

func (p *Producer) SendResponse(response *models.CheckResponse) error {
	log.Debug("Sending response: ", response.Values)
	var err error
	value, err := json.Marshal(response)
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

func init() {
	viper.SetDefault("QUEUE.KAFKA.TLS", true)
	viper.SetDefault("QUEUE.KAFKA.SASL_USER", "")
	viper.SetDefault("QUEUE.KAFKA.SASL_PASSWORD", "")
}
