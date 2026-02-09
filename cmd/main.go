package main

import (
	"log"
	"net/http"
	"reconTool/global"
	"reconTool/internal/api"
	"reconTool/internal/database"
	"reconTool/internal/initialize"

	"reconTool/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initialize.InitLogger()

	db := database.InitDatabase()

	if err := database.DBMigrator(db); err != nil {
		log.Fatal("Migration Failed:", err)
	}
	initialize.SyncConfig(db)

	cfg, _ := service.GetGlobalConfig()

	initialize.InitWorkerPools(cfg)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	r.LoadHTMLGlob("template/*")

	r.StaticFile("/index.js", "./template/index.js")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Recon Dashboard",
		})
	})
	v1 := r.Group("/api")

	api.InitConfigRoutes(v1)
	api.InitSimRoutes(v1)

	global.Logger.Info("Server is running on :8000")
	r.Run(":8000")
}
