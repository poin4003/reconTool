package initialize

import (
	"reconTool/global"
	"reconTool/internal/model"
	"reconTool/internal/worker"
)

var (
	ImportSimPool *worker.PoolManager[model.SimModel]
	SyncSimPool   *worker.PoolManager[model.SimModel]
)

func InitWorkerPools(cfg *model.ConfigModel) {
	ImportSimPool = worker.NewPoolManager(
		"ImportSimPool",
		global.ImportSimJobChan,
		&worker.ImportSimProcessor{},
	)
	ImportSimPool.Adjust(int(cfg.ImportSimConcurrency))

	SyncSimPool = worker.NewPoolManager(
		"SyncSimPool",
		global.SyncSimJobChan,
		&worker.SyncSimProcessor{},
	)
	SyncSimPool.Adjust(int(cfg.SyncSimConcurrency))
}
