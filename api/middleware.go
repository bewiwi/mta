package api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"strconv"
)

func (aw *apiWorker) CheckServiceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceId, err := strconv.Atoi(c.Param("serviceId"))
		if err != nil {
			raise404(c, "Wrong service id", err)
			return
		}
		service, err := aw.scheduler.GetService(serviceId)
		if err != nil {
			raise404(c, "Wrong service id", err)
			return
		}
		c.Set("service", service)

	}
}

func (aw *apiWorker) CheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkId, err := strconv.Atoi(c.Param("checkId"))
		if err != nil {
			c.String(http.StatusBadRequest, "Wrong id", checkId)
			return
		}
		check, err := aw.scheduler.GetCheck(checkId)
		if err != nil {
			c.String(http.StatusBadRequest, "Wrong id", checkId)
			return
		}
		c.Set("check", check)
	}
}
