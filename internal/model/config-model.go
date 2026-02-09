package model

import "gorm.io/gorm"

type ConfigModel struct {
	gorm.Model
	Code                 string `gorm:"uniqueIndex;size:50" json:"Code"`
	ImportSimConcurrency int    `gorm:"column:import_sim_concurrency" json:"ImportSimConcurrency"`
	SyncSimConcurrency   int    `gorm:"column:sync_sim_concurrency" json:"SyncSimConcurrency"`
	PartnerUrl           string `gorm:"column:api_url" json:"ApiUrl"`
	DeviceId             string `gorm:"column:device_id" json:"DeviceId"`
	Token                string `gorm:"column:token" json:"Token"`
}

func NewDefaultConfig() *ConfigModel {
	return &ConfigModel{
		Code:                 "global_config",
		ImportSimConcurrency: 5,
		SyncSimConcurrency:   10,
		PartnerUrl:           "http://localhost:8000",
		DeviceId:             "",
		Token:                "",
	}
}
