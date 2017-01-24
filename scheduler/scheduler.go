package scheduler

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/database"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/queue"
)

type scheduler struct {
	Queue queue.QueueInterface
	Wait  sync.WaitGroup
}

func (s *scheduler) schedule(check models.CheckRequestV1) {
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
	checks := database.GetChecks()
	for _, check := range checks {
		s.Wait.Add(1)
		s.schedule(check)
	}
	s.Wait.Wait()
	log.Debug("Scheduler quit")
}

func Run() {
	scheduler := scheduler{}
	scheduler.Queue = queue.GetQueue()
	scheduler.Queue.InitProducer()
	scheduler.RunLoopSchedule()

}
