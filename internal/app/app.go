package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxhelias/golol/internal/config"
	"github.com/maxhelias/golol/internal/lcu"
	"github.com/maxhelias/golol/internal/logger"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap/zapcore"
)

const initialRefreshInterval = 60 * time.Second

var (
	defaultOps = &opt{
		debug:           false,
		refreshInterval: initialRefreshInterval, // Default interval
	}
)

// Option instances allow to configure the library.
type Option func(a *opt) error

func WithDebug() Option {
	return func(a *opt) error {
		a.debug = true

		return nil
	}
}

func WithRefreshInterval(interval time.Duration) Option {
	return func(a *opt) error {
		a.refreshInterval = interval

		return nil
	}
}

type opt struct {
	debug           bool
	refreshInterval time.Duration
}

type App struct {
	*opt

	ctx              context.Context
	cancel           func()
	lockfile         *os.File
	currentLcuTicker *time.Ticker

	Context *config.Context
}

func New(options ...Option) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	a := &App{
		ctx:     ctx,
		cancel:  cancel,
		opt:     defaultOps,
		Context: &config.Context{},
	}

	for _, o := range options {
		if err := o(a.opt); err != nil {
			return nil, err
		}
	}

	var err error
	a.Context.Logger, err = logger.NewLogger(a.opt.debug)
	if err != nil {
		return nil, fmt.Errorf("error when creating logger: %w", err)
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		if r := recover(); r != nil {
			a.Context.Logger.Error("Application panicked", logger.LogField{Key: "panic", Type: zapcore.StringType, String: fmt.Sprintf("%v", r)})
		}
		defer a.unlock()
	}()

	if err := a.lock(); err != nil {
		a.Context.Logger.Error("Initialization error", logger.LogField{Key: "error", Type: zapcore.StringType, String: err.Error()})

		return err
	}

	return a.handle()
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}

	return nil
}

func (a *App) lock() error {
	_, err := os.Stat(config.LockfilePath)
	if err == nil {
		return fmt.Errorf("app is already running")
	}

	lockfile, err := os.OpenFile(config.LockfilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error when creating lockfile: %w", err)
	}

	a.lockfile = lockfile

	a.Context.Logger.Info("App locked")

	return nil
}

func (a *App) unlock() {
	if a.lockfile != nil {
		a.lockfile.Close()
		os.Remove(config.LockfilePath)

		a.Context.Logger.Info("Lockfile removed")
	}
}

func (a *App) handle() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	g, c := errgroup.WithContext(a.ctx)

	// Check signals
	g.Go(func() error {
		for {
			select {
			case <-c.Done():
				return c.Err()
			case sig := <-signalChan:
				a.Context.Logger.Info("Signal received", logger.LogField{Key: "signal", Type: zapcore.StringType, String: sig.String()})
				_ = a.Stop()
			}
		}
	})

	// Check the LCU
	g.Go(func() error {
		a.Context.Logger.Info("Starting LCU check")

		lcu.CheckLcu(a.Context)

		ticker := time.NewTicker(a.opt.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-c.Done():
				return c.Err()
			case <-ticker.C:
				lcu.CheckLcu(a.Context)

				// TODO
				if a.Context.LcuInfo != nil {
					var newInterval time.Duration
					if a.Context.LcuInfo != nil {
						newInterval = 5 * time.Minute
					} else {
						newInterval = initialRefreshInterval
					}

					if a.opt.refreshInterval != newInterval {
						a.opt.refreshInterval = newInterval

						ticker.Stop()
						ticker = time.NewTicker(a.opt.refreshInterval)

						a.Context.Logger.Info("Refresh interval updated", logger.LogField{Key: "interval", Type: zapcore.Int64Type, Integer: int64(a.opt.refreshInterval.Seconds())})
					}
				}
			}
		}
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
