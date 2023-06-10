package main

import (
	goWorker "github.com/voonik/goFramework/pkg/worker"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/utils"
	"github.com/voonik/work"
)

func initGoWorker() {
	goWorker.InitWorkerPool(10)
	jobInstance := goWorker.GetJobInstance()
	schedulerInstance := goWorker.GetSchedulerInstance()
	opts := []goWorker.AccountOption{
		goWorker.WithAccountandPortal(int64(1), int64(1)),
		goWorker.WithAccountandPortal(int64(2), int64(2)),
	}

	jobInstance.RegisterJobWithOptions(utils.ChangePendingSupplierStatus, work.JobOptions{
		MaxConcurrency: 0, Priority: 10, MaxFails: 1,
	}, helpers.ChangePendingState, opts...)
	schedulerInstance.SchedulePeriodicJob(utils.ScheduleEveryDay, utils.ChangePendingSupplierStatus, map[string]interface{}{}, opts...)

	goWorker.Start()
}
