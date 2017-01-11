package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var _db *sql.DB

func getDB() *sql.DB {
	var err error
	if _db != nil {
		return _db
	}

	_db, err = sql.Open(viper.GetString("DB.driver"),
		viper.GetString("DB.datasource"))
	if err != nil {
		log.Fatal(err)
	}
	return _db
}

func init() {
	viper.SetDefault("DB.driver", "postgres")
	viper.SetDefault("DB.datasource", "user=mta dbname=mta")
}
