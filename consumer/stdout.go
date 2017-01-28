package consumer

import (
	"github.com/bewiwi/mta/models"
)

type Stdout struct {}

func (s *Stdout) Init(){}

func (s *Stdout) Push(ca *models.CheckResponse, ackFunc func()) {
	ca.Print()
	ackFunc()
}
