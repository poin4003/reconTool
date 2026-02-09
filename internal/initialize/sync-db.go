package initialize

import (
	"log"
	"reconTool/internal/model"

	"gorm.io/gorm"
)

func SyncConfig(db *gorm.DB) {
	var cfg model.ConfigModel

	err := db.Where("code = ?", "global_config").First(&cfg).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("-> [Config] Not exists, initialize default config data...")
			defaultCfg := model.NewDefaultConfig()

			if err := db.Create(defaultCfg).Error; err != nil {
				log.Printf("!!! [Config] Error initialize: %v", err)
			} else {
				log.Println("-> [Config] Success init global_config.")
			}
		} else {
			log.Printf("!!! [Config] Query error: %v", err)
		}
		return
	}
}
