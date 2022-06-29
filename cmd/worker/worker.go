package main

import (
	goWorker "github.com/voonik/goFramework/pkg/worker"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/work"
)

func main() {
	goWorker.InitWorkerPool(10)
	jobInstance := goWorker.GetJobInstance()
	schedulerInstance := goWorker.GetSchedulerInstance()
	opts := []goWorker.AccountOption{
		goWorker.WithAccountandPortal(int64(1), int64(1)),
		goWorker.WithAccountandPortal(int64(2), int64(2)),
	}

	jobInstance.RegisterJobWithOptions("change_pending_supplier_status", work.JobOptions{
		MaxConcurrency: 0, Priority: 10, MaxFails: 1,
	}, helpers.ChangePendingState, opts...)
	schedulerInstance.SchedulePeriodicJob("0 0 * * *", "change_pending_supplier_status", map[string]interface{}{}, opts...)

	goWorker.Start()
}
