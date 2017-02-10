package api

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func raise(c *gin.Context, httpCode int, msg string, err error) {
	if err != nil {
		log.WithError(err).Info(msg)
	}
	c.JSON(httpCode, gin.H{"message": msg})
	c.Abort()
}

func raise400(c *gin.Context, msg string, err error) {
	raise(c, http.StatusBadRequest, msg, err)
}

func raise404(c *gin.Context, msg string, err error) {
	raise(c, http.StatusNotFound, msg, err)
}

func raise500(c *gin.Context, msg string, err error) {
	raise(c, http.StatusInternalServerError, msg, err)

}
