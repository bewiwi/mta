package queue

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/spf13/viper"
	"gopkg.in/bsm/sarama-cluster.v2"
)

type Kafka struct {
	Queue
	syncProducer    sarama.SyncProducer
	requestConsumer *cluster.Consumer
}

func (k *Kafka) Init() {
	log.Info("Use KAFKA as Queue system")
}

func (k *Kafka) getConfig() *sarama.Config {
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

func (k *Kafka) getClusterConfig() *cluster.Config {
	config := k.getConfig()
	clusterConfig := cluster.NewConfig()
	clusterConfig.Config = *config
	clusterConfig.Consumer.Return.Errors = true
	clusterConfig.Group.Return.Notifications = true
	return clusterConfig
}

func (k *Kafka) initProducer() {
	log.Debug("Create new producer")
	config := k.getConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("QUEUE.KAFKA.HOSTS"), config)
	if err != nil {
		log.WithError(err).Fatal("Error during connecting producer")
	}
	k.syncProducer = producer
}

func (k *Kafka) InitRequestProducer() {
	k.initProducer()
}

func (k *Kafka) InitResponseProducer() {
	k.initProducer()
}

func (k *Kafka) getConsumer(topic []string) *cluster.Consumer {
	if k.requestConsumer != nil {
		return k.requestConsumer
	}
	log.Debug("Create new request consumer")

	consumer, err := cluster.NewConsumer(
		viper.GetStringSlice("QUEUE.KAFKA.HOSTS"),
		viper.GetString("QUEUE.KAFKA.GROUPID"),
		topic,
		k.getClusterConfig())
	if err != nil {
		log.WithError(err).Fatal("Error during consumption")
	}
	k.requestConsumer = consumer

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

	return consumer
}

func (k *Kafka) PushCheckRequest(check models.CheckV1) {
	value, err := json.Marshal(check)
	if err != nil {
		log.WithError(err).Error("error jsonify")
	}

	msg := &sarama.ProducerMessage{
		Topic: viper.GetString("QUEUE.KAFKA.TOPIC_REQUEST"),
		Value: sarama.StringEncoder(value),
	}

	producer := k.syncProducer
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.WithError(err).Error("error sendig")
	}
}

func (k *Kafka) PushCheckResponse(response *models.CheckResponse) error {
	log.Debug("Sending response: ", response.Values)
	var err error
	value, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).Error("error jsonify")
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: viper.GetString("QUEUE.KAFKA.TOPIC_ANSWER"),
		Value: sarama.StringEncoder(value),
	}
	_, _, err = k.syncProducer.SendMessage(msg)
	if err != nil {
		log.WithError(err).Error("error sending")
		return err
	}
	return nil
}

func (k *Kafka) GetNextCheckRequest() (*models.CheckV1, func()) {
	consumer := k.getConsumer(viper.GetStringSlice("QUEUE.KAFKA.TOPIC_REQUEST"))
	msg := <-consumer.Messages()
	var checkRequest models.CheckV1
	err := json.Unmarshal(msg.Value, &checkRequest)
	if err != nil {
		log.WithError(err).Error("error unmarchal")
	}
	f := func() {
		consumer.MarkOffset(msg, "")
	}
	return &checkRequest, f
}

func (k *Kafka) GetNextCheckResponse() (*models.CheckResponse, func()) {
	consumer := k.getConsumer(viper.GetStringSlice("QUEUE.KAFKA.TOPIC_REQUEST"))
	msg := <-consumer.Messages()

	var checkResponse models.CheckResponse
	err := json.Unmarshal(msg.Value, &checkResponse)
	if err != nil {
		log.WithError(err).Error("error unmarchal check response")
	}
	f := func() {
		consumer.MarkOffset(msg, "")
	}
	return &checkResponse, f
}

func init() {
	viper.SetDefault("QUEUE.KAFKA.TLS", true)
	viper.SetDefault("QUEUE.KAFKA.SASL_USER", "")
	viper.SetDefault("QUEUE.KAFKA.SASL_PASSWORD", "")
}
