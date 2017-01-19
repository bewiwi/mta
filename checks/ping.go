package checks

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/sparrc/go-ping"
	"time"
)

type Ping struct {
	Host string
}

func NewPing(host string) *Ping {
	return &Ping{
		Host: host,
	}
}

func (p *Ping) Run() (*models.CheckResponse, error) {
	var err error
	response := models.NewCheckResponse()

	log.Debug("PING: ", p.Host)
	pinger, err := ping.NewPinger(p.Host)
	if err != nil {
		return handleError(&response, err)
	}
	pinger.SetPrivileged(true)
	pinger.Timeout = 2 * time.Second
	pinger.Count = 1
	pinger.Run()

	stats := pinger.Statistics()
	if stats.PacketsRecv < pinger.Count {
		err = errors.New(fmt.Sprintf("Timeout (%s)", pinger.Timeout.String()))
		return handleError(&response, err)
	} else {
		response.Values = map[string]float64{
			"rtts": stats.AvgRtt.Seconds(),
		}
	}

	return &response, err

}

func handleError(response *models.CheckResponse, err error) (*models.CheckResponse, error) {
	log.WithError(err).Error("Error on ping")
	response.Error = err.Error()
	return response, err
}
