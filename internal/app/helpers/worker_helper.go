package helpers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	worker "github.com/shopuptech/go-jobs/v3"
	"github.com/shopuptech/go-jobs/v3/backend"
	"github.com/shopuptech/go-jobs/v3/backend/rabbitmq"
	backendRedis "github.com/shopuptech/go-jobs/v3/backend/redis"
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

var wrk *worker.Worker

func GetWorkerInstance() *worker.Worker {
	if wrk == nil {
		InitGoJobsWorker()
	}
	return wrk
}

func InitGoJobsWorker() {
	var backend backend.Backend
	var err error
	if config.AsynqConfigEnabled() {
		backend, err = initRabbitMQBackend()
	} else {
		backend, err = initRedisBackend()
	}

	if err != nil {
		logger.Log().Errorf("failed to initialise backend: %+v", err)
	}
	if config.StatsdConfigEnabled() {
		defer worker.SetReporterHost(config.StatsdConfigHost() + ":" + config.StatsdConfigPort())()
	}
	wrk, err = worker.NewWorker(backend, config.AsynqConfigTeamName(), config.AsynqConfigServiceName())
	if err != nil {
		logger.Log().Errorf("failed to initialise worker: %+v", err)
	}
	registerJobs(wrk)
}

func initRedisBackend() (backend.Backend, error) {
	redisPool := initRedis(config.JobsBackendConnString())
	return backendRedis.NewRedisBackend(backendRedis.VaccountContext{}, uint(utils.Ten), config.AsynqConfigServiceName(), redisPool)
}

func initRedis(address string) *redis.Pool {
	opts := []redis.DialOption{redis.DialDatabase(12)} //nolint:revive,gomnd
	redisPool := &redis.Pool{
		MaxActive:   10,                //nolint:revive,gomnd
		MaxIdle:     10,                //nolint:revive,gomnd
		IdleTimeout: 300 * time.Second, //nolint:revive,gomnd
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address, opts...)
		},
	}

	testRedisConn(redisPool)
	return redisPool
}

func testRedisConn(pool *redis.Pool) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		logger.Log().Errorf("Not able to connect to redis pool with connection string, check the connection")
		return
	}
	logger.Log().Infof("Connection to redis successful")
}

func initRabbitMQBackend() (backend.Backend, error) {
	return rabbitmq.NewRMQBackend(config.JobsBackendConnString(),
		config.AsynqConfigPoolSize(), config.AsynqConfigTeamName(), config.AsynqConfigServiceName())
}

func registerJobs(wrk *worker.Worker) {

}

func ScheduleJobs() {
	h := handler.NewHandler(changePendingState, opts.WithPriority(utils.Ten))
	GetWorkerInstance().RegisterJob(utils.ChangePendingSupplierStatus, h)

	vcontext := []*goJobsConfig.VContext{{VAccountID: 1, PortalID: 1}, {VAccountID: 2, PortalID: 2}}
	GetWorkerInstance().ScheduleJob(utils.ChangePendingSupplierStatus, vcontext, utils.ScheduleEveryDay, map[string]interface{}{})
}

func StartWorker() {
	err := GetWorkerInstance().Start()
	if err != nil {
		panic(fmt.Errorf("failed to start go-jobs consumer: %w", err))
	}

	c := make(chan os.Signal, 1)                                       //nolint:revive
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) //nolint:staticcheck
	<-c
	GetWorkerInstance().Stop()
}

func changePendingState(ctx context.Context, _ jobs.Job) error {
	thread := misc.ExtractThreadObjectWithDefault(ctx)
	vcontext := &worker2.VaccountContext{VaccountID: thread.VaccountId, PortalID: thread.PortalId}

	return ChangePendingState(vcontext, &work.Job{})
}
