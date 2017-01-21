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

func getInfluxBp() client.BatchPoints{
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  viper.GetString("INFLUX.DB"),
		Precision: "us",
	})
	if err != nil {
		log.WithError(err).Fatal("Can't create influx batch point influx")
	}
	return bp
}

func Run() {
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     viper.GetString("INFLUX.HOST"),
		Username: viper.GetString("INFLUX.USER"),
		Password: viper.GetString("INFLUX.PASSWORD"),
	})
	if err != nil {
		log.WithError(err).Fatal("Can't init influx client")
	}

	bp := getInfluxBp()
	canConsume := make(chan bool)

	go func(){
		for{
			time.Sleep(1*time.Second)
			canConsume <- false
			log.Debug("Push batch influx : ", len(bp.Points()))
			err = clnt.Write(bp)
			if err != nil {
				log.WithError(err).Fatal("Can't write influx point")
			}
			bp = getInfluxBp()
			canConsume <- true
		}


	}()

	consumer.Consume(func(ca models.CheckResponse)error{
		select {
		case value := <-canConsume:
			if value == false {
				//Wait OK
				log.Debug("Wait influx push")
				for {
					value := <- canConsume
					if value == true{
						log.Debug("Consume again")
						break
					}
				}

			}
		default:

			if ca.Error != "" {
				log.Debug("Check in error, pass")
				return nil
			}
			tags := map[string]string{
				"check_id": strconv.Itoa(ca.CheckMetadata.Id),
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
				return err
			}
			bp.AddPoint(pt)
		}
		return nil
	})
}
