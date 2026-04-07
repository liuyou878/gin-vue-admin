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

var validFirmwareStatuses = map[string]bool{
	"draft":       true,
	"testing":     true,
	"tested_pass": true,
	"stable":      true,
	"deprecated":  true,
}

// CreateFirmwareVersion 创建固件版本
func (s *FirmwareVersionService) CreateFirmwareVersion(firmware *deviceModel.FirmwareVersion) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(firmware).Error; err != nil {
			return err
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, "upload", "", firmware.Status, firmware.UploadedBy, "创建固件版本")
	})
}

// DeleteFirmwareVersion 删除固件版本
func (s *FirmwareVersionService) DeleteFirmwareVersion(id string) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.ModelFirmwareRel{}).Where("firmware_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该固件版本已被设备型号使用，不能删除")
	}
	return global.GVA_DB.Delete(&deviceModel.FirmwareVersion{}, "id = ?", id).Error
}

// DeleteFirmwareVersionByIds 批量删除固件版本
func (s *FirmwareVersionService) DeleteFirmwareVersionByIds(ids commonReq.IdsReq) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.ModelFirmwareRel{}).Where("firmware_id in ?", ids.Ids).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("所选固件版本中存在已被设备型号使用的数据，不能删除")
	}
	return global.GVA_DB.Delete(&[]deviceModel.FirmwareVersion{}, "id in ?", ids.Ids).Error
}

// UpdateFirmwareVersion 更新固件版本
func (s *FirmwareVersionService) UpdateFirmwareVersion(firmware deviceModel.FirmwareVersion) error {
	updates := map[string]interface{}{
		"version_code": firmware.VersionCode,
		"version_name": firmware.VersionName,
		"package_url":  firmware.PackageURL,
		"package_name": firmware.PackageName,
		"checksum":     firmware.Checksum,
		"status":       firmware.Status,
		"release_note": firmware.ReleaseNote,
		"test_summary": firmware.TestSummary,
		"is_latest":    firmware.IsLatest,
		"is_stable":    firmware.IsStable,
		"uploaded_by":  firmware.UploadedBy,
		"uploaded_at":  firmware.UploadedAt,
		"remark":       firmware.Remark,
	}
	return global.GVA_DB.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", firmware.ID).Updates(updates).Error
}

// GetFirmwareVersion 获取固件版本详情
func (s *FirmwareVersionService) GetFirmwareVersion(id string) (firmware deviceModel.FirmwareVersion, err error) {
	err = global.GVA_DB.Preload("ChangeItems").Preload("Tags.Tag").Where("id = ?", id).First(&firmware).Error
	return
}

// GetFirmwareVersionInfoList 获取固件版本分页列表
func (s *FirmwareVersionService) GetFirmwareVersionInfoList(info deviceReq.FirmwareVersionSearch) (list []deviceModel.FirmwareVersion, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.FirmwareVersion{})
	if info.VersionCode != "" {
		db = db.Where("version_code LIKE ?", "%"+info.VersionCode+"%")
	}
	if info.VersionName != "" {
		db = db.Where("version_name LIKE ?", "%"+info.VersionName+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.IsLatest != nil {
		db = db.Where("is_latest = ?", *info.IsLatest)
	}
	if info.IsStable != nil {
		db = db.Where("is_stable = ?", *info.IsStable)
	}
	if len(info.UploadedAtRange) == 2 {
		db = db.Where("uploaded_at BETWEEN ? AND ?", info.UploadedAtRange[0], info.UploadedAtRange[1])
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}
	err = db.Preload("Tags.Tag").Order("uploaded_at desc, id desc").Find(&list).Error
	return
}

// ChangeFirmwareVersionStatus 更新固件版本状态
func (s *FirmwareVersionService) ChangeFirmwareVersionStatus(req deviceReq.ChangeFirmwareVersionStatusRequest) error {
	if !validFirmwareStatuses[req.Status] {
		return errors.New("固件状态不合法")
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		fromStatus := firmware.Status
		updates := map[string]interface{}{
			"status":    req.Status,
			"is_stable": req.Status == "stable",
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		action := convertStatusToAction(req.Status)
		operator := req.Operator
		if operator == "" {
			operator = firmware.UploadedBy
		}
		content := req.Content
		if content == "" {
			content = "更新固件状态"
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, action, fromStatus, req.Status, operator, content)
	})
}

func convertStatusToAction(status string) string {
	switch status {
	case "testing":
		return "start_testing"
	case "tested_pass":
		return "test_pass"
	case "stable":
		return "set_stable"
	case "deprecated":
		return "deprecate"
	default:
		return "upload"
	}
}

func createFirmwareVersionLog(tx *gorm.DB, firmwareID uint, modelID *uint, action, fromStatus, toStatus, operator, content string) error {
	now := time.Now()
	log := deviceModel.FirmwareVersionLog{
		FirmwareID: firmwareID,
		ModelID:    modelID,
		Action:     action,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Operator:   operator,
		OperateAt:  &now,
		Content:    content,
	}
	return tx.Create(&log).Error
}
