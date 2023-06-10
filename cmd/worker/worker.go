package main

import (
	"github.com/shopuptech/go-libs/logger"
	"github.com/voonik/goFramework/pkg/config"
)

func main() {
	if config.AsynqConfigEnabled() {
		logger.Log().Infof("initialising go-jobs")
		initGoJobs()
	} else {
		logger.Log().Infof("initialising go-worker")
		initGoWorker()
	}
}
