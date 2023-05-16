package helpers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	worker "github.com/shopuptech/go-jobs/v3"
	"github.com/shopuptech/go-jobs/v3/backend"
	"github.com/shopuptech/go-jobs/v3/backend/rabbitmq"
	goJobsConfig "github.com/shopuptech/go-jobs/v3/config"
	"github.com/shopuptech/go-jobs/v3/handler"
	"github.com/shopuptech/go-jobs/v3/jobs"
	"github.com/shopuptech/go-jobs/v3/opts"
	"github.com/shopuptech/go-libs/logger"
	"github.com/voonik/goFramework/pkg/config"
	"github.com/voonik/goFramework/pkg/misc"
	worker2 "github.com/voonik/goFramework/pkg/worker"
	"github.com/voonik/ss2/internal/app/utils"
	"github.com/voonik/work"
)

var Wrk *worker.Worker

func InitGoJobsWorker() {
	backend, err := initRabbitMQ()
	if err != nil {
		logger.Log().Panicf("failed to initialise backend: %+v", err)
	}

	if config.StatsdConfigEnabled() {
		defer worker.SetReporterHost(config.StatsdConfigHost() + ":" + config.StatsdConfigPort())()
	}

	Wrk, err = worker.NewWorker(backend, config.AsynqConfigTeamName(), config.AsynqConfigServiceName())
	if err != nil {
		logger.Log().Panicf("failed to initialise worker: %+v", err)
	}

	RegisterJobs(Wrk)
}

func initRabbitMQ() (backend.Backend, error) {
	return rabbitmq.NewRMQBackend(config.AsynqConfigRMQConnString(),
		config.AsynqConfigPoolSize(), config.AsynqConfigTeamName(), config.AsynqConfigServiceName())
}

func RegisterJobs(wrk *worker.Worker) {
	h := handler.NewHandler(changePendingState, opts.WithPriority(10))
	wrk.RegisterJob(utils.ChangePendingSupplierStatus, h)

	vcontext := []*goJobsConfig.VContext{{VAccountID: 1, PortalID: 1}, {VAccountID: 2, PortalID: 2}}
	wrk.ScheduleJob(utils.ChangePendingSupplierStatus, vcontext, utils.ScheduleEveryDay, map[string]interface{}{})

}

func StartWorker() {
	err := Wrk.Start()
	if err != nil {
		panic(fmt.Errorf("failed to start go-jobs consumer: %w", err))
	}

	c := make(chan os.Signal, 1)                                       //nolint:revive
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) //nolint:staticcheck
	<-c
	Wrk.Stop()
}

func changePendingState(ctx context.Context, _ jobs.Job) error {
	thread := misc.ExtractThreadObjectWithDefault(ctx)
	vcontext := &worker2.VaccountContext{VaccountID: thread.VaccountId, PortalID: thread.PortalId}

	return ChangePendingState(vcontext, &work.Job{})
}
