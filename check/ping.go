package check

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/sparrc/go-ping"
	"time"
)

type Ping struct {
	Host string `json:"host"`
}

func (p *Ping) ValidParam() (error) {
	if p.Host == "" {
		return errors.New("host can't be null")
	}
	return nil
}

func (p *Ping) Run() (*models.CheckResponse, error) {
	var err error
	response := models.NewCheckResponse()

	log.Debug("PING: ", p.Host)
	pinger, err := ping.NewPinger(p.Host)
	if err != nil {
		return HandleError(&response, err)
	}
	pinger.SetPrivileged(true)
	pinger.Timeout = 2 * time.Second
	pinger.Count = 1
	pinger.Run()

	stats := pinger.Statistics()
	if stats.PacketsRecv < pinger.Count {
		err = errors.New(fmt.Sprintf("Timeout (%s)", pinger.Timeout.String()))
		return HandleError(&response, err)
	} else {
		response.Values = map[string]float64{
			"rtts": stats.AvgRtt.Seconds(),
		}
	}

	return &response, err

}
