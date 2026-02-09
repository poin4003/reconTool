package api

import (
	"fmt"
	"net/http"
	"reconTool/global"
	"reconTool/internal/initialize"
	"reconTool/internal/model"
	"reconTool/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func InitConfigRoutes(rg *gin.RouterGroup) {
	config := rg.Group("/config")
	{
		config.GET("", GetConfigList)
		config.PATCH("/:id", PatchConfig)
		config.GET("/queue", GetQueueStatus)
	}
}

func GetConfigList(c *gin.Context) {
	var configs []model.ConfigModel

	if err := global.DB.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func PatchConfig(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(updates, "id")
	delete(updates, "code")

	updatedCfg, err := service.UpdateGlobalConfig(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if initialize.ImportSimPool != nil {
		initialize.ImportSimPool.Adjust(updatedCfg.ImportSimConcurrency)
	}
	if initialize.SyncSimPool != nil {
		initialize.SyncSimPool.Adjust(updatedCfg.SyncSimConcurrency)
	}

	c.Status(http.StatusNoContent)
}

func GetQueueStatus(c *gin.Context) {
	importPending := len(global.ImportSimJobChan)
	importCap := cap(global.ImportSimJobChan)
	importWorkers := 0
	if initialize.ImportSimPool != nil {
		importWorkers = initialize.ImportSimPool.GetWorkerCount()
	}

	syncPending := len(global.SyncSimJobChan)
	syncCap := cap(global.SyncSimJobChan)
	syncWorkers := 0
	if initialize.SyncSimPool != nil {
		syncWorkers = initialize.SyncSimPool.GetWorkerCount()
	}

	c.JSON(http.StatusOK, gin.H{
		"import_queue": gin.H{
			"pending":        importPending,
			"capacity":       importCap,
			"active_workers": importWorkers,
			"usage_percent":  fmt.Sprintf("%.2f%%", float64(importPending)/float64(importCap)*100),
		},
		"sync_queue": gin.H{
			"pending":        syncPending,
			"capacity":       syncCap,
			"active_workers": syncWorkers,
			"usage_percent":  fmt.Sprintf("%.2f%%", float64(syncPending)/float64(syncCap)*100),
		},
		"system_time": time.Now().Format("15:04:05"),
	})
}
