package db_consumer

import (
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/consumer"
)

func Run() {
	consumer.Consume(func(ca models.CheckResponse) {
		database.InsertResult(&ca)
	})
}
