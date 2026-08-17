package webhook

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	service  *Service
	clock    Clock
	interval time.Duration
	wg       sync.WaitGroup
}

func NewScheduler(service *Service, clock Clock, interval time.Duration) *Scheduler {
	return &Scheduler{service: service, clock: clock, interval: interval}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := s.clock.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C():
				s.service.DispatchDue(ctx)
			}
		}
	}()
}

func (s *Scheduler) Wait() { s.wg.Wait() }
