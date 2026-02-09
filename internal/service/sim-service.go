package service

import (
	"fmt"
	"reconTool/global"
	"reconTool/internal/dto"
	"reconTool/internal/model"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CheckAndInsertSim(sim model.SimModel) error {
	var existing model.SimModel

	err := global.DB.Where("Isdn = ? AND Serial = ?", sim.Isdn, sim.Serial).First(&existing).Error

	if err == nil {
		return fmt.Errorf("[DUPLICATE] SIM %s already exists, skipping.", sim.Isdn)
	}

	if err := global.DB.Create(&sim).Error; err != nil {
		return fmt.Errorf("Error insert Sim %s: %v", sim.Isdn, err)
	}

	return nil
}

func applyFilters(db *gorm.DB, filter dto.SimFilter) *gorm.DB {
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.Note != "" {
		db = db.Where("note LIKE ?", "%"+filter.Note+"%")
	}
	return db
}

func GetSims(page, limit int, filter dto.SimFilter) (dto.PaginationResult[model.SimModel], error) {
	var sims []model.SimModel
	var total int64

	query := applyFilters(global.DB.Model(&model.SimModel{}), filter)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&sims).Error

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return dto.PaginationResult[model.SimModel]{
		Items: sims, Total: total, Page: page, Limit: limit, TotalPages: totalPages,
	}, err
}

func GetSimsForExport(filter dto.SimFilter) ([]model.SimModel, error) {
	var sims []model.SimModel
	err := applyFilters(global.DB.Model(&model.SimModel{}), filter).Order("id DESC").Find(&sims).Error
	return sims, err
}

func DeleteSims(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Unscoped().Where("id IN ?", ids).Delete(&model.SimModel{}).Error
}

func SyncSim(sim model.SimModel) error {
	res, err := CallAsimApi[dto.GetCurrentPackageRequest, dto.GetCurrentPackageResponse](
		"POST",
		"get-pck",
		dto.GetCurrentPackageRequest{Isdn: sim.Isdn},
	)

	if err != nil {
		global.Logger.Error("System Error khi Sync SIM", zap.String("isdn", sim.Isdn), zap.Error(err))
		return err
	}

	var updateData map[string]interface{}

	if res == nil || !res.IsSucceeded || res.Data == nil {
		msg := "Partner Error"
		if res != nil {
			msg = res.ErrorMessage
		}

		updateData = map[string]interface{}{
			"status": 4,
			"note":   msg,
		}
	} else {
		packageString := extractUniquePackageNames(res.Data)

		if packageString == "" {
			updateData = map[string]interface{}{
				"status": model.SimStatusNoPlan,
				"note":   "SIM không có gói",
			}
		} else {
			updateData = map[string]interface{}{
				"status": model.SimStatusHasPlan,
				"note":   packageString,
			}
		}
	}

	err = global.DB.Model(&sim).Updates(updateData).Error
	if err != nil {
		global.Logger.Error("DB Update Failed", zap.String("isdn", sim.Isdn), zap.Error(err))
		return err
	}

	return nil
}

func extractUniquePackageNames(data *dto.GetCurrentPackageResponse) string {
	if data == nil {
		return ""
	}

	nameSet := make(map[string]struct{})

	allLists := [][]dto.PackageForList{
		data.MainPackages,
		data.AddonPackages,
		data.IpPackages,
		data.OtherPackages,
	}

	for _, list := range allLists {
		for _, pkg := range list {
			if pkg.GoodName != "" {
				nameSet[pkg.GoodName] = struct{}{}
			}
		}
	}

	if len(nameSet) == 0 {
		return ""
	}

	var uniqueNames []string
	for name := range nameSet {
		uniqueNames = append(uniqueNames, name)
	}

	return strings.Join(uniqueNames, "|")
}
