package utility

import (
	"context"
	"fmt"
	"galasejahtera/pkg/logger"
	"sync"
	"time"

	"go.uber.org/zap"
)

type DailyScheduler struct {
	Enabled   bool
	HaveStart bool
	Job       func(ctx context.Context) error
	Wg        sync.WaitGroup
}

func (it *DailyScheduler) isr() {
	fmt.Printf("Scheduler : Job called\n")
	if it.Enabled {
		now := MalaysiaTime(time.Now())
		ctx := context.Background()
		if it.HaveStart {
			// generate report
			err := it.Job(ctx)
			if err != nil {
				logger.Log.Error("error in scheduler job function", zap.String("reason", err.Error()))
			}

			ctx.Done()
		} else {
			it.HaveStart = true
			// generate report
			err := it.Job(ctx)
			if err != nil {
				logger.Log.Error("error in scheduler job function", zap.String("reason", err.Error()))
			}
		}

		t := now.Add(6 * time.Hour)
		fmt.Printf("\tJob interval %v\n", t.Sub(now))
		time.AfterFunc(t.Sub(now), it.isr)
	} else {
		it.Wg.Done()
	}
}

//trigger
func (it *DailyScheduler) Start() {
	if it.Enabled {
		it.Wg.Add(1)
		it.isr()
	}
}
