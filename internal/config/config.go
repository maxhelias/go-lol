package config

import (
	"sync"

	"github.com/maxhelias/golol/internal/api/models"
	"github.com/maxhelias/golol/internal/logger"
	"github.com/maxhelias/golol/pkg/client"
	"github.com/maxhelias/golol/pkg/process"
	"go.uber.org/zap/zapcore"
)

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

type Context struct {
	mu sync.Mutex

	Logger    *logger.ZapLogger
	LcuInfo   *process.LcuConnectInfo
	LcuClient *client.Client

	CurrSummoner models.CurrSummoner
}

func (c *Context) SetLcuInfo(lcuInfo *process.LcuConnectInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Set the LCU info
	c.LcuInfo = lcuInfo

	// Create a new LCU client
	options := []client.Option{}
	options = append(options, client.WithPort(lcuInfo.Port))
	options = append(options, client.WithAuthToken(lcuInfo.AuthToken))
	options = append(options, client.WithLcuCertPath(LcuCertPath))

	lcuClient, err := client.NewClient(options...)
	if err != nil {
		c.Logger.Error("Error creating LCU client", logger.LogField{Key: "error", Type: zapcore.StringType, String: err.Error()})
	}

	c.LcuClient = lcuClient
}

func (c *Context) SetCurrSummoner(summoner models.CurrSummoner) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CurrSummoner = summoner
}
