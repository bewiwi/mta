package database


import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/models"
)
var _db *sql.DB

func getDB() *sql.DB{
	if _db != nil {
		return  _db
	}

	db, err := sql.Open(viper.GetString("DB.driver"),
		viper.GetString("DB.datasource"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetChecks() []models.CheckRequestV1{
	db := getDB()
	rows, err := db.Query("SELECT id, type, config, freqence FROM checks")
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
		//log.Println(string(check.Param))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return checks
}

func init()  {
	viper.SetDefault("DB.driver", "postgres")
	viper.SetDefault("DB.datasource", "user=mta dbname=mta")
}