package service

import (
	"fmt"
	"reconTool/global"
	"reconTool/internal/model"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c              = cache.New(5*time.Minute, 10*time.Minute)
	configCacheKey = "global_config"
)

func GetGlobalConfig() (*model.ConfigModel, error) {
	if val, found := c.Get(configCacheKey); found {
		return val.(*model.ConfigModel), nil
	}

	var cfg model.ConfigModel
	if err := global.DB.Where("code = ?", "global_config").First(&cfg).Error; err != nil {
		return nil, err
	}

	c.Set(configCacheKey, &cfg, cache.DefaultExpiration)
	return &cfg, nil
}

func UpdateGlobalConfig(id string, updates map[string]interface{}) (*model.ConfigModel, error) {
	result := global.DB.Model(&model.ConfigModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("config not found")
	}

	c.Delete(configCacheKey)

	return GetGlobalConfig()
}
