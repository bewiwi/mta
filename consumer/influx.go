package consumer

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
	"strconv"
)

type InfluxDB struct {
	clnt client.Client
	canConsume chan bool
	bp client.BatchPoints
	ackFuncs []func()
}


func (i *InfluxDB) resetInfluxBp() {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  viper.GetString("CONSUMER.INFLUX.DB"),
		Precision: "us",
	})
	if err != nil {
		log.WithError(err).Fatal("Can't create influx batch point influx")
	}
	i.bp = bp
	i.ackFuncs = []func(){}
}

func (i *InfluxDB) Init() {
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     viper.GetString("CONSUMER.INFLUX.HOST"),
		Username: viper.GetString("CONSUMER.INFLUX.USER"),
		Password: viper.GetString("CONSUMER.INFLUX.PASSWORD"),
	})
	if err != nil {
		log.WithError(err).Fatal("Can't init influx client")
	}
	i.clnt = clnt

	i.resetInfluxBp()
	i.canConsume = make(chan bool)

	go func(){
		for{
			i.canConsume <- false
			log.Debug("Push batch influx : ", len(i.bp.Points()))
			err = i.clnt.Write(i.bp)
			if err != nil {
				log.WithError(err).Fatal("Can't write influx point")
			}
			for _, f := range i.ackFuncs {
				f()
			}
			i.resetInfluxBp()
			i.canConsume <- true
			time.Sleep(5*time.Second)
		}


	}()
}

func (i *InfluxDB) Push(ca *models.CheckResponse, ackFunc func()) {
	if ca.Error != "" {
		//log.Debug("Check in error, pass")
		ackFunc()
		return
	}
	tags := map[string]string{
		"check_id": strconv.Itoa(ca.CheckMetadata.Id),
		"service_id": strconv.Itoa(ca.CheckMetadata.ServiceId),
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

	select {
	case value := <-i.canConsume:
		if value == false {
			//Wait OK
			log.Debug("Wait influx push")
			for {
				value := <- i.canConsume
				if value == true{
					log.Debug("Consume again")
					break
				}
			}

		}
	default:

	}

	i.bp.AddPoint(pt)
	i.ackFuncs = append(i.ackFuncs, ackFunc)
}
