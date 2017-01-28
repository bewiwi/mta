package scheduler

import (
	"log"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/models"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func (d *DB) Init(){
	var err error

	db, err := sql.Open(viper.GetString("SCHEDULER.DB.driver"),
		viper.GetString("SCHEDULER.DB.datasource"))
	if err != nil {
		log.Fatal(err)
	}
	d.db = db
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS checks
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

func (d *DB) GetChecks() []models.CheckRequestV1 {
	rows, err := d.db.Query("SELECT id, type, config, freq FROM checks")
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

func init() {
	viper.SetDefault("SCHEDULER.DB.driver", "postgres")
	viper.SetDefault("SCHEDULER.DB.datasource", "user=mta dbname=mta")
}