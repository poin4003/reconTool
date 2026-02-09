package model

import "gorm.io/gorm"

const (
	SimStatusNormal     = 1
	SimStatusHasPlan    = 2
	SimStatusNoPlan     = 3
	SimStatusNoBasePlan = 4
)

type SimModel struct {
	gorm.Model
	Isdn   string  `gorm:"column:isdn;size:20" json:"Isdn"`
	Serial string  `gorm:"column:serial;size:20" json:"Serial"`
	Status int64   `gorm:"column:status;default:1" json:"status" excel:"Status"`
	Note   *string `gorm:"column:note;type:text" json:"note" excel:"Note"`
}
