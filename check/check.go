package check

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
)

type CheckRunInterface interface {
	Run() (*models.CheckResponse, error)
	ValidParam() (error)
}

func HandleError(response *models.CheckResponse, err error) (*models.CheckResponse, error) {
	log.WithError(err).Error("Error on check")
	response.Error = err.Error()
	return response, err
}
