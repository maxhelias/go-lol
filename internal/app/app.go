package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxhelias/golol/internal/logger"
	"github.com/maxhelias/golol/pkg/client"
	"github.com/maxhelias/golol/pkg/process"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Static variables
var (
	// Version of the application
	Version string
	// Commit of the application
	Commit string
	// BuildDate of the application
	BuildDate string

	// LockfilePath is the path to the lockfile
	LockfilePath = "go-lol.lock"
	// LcuCertPath is the path to the LCU certificate
	LcuCertPath = "config/lcu.pem"
)

// Option instances allow to configure the library.
type Option func(a *opt) error

func WithDebug() Option {
	return func(a *opt) error {
		a.debug = true

		return nil
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(a *opt) error {
		a.logger = logger

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
	logger          logger.Logger
	refreshInterval time.Duration
}

type App struct {
	*opt

	LcuInfo   *process.LcuConnectInfo
	LcuClient *client.Client

	lockfile *os.File
}

func New(options ...Option) (*App, error) {
	opt := &opt{
		refreshInterval: 60 * time.Second, // Default interval
	}

	for _, o := range options {
		if err := o(opt); err != nil {
			return nil, err
		}
	}

	if opt.logger == nil {
		var (
			l   logger.Logger
			err error
		)
		if opt.debug {
			l, err = zap.NewDevelopment()
		} else {
			l, err = zap.NewProduction()
		}

		if err != nil {
			return nil, fmt.Errorf("error when creating logger: %w", err)
		}

		opt.logger = l

	}

	return &App{opt: opt}, nil
}

func (a *App) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("Application panicked", logger.LogField{Key: "panic", Type: zapcore.StringType, String: fmt.Sprintf("%v", r)})
		}
		defer a.cleanup()
	}()

	if err := a.initialize(); err != nil {
		a.logger.Error("Initialization error", logger.LogField{Key: "error", Type: zapcore.StringType, String: err.Error()})

		return
	}

	// Channel to capture system signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Channel to stop the ticker goroutine
	tickerStop := make(chan bool)

	go a.startTicker(tickerStop)

	<-stop // Wait for a termination signal

	a.logger.Info("Shutting down application")

	tickerStop <- true // Stop the ticker goroutine
}

func (a *App) initialize() error {
	_, err := os.Stat(LockfilePath)
	if err == nil {
		return fmt.Errorf("app is already running")
	}

	lockfile, err := os.OpenFile(LockfilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error when creating lockfile: %w", err)
	}
	a.lockfile = lockfile

	a.logger.Info("App initialized")
	return nil
}

func (a *App) cleanup() {
	if a.lockfile != nil {
		a.lockfile.Close()
		os.Remove(LockfilePath)

		a.logger.Info("Lockfile removed")
	}
}

func (a *App) startTicker(stop chan bool) {
	ticker := time.NewTicker(a.refreshInterval)
	defer ticker.Stop()

	a.checkLcu()

	for {
		select {
		case <-ticker.C:
			a.checkLcu()
		case <-stop:
			return
		}
	}
}

func (a *App) checkLcu() {
	newLcuInfo, err := process.FindLcuConnectInfo()
	if err != nil {
		a.logger.Error("Error fetching LCU info", logger.LogField{Key: "error", Type: zapcore.StringType, String: err.Error()})

		return
	}

	// If the process ID is different, update the LCU info
	if a.LcuInfo == nil || newLcuInfo.ProcessID != a.LcuInfo.ProcessID {
		a.LcuInfo = newLcuInfo
		a.logger.Info("LCU Info", logger.LogField{Key: "port", Type: zapcore.Int64Type, Integer: int64(newLcuInfo.Port)},
			logger.LogField{Key: "token", Type: zapcore.StringType, String: newLcuInfo.AuthToken})

		// Create a new LCU client
		a.LcuClient = client.NewClient(newLcuInfo.Port, newLcuInfo.AuthToken, LcuCertPath)
	}
}
