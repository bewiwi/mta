package api

import (
	"github.com/bewiwi/mta/consumer"
	"github.com/bewiwi/mta/scheduler"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

type apiWorker struct {
	scheduler scheduler.SchedulerInterface
	consumer  consumer.ConsumerInterface
}

func (aw *apiWorker) RunWorker() {
	router := gin.Default()
	// Simple group: v1
	v1 := router.Group("/v1")
	{
		// Service
		v1.POST("/services", aw.postService)
		v1Services := v1.Group("/services/:serviceId")
		{
			v1Services.Use(aw.CheckServiceMiddleware())
			v1Services.GET("", aw.getService)
			v1Services.DELETE("", aw.deleteService)
			v1Services.PUT("", aw.putService)
			v1Checks := v1Services.Group("/checks")
			{
				v1Checks.GET("", aw.getServiceChecks)
				v1Checks.POST("", aw.postCheck)
				v1CheckId := v1Checks.Group("/:checkId")
				{
					v1CheckId.Use(aw.CheckMiddleware())
					v1CheckId.GET("", aw.getCheck)
					v1CheckId.DELETE("", aw.deleteCheck)
				}
			}
		}
	}

	router.Run(viper.GetString("API.LISTEN"))
}

func Run() {
	worker := apiWorker{
		scheduler: scheduler.GetScheduler(),
		consumer:  consumer.GetConsumer(),
	}
	worker.RunWorker()
	defer worker.scheduler.Close()
}

func init() {
	viper.SetDefault("API.LISTEN", ":8099")
}
