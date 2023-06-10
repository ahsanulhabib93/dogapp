package main

import (
	"github.com/voonik/ss2/internal/app/helpers"
)

func initGoJobs() {
	helpers.InitGoJobsWorker()
	helpers.StartWorker()
}
