package consumer

import (
	"errors"
	"github.com/bewiwi/mta/models"
)

type Stdout struct{}

func (s *Stdout) Init() {}

func (s *Stdout) Push(ca *models.CheckResponse, ackFunc func()) {
	ca.Print()
	ackFunc()
}

func (s *Stdout) GetValues(check models.CheckV1, startDate int, endDate int) (map[int]float64, error) {
	return map[int]float64{}, errors.New("Stdout backend not compatible")
}
