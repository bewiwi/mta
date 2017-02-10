package api

import (
	"net/http"

	"github.com/bewiwi/mta/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func (aw *apiWorker) getServiceChecks(c *gin.Context) {
	serviceContext, _ := c.Get("service")
	service := serviceContext.(models.Service)
	checks, err := aw.scheduler.GetChecks(service.Id)
	if err != nil {
		raise400(c, "Wrong id", err)
		return
	}
	c.JSON(http.StatusOK, checks)
}

func (aw *apiWorker) getCheck(c *gin.Context) {
	checkContext, _ := c.Get("check")
	check := checkContext.(models.CheckV1)
	c.JSON(http.StatusOK, check)
}

func (aw *apiWorker) deleteCheck(c *gin.Context) {
	checkContext, _ := c.Get("check")
	check := checkContext.(models.CheckV1)
	err := aw.scheduler.DeleteCheck(check.Metadata.Id)
	if err != nil {
		raise400(c, "Error on deletion", err)
		return
	}
	c.JSON(http.StatusOK, check)
}

func (aw *apiWorker) postCheck(c *gin.Context) {
	var check CheckApi
	err := c.BindJSON(&check)
	if err != nil {
		raise400(c, "Parameters wrong", err)
		return
	}
	modelCheck, err := check.GetCheck()
	if err != nil {
		raise400(c, "Can't convert", err)
		return
	}

	serviceContext, _ := c.Get("service")
	service := serviceContext.(models.Service)
	modelCheck.Metadata.ServiceId = service.Id
	modelCheck, err = aw.scheduler.CreateCheck(modelCheck)
	if err != nil {
		c.String(http.StatusBadRequest, "Wrong id", err.Error())
		raise500(c, "Insert error", err)
		return
	}
	c.JSON(http.StatusOK, modelCheck)

}