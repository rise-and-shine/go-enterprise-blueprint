package consumer

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/kafka"
	"golang.org/x/sync/errgroup"
)

// Config holds the configs of consumers this controller is responsible for.
type Config struct {
	// SomeConsumer kafka.ConsumerConfig `yaml:"some_consumer"`
}

type Controller struct {
	cfg              Config
	brokerConfig     kafka.BrokerConfig
	usecaseContainer *usecase.Container

	// Add your consumers here...
	someConsumer *kafka.Consumer
}

func NewController(
	cfg Config,
	brokerConfig kafka.BrokerConfig,
	usecaseContainer *usecase.Container,
) (*Controller, error) {
	ctrl := &Controller{
		cfg:              cfg,
		brokerConfig:     brokerConfig,
		usecaseContainer: usecaseContainer,
	}

	err := ctrl.initConsumers()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return ctrl, nil
}

// Start starts all the consumers in separate goroutines and
// blocks until all of them are done or one of them fails.
func (c *Controller) Start() error {
	var g errgroup.Group

	// Run your consumers here...
	// g.Go(c.someConsumer.Start)

	err := g.Wait()
	return errx.Wrap(err)
}

// Shutdown parallelly stops all consumers gracefully and
// blocks until all of them are done.
func (c *Controller) Shutdown() error {
	// errs := make(chan error, 1) // buffer size == consumer count

	// Stop your consumers here...
	// go func() { errs <- c.someConsumer.Stop() }()

	// return errx.Wrap(errors.Join(<-errs))
	return nil
}

func (c *Controller) initConsumers() error {
	// var err error

	// Add your consumers here...
	// c.someConsumer, err = kafka.NewConsumer(
	// 	c.brokerConfig,
	// 	c.cfg.SomeConsumer,
	// 	forward.ToEventSubscriber(c.usecaseContainer.SomeSubscriber()),
	// )
	// if err != nil {
	// 	return errx.Wrap(err)
	// }

	return nil
}
