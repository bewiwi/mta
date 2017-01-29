package scheduler

import (
	"log"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/models"
	"encoding/json"
	"io/ioutil"
)

type Json struct {
	db *sql.DB
}

func (c *Json) Init(){}

func (c *Json) GetChecks() []models.CheckRequestV1 {
	var checks []models.CheckRequestV1
	var err error
	file, err := ioutil.ReadFile(viper.GetString("SCHEDULER.JSON.FILE"))
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &checks)
	if err != nil {
		log.Fatal(err)
	}
	return checks
}

func init() {
	viper.SetDefault("SCHEDULER.JSON.FILE", "check.csv")
}