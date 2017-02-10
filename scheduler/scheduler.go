package scheduler

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/queue"
	"github.com/spf13/viper"
)

type SchedulerInterface interface {
	Init()
	Close()
	CreateService(service models.Service) (models.Service, error)
	GetService(id int) (models.Service, error)
	DeleteService(id int) error
	UpdateService(id int, service models.Service) (models.Service, error)
	CreateCheck(check models.CheckV1) (models.CheckV1, error)
	DeleteCheck(id int) error
	GetChecks(serviceId int) ([]models.CheckV1, error)
	GetAllChecks() ([]models.CheckV1, error)
	GetCheck(checkId int) (models.CheckV1, error)
}

type scheduler struct {
	Queue queue.QueueInterface
	Wait  sync.WaitGroup
}

func (s *scheduler) schedule(check models.CheckV1) {
	go func() {
		defer s.Wait.Done()

		for {
			log.Debug("Schedule check id : ", check.Metadata.Id)

			check.Metadata.Timestamp = time.Now().Unix()
			s.Queue.PushCheckRequest(check)
			time.Sleep(time.Duration(check.Metadata.Freq) * time.Second)
		}
	}()
}

func (s *scheduler) RunLoopSchedule() {
	scheduler := GetScheduler()
	defer scheduler.Close()
	checks, err := scheduler.GetAllChecks()
	if err != nil {
		log.WithError(err).Fatal("Can't start scheduler")
	}
	for _, check := range checks {
		s.Wait.Add(1)
		s.schedule(check)
	}
	s.Wait.Wait()
	log.Debug("Scheduler quit")
}

func GetScheduler() SchedulerInterface {
	schedulerType := viper.GetString("SCHEDULER_TYPE")
	if schedulerType == "DB" {
		db := DB2{}
		db.Init()
		return &db
	}
	log.Fatal("Invalid SCHEDULER_TYPE: ", schedulerType)
	return nil
}

func Run() {
	scheduler := scheduler{}
	scheduler.Queue = queue.GetQueue()
	scheduler.Queue.InitRequestProducer()
	scheduler.RunLoopSchedule()

}
