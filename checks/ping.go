package checks

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sparrc/go-ping"
	"github.com/bewiwi/mta/models"
	"errors"
	"time"
	"fmt"
)

type Ping struct {
	Host string
}

func NewPing(host string) *Ping {
	return &Ping{
		Host: host,
	}
}

func (p *Ping) Run() (*models.CheckAnswer, error){
	var err error
	answer := models.NewCheckAnswer()

	log.Debug("PING: ", p.Host)
	pinger, err := ping.NewPinger(p.Host)
	if err != nil {
		return handleError(&answer, err)
	}
	pinger.Timeout = 2 * time.Second
	pinger.Count = 1
	pinger.Run()

	stats := pinger.Statistics()
	if stats.PacketsRecv < pinger.Count {
		err = errors.New(fmt.Sprintf("Timeout (%s)",pinger.Timeout.String()))
		return handleError(&answer, err)
	}else{
		answer.Values = map[string]float64{
			"rtts": stats.AvgRtt.Seconds(),
		}
	}

	return &answer, err

}

func handleError(answer *models.CheckAnswer,err error) (*models.CheckAnswer,error){
	log.WithError(err).Error("Error on ping")
	answer.Error = err.Error()
	return answer, err
}