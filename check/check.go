package check

import (
	"github.com/bewiwi/mta/models"
	log "github.com/Sirupsen/logrus"
)

type CheckRun interface {
	Run() (*models.CheckResponse, error)
}


func HandleError(response *models.CheckResponse, err error) (*models.CheckResponse, error) {
	log.WithError(err).Error("Error on check")
	response.Error = err.Error()
	return response, err
}
