package main

import (
	"github.com/shopuptech/go-libs/logger"
	"github.com/voonik/ss2/internal/app/helpers"
)

func main() {
	logger.Log().Infof("initialising go-jobs")
	helpers.InitGoJobsWorker()
	helpers.ScheduleJobs()
	helpers.StartWorker()
}
