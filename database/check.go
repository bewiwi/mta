package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
)

var onceCreateCheckTable bool

func createCheckTable() {
	// Just call once
	if onceCreateCheckTable == true {
		return
	}
	onceCreateCheckTable = true

	db := getDB()
	_, err := db.Query(`CREATE TABLE IF NOT EXISTS checks
	(
		id serial NOT NULL,
		type character varying(25),
		config jsonb,
		freq integer,
		CONSTRAINT checks_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetChecks() []models.CheckRequestV1 {
	createCheckTable()
	db := getDB()
	rows, err := db.Query("SELECT id, type, config, freq FROM checks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var checks []models.CheckRequestV1
	for rows.Next() {
		check := models.CheckRequestV1{}
		err := rows.Scan(&check.Metadata.Id, &check.Metadata.Type, &check.Param, &check.Metadata.Freq)
		if err != nil {
			log.Fatal(err)
		}
		checks = append(checks, check)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return checks
}
