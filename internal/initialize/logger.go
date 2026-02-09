package initialize

import (
	"reconTool/global"
	"reconTool/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger()
}
