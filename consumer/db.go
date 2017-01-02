package consumer

import (
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/models"
)

func RunDB() {
	consume(func(ca models.CheckAnswer) {
		database.InsertResult(&ca)
	})
}
