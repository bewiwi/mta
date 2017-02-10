package api

import (
	"github.com/bewiwi/mta/models"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type ServiceAPI struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

func (aw *apiWorker) postService(c *gin.Context) {
	var service models.Service
	err := c.BindJSON(&service)
	if err != nil {
		raise400(c, "Parameters wrong", err)
		return
	}

	service, err = aw.scheduler.CreateService(service)
	if err != nil {
		c.String(http.StatusBadRequest, "Wrong id", err.Error())
		raise500(c, "Insert erro", err)
		return
	}
	c.JSON(http.StatusOK, service)

}

func (aw *apiWorker) getService(c *gin.Context) {
	serviceContext, _ := c.Get("service")
	service := serviceContext.(models.Service)
	c.JSON(http.StatusOK, service)
}

func (aw *apiWorker) deleteService(c *gin.Context) {
	serviceContext, _ := c.Get("service")
	service := serviceContext.(models.Service)
	aw.scheduler.DeleteService(service.Id)
	c.JSON(http.StatusOK, service)
}

func (aw *apiWorker) putService(c *gin.Context) {
	serviceContext, _ := c.Get("service")
	service := serviceContext.(models.Service)

	var serviceData models.Service
	err := c.BindJSON(&serviceData)
	if err != nil {
		raise400(c, "Parameters wrong", err)
		return
	}
	service, err = aw.scheduler.UpdateService(service.Id, serviceData)
	if err != nil {
		c.String(http.StatusBadRequest, "Wrong id", err.Error())
		raise500(c, "Insert erro", err)
		return
	}
	c.String(http.StatusOK, "Hello %s", service)

}
