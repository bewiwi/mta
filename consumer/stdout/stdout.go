package stdout_consumer

import (
	"github.com/bewiwi/mta/consumer"
	"github.com/bewiwi/mta/models"
)

func Run() {
	consumer.Consume(func(ca *models.CheckResponse) error {
		ca.Print()
		return nil
	})
}
