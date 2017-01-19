package influx_consumer

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/consumer"
	"github.com/bewiwi/mta/models"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
	"strconv"
)

func Run() {
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     viper.GetString("INFLUX.HOST"),
		Username: viper.GetString("INFLUX.USER"),
		Password: viper.GetString("INFLUX.PASSWORD"),
	})
	if err != nil {
		log.WithError(err).Fatal("Can't init influx client")
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  viper.GetString("INFLUX.DB"),
		Precision: "us",
	})
	if err != nil {
		log.WithError(err).Fatal("Can't create influx batch point influx")
	}

	consumer.Consume(func(ca models.CheckResponse) {
		if ca.Error != "" {
			log.Debug("Check in error, pass")
			return
		}
		tags := map[string]string{
			"check_id": strconv.Itoa(123),
			"worker":   ca.Hostname,
		}

		fields := make(map[string]interface{})
		for key, value := range ca.Values {
			fields[key] = value
		}

		pt, err := client.NewPoint(
			ca.CheckMetadata.Type,
			tags,
			fields,
			time.Unix(ca.Timestamp/1000000000, 0),
		)

		if err != nil {
			log.WithError(err).Fatal("Can't create influx point")
		}

		bp.AddPoint(pt)

		log.Debug("Push influx")
		err = clnt.Write(bp)
		if err != nil {
			log.WithError(err).Fatal("Can't write influx point")
		}
	})
}
