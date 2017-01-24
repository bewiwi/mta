package db_consumer

import (
	"github.com/bewiwi/mta/consumer"
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/models"
)

func Run() {
	consumer.Consume(func(ca *models.CheckResponse) error {
		database.InsertResult(ca)
		return nil
	})
}
