package device

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"gorm.io/gorm"
)

var validModelTestResults = map[string]bool{
	"pending": true,
	"testing": true,
	"passed":  true,
	"failed":  true,
}

// CreateModelFirmwareRel 创建型号固件关系
func (s *ModelFirmwareRelService) CreateModelFirmwareRel(rel *deviceModel.ModelFirmwareRel) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := ensureDeviceModelExists(tx, rel.ModelID); err != nil {
			return err
		}
		if err := ensureFirmwareVersionExists(tx, rel.FirmwareID); err != nil {
			return err
		}
		if err := cleanupOrRejectDuplicateModelFirmwareRel(tx, rel.ModelID, rel.FirmwareID, 0); err != nil {
			return err
		}
		if err := tx.Create(rel).Error; err != nil {
			return err
		}
		modelID := rel.ModelID
		return createFirmwareVersionLog(tx, rel.FirmwareID, &modelID, "bind_model", "", rel.TestResult, rel.Tester, "绑定设备型号")
	})
}

// DeleteModelFirmwareRel 删除型号固件关系
func (s *ModelFirmwareRelService) DeleteModelFirmwareRel(id string) error {
	return global.GVA_DB.Unscoped().Delete(&deviceModel.ModelFirmwareRel{}, "id = ?", id).Error
}

// DeleteModelFirmwareRelByIds 批量删除型号固件关系
func (s *ModelFirmwareRelService) DeleteModelFirmwareRelByIds(ids commonReq.IdsReq) error {
	return global.GVA_DB.Unscoped().Delete(&[]deviceModel.ModelFirmwareRel{}, "id in ?", ids.Ids).Error
}

// UpdateModelFirmwareRel 更新型号固件关系
func (s *ModelFirmwareRelService) UpdateModelFirmwareRel(rel deviceModel.ModelFirmwareRel) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var current deviceModel.ModelFirmwareRel
		if err := tx.Where("id = ?", rel.ID).First(&current).Error; err != nil {
			return err
		}
		if err := ensureDeviceModelExists(tx, rel.ModelID); err != nil {
			return err
		}
		if err := ensureFirmwareVersionExists(tx, rel.FirmwareID); err != nil {
			return err
		}
		if err := cleanupOrRejectDuplicateModelFirmwareRel(tx, rel.ModelID, rel.FirmwareID, rel.ID); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"model_id":       rel.ModelID,
			"firmware_id":    rel.FirmwareID,
			"is_supported":   rel.IsSupported,
			"is_recommended": rel.IsRecommended,
			"test_result":    rel.TestResult,
			"tested_at":      rel.TestedAt,
			"tester":         rel.Tester,
			"remark":         rel.Remark,
		}
		return tx.Model(&deviceModel.ModelFirmwareRel{}).Where("id = ?", rel.ID).Updates(updates).Error
	})
}

// GetModelFirmwareRel 获取型号固件关系详情
func (s *ModelFirmwareRelService) GetModelFirmwareRel(id string) (rel deviceModel.ModelFirmwareRel, err error) {
	err = global.GVA_DB.Preload("Model.Category").Preload("Firmware").Where("id = ?", id).First(&rel).Error
	return
}

// GetModelFirmwareRelInfoList 获取型号固件关系分页列表
func (s *ModelFirmwareRelService) GetModelFirmwareRelInfoList(info deviceReq.ModelFirmwareRelSearch) (list []deviceModel.ModelFirmwareRel, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.ModelFirmwareRel{})
	if info.ModelID > 0 {
		db = db.Where("model_id = ?", info.ModelID)
	}
	if info.FirmwareID > 0 {
		db = db.Where("firmware_id = ?", info.FirmwareID)
	}
	if info.TestResult != "" {
		db = db.Where("test_result = ?", info.TestResult)
	}
	if info.IsRecommended != nil {
		db = db.Where("is_recommended = ?", *info.IsRecommended)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}
	err = db.Preload("Model.Category").Preload("Firmware").Order("id desc").Find(&list).Error
	return
}

// SetModelFirmwareRecommended 设置推荐版本
func (s *ModelFirmwareRelService) SetModelFirmwareRecommended(req deviceReq.SetModelFirmwareRecommendedRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var rel deviceModel.ModelFirmwareRel
		if err := tx.Where("id = ?", req.ID).First(&rel).Error; err != nil {
			return err
		}
		if err := tx.Model(&deviceModel.ModelFirmwareRel{}).Where("model_id = ?", rel.ModelID).Update("is_recommended", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&deviceModel.ModelFirmwareRel{}).Where("id = ?", rel.ID).Update("is_recommended", true).Error; err != nil {
			return err
		}
		modelID := rel.ModelID
		content := req.Content
		if content == "" {
			content = "设置推荐版本"
		}
		return createFirmwareVersionLog(tx, rel.FirmwareID, &modelID, "set_recommended", "", "", req.Operator, content)
	})
}

// SetModelFirmwareTestResult 设置型号固件测试结果
func (s *ModelFirmwareRelService) SetModelFirmwareTestResult(req deviceReq.SetModelFirmwareTestResultRequest) error {
	if !validModelTestResults[req.TestResult] {
		return errors.New("测试结果不合法")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var rel deviceModel.ModelFirmwareRel
		if err := tx.Where("id = ?", req.ID).First(&rel).Error; err != nil {
			return err
		}
		testedAt := req.TestedAt
		if testedAt == nil {
			now := time.Now()
			testedAt = &now
		}
		updates := map[string]interface{}{
			"test_result": req.TestResult,
			"tester":      req.Tester,
			"tested_at":   testedAt,
		}
		if err := tx.Model(&deviceModel.ModelFirmwareRel{}).Where("id = ?", rel.ID).Updates(updates).Error; err != nil {
			return err
		}
		modelID := rel.ModelID
		content := req.Content
		if content == "" {
			content = "更新型号固件测试结果"
		}
		return createFirmwareVersionLog(tx, rel.FirmwareID, &modelID, convertTestResultToAction(req.TestResult), rel.TestResult, req.TestResult, req.Operator, content)
	})
}

func convertTestResultToAction(testResult string) string {
	switch testResult {
	case "testing":
		return "start_testing"
	case "passed":
		return "test_pass"
	case "failed":
		return "test_fail"
	default:
		return "bind_model"
	}
}

func ensureDeviceModelExists(tx *gorm.DB, modelID uint) error {
	var count int64
	if err := tx.Model(&deviceModel.DeviceModel{}).Where("id = ?", modelID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("设备型号不存在")
	}
	return nil
}

func ensureFirmwareVersionExists(tx *gorm.DB, firmwareID uint) error {
	var count int64
	if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", firmwareID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("固件版本不存在")
	}
	return nil
}

func cleanupOrRejectDuplicateModelFirmwareRel(tx *gorm.DB, modelID, firmwareID, excludeID uint) error {
	var duplicates []deviceModel.ModelFirmwareRel
	query := tx.Unscoped().Where("model_id = ? AND firmware_id = ?", modelID, firmwareID)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Find(&duplicates).Error; err != nil {
		return err
	}

	var softDeletedIDs []uint
	for _, duplicate := range duplicates {
		if duplicate.DeletedAt.Valid {
			softDeletedIDs = append(softDeletedIDs, duplicate.ID)
			continue
		}
		return errors.New("该型号已经关联了这个固件版本，请直接编辑现有关系")
	}

	if len(softDeletedIDs) == 0 {
		return nil
	}

	return tx.Unscoped().Delete(&deviceModel.ModelFirmwareRel{}, "id in ?", softDeletedIDs).Error
}
