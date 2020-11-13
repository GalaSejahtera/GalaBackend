package utility

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"safeworkout/pkg/logger"
	"sync"
	"time"
)

// Scheduler ...
type Scheduler struct {
	Enabled   bool
	HaveStart bool
	Job       func(ctx context.Context) error
	Wg        sync.WaitGroup
}

func (it *Scheduler) isr() {
	fmt.Printf("Scheduler : Job called\n")
	if it.Enabled {
		// disable inactive user
		err := it.Job(context.Background())
		if err != nil {
			logger.Log.Error("error in scheduler revoke job function", zap.String("reason", err.Error()))
		}

		now := MalaysiaTime(time.Now())
		t := now.Add(10 * time.Minute)
		time.AfterFunc(t.Sub(now), it.isr)
	} else {
		it.Wg.Done()
	}
}

// Start ...
func (it *Scheduler) Start() {
	if it.Enabled {
		it.Wg.Add(1)
		it.isr()
	}
}
