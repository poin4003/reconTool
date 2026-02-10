package global

import (
	"reconTool/internal/model"
	"reconTool/pkg/logger"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Logger *logger.LoggerZap

	ImportSimJobChan = make(chan model.SimModel, 10000)
	SyncSimJobChan   = make(chan model.SimModel, 10000)
)
