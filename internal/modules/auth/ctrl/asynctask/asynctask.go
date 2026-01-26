package asynctask

import (
	"context"
	"errors"
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/pg/hooks"
	"github.com/rise-and-shine/pkg/taskmill/scheduler"
	"github.com/rise-and-shine/pkg/taskmill/worker"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Controller struct {
	worker           worker.Worker
	scheduler        scheduler.Scheduler
	usecaseContainer *usecase.Container
}

func NewController(
	dbConn *bun.DB,
	queueName string,
	usecaseContainer *usecase.Container,
) (*Controller, error) {
	worker, err := worker.New(dbConn, queueName, worker.WithPollInterval(5*time.Second))
	if err != nil {
		return nil, errx.Wrap(err)
	}

	scheduler, err := scheduler.New(dbConn, queueName)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	ctrl := &Controller{
		worker,
		scheduler,
		usecaseContainer,
	}

	ctrl.registerTasks()

	err = ctrl.registerSchedules()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return ctrl, nil
}

// Start starts taskmill worker and scheduler in separate goroutines and
// blocks until both of them are done or one of them fails.
func (c *Controller) Start() error {
	var g errgroup.Group

	ctx := context.Background()

	g.Go(func() error { return c.worker.Start(ctx) })
	logger.
		With("module", "auth").
		Info("taskmill worker is running . . .")

	g.Go(func() error { return c.scheduler.Start(ctx) })
	logger.
		With("module", "auth").
		Info("taskmill scheduler is running . . .")

	err := g.Wait()
	return errx.Wrap(err)
}

// Shutdown parallelly stops taskmill worker and scheduler gracefully and
// blocks until both of them are done.
func (c *Controller) Shutdown() error {
	errs := make(chan error, 2)

	go func() { errs <- c.worker.Stop() }()
	go func() { errs <- c.scheduler.Stop() }()

	return errx.Wrap(errors.Join(<-errs, <-errs))
}

func (c *Controller) registerTasks() {
	// Register async tasks here...
	// worker.ForwardToAsyncTask(c.worker, c.usecaseContainer.SomeAsyncTask())
}

func (c *Controller) registerSchedules() error {
	const (
		registerTimeout = 30 * time.Second
	)

	ctx, cancel := context.WithTimeout(hooks.WithSuppressedQueryLogs(context.Background()), registerTimeout)
	defer cancel()

	err := c.scheduler.RegisterSchedules(
		ctx,
		// Register cron schedules here...
		// scheduler.Schedule{
		// 	CronPattern: "* * * * *", // every minute
		// 	OperationID: c.usecaseContainer.SomeAsyncTask().OperationID(),
		// },
	)

	return errx.Wrap(err)
}
