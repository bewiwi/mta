package database

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
)

var onceCreateResultTable bool

func createResultTable() {
	// Just call once
	if onceCreateResultTable == true {
		return
	}
	onceCreateResultTable = true

	db := getDB()
	_, err := db.Query(`CREATE TABLE IF NOT EXISTS results
	(
		id serial NOT NULL,
		worker VARCHAR,
		timestamp NUMERIC,
		check_id integer ,
		result jsonb,
		CONSTRAINT results_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertResult(c *models.CheckAnswer) {
	createResultTable()
	db := getDB()
	log.Debug("Insert row in result table")
	stmt, err := db.Prepare("INSERT INTO results (worker, check_id, timestamp, result) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	result, err := json.Marshal(c.Values)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(c.Hostname, c.CheckID, c.Timestamp, result)
	if err != nil {
		log.Fatal(err)
	}

}
