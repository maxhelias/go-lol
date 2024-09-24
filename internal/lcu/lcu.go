package lcu

import (
	"github.com/maxhelias/golol/internal/config"
	"github.com/maxhelias/golol/internal/logger"
	"github.com/maxhelias/golol/pkg/process"
	"go.uber.org/zap/zapcore"
)

func CheckLcu(c *config.Context) {
	// Check the LCU
	newLcuInfo, err := process.FindLcuConnectInfo()
	if err != nil {
		c.Logger.Error("Error fetching LCU info", logger.LogField{Key: "error", Type: zapcore.StringType, String: err.Error()})

		return
	}

	// If the process ID is different, update the LCU info
	if c.LcuInfo == nil || newLcuInfo.ProcessID != c.LcuInfo.ProcessID {
		c.Logger.Info("LCU Info", logger.LogField{Key: "port", Type: zapcore.Int64Type, Integer: int64(newLcuInfo.Port)},
			logger.LogField{Key: "token", Type: zapcore.StringType, String: newLcuInfo.AuthToken})

		c.SetLcuInfo(newLcuInfo)
	}
}
