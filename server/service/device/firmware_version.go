package device

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	deviceResp "github.com/flipped-aurora/gin-vue-admin/server/model/device/response"
	exampleModel "github.com/flipped-aurora/gin-vue-admin/server/model/example"
	emailUtils "github.com/flipped-aurora/gin-vue-admin/server/plugin/email/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var validFirmwareStatuses = map[string]bool{
	"pending_test":    true,
	"testing":         true,
	"tested_pass":     true,
	"test_failed":     true,
	"pending_release": true,
}

var validPublishStatuses = map[string]bool{
	"unpublished": true,
	"published":   true,
	"voided":      true,
	"removed":     true,
}

var firmwareVersionCodePattern = regexp.MustCompile(`^\d+(\.\d+){2}$`)

func validateFirmwareVersionCode(versionCode string) error {
	versionCode = strings.TrimSpace(versionCode)
	if versionCode == "" {
		return errors.New("请填写版本号")
	}
	return nil
}

func compareFirmwareVersionCode(a, b string) (int, error) {
	partsA := strings.Split(strings.TrimSpace(a), ".")
	partsB := strings.Split(strings.TrimSpace(b), ".")
	if len(partsA) != 3 || len(partsB) != 3 {
		return 0, errors.New("版本号格式不合法")
	}
	for i := 0; i < 3; i++ {
		left, err := strconv.ParseUint(partsA[i], 10, 64)
		if err != nil {
			return 0, errors.New("版本号格式不合法")
		}
		right, err := strconv.ParseUint(partsB[i], 10, 64)
		if err != nil {
			return 0, errors.New("版本号格式不合法")
		}
		switch {
		case left > right:
			return 1, nil
		case left < right:
			return -1, nil
		}
	}
	return 0, nil
}

type modelLatestPublishedVersionInfo struct {
	ModelID     uint   `gorm:"column:model_id"`
	ModelName   string `gorm:"column:model_name"`
	VersionCode string `gorm:"column:version_code"`
}

func ensureFirmwareVersionNotLowerThanPublished(tx *gorm.DB, versionCode string, modelIDs []uint) error {
	if len(modelIDs) == 0 {
		return nil
	}
	if !firmwareVersionCodePattern.MatchString(strings.TrimSpace(versionCode)) {
		return nil
	}
	var rows []modelLatestPublishedVersionInfo
	err := tx.Table("alpha_model_firmware_rels").
		Select("alpha_model_firmware_rels.model_id, alpha_device_models.model_name, alpha_firmware_versions.version_code").
		Joins("JOIN alpha_firmware_versions ON alpha_firmware_versions.id = alpha_model_firmware_rels.firmware_id").
		Joins("JOIN alpha_device_models ON alpha_device_models.id = alpha_model_firmware_rels.model_id").
		Where("alpha_model_firmware_rels.model_id IN ? AND alpha_firmware_versions.publish_status = ?", modelIDs, "published").
		Find(&rows).Error
	if err != nil {
		return err
	}
	latestByModel := make(map[uint]modelLatestPublishedVersionInfo, len(rows))
	for _, row := range rows {
		if !firmwareVersionCodePattern.MatchString(strings.TrimSpace(row.VersionCode)) {
			continue
		}
		current, exists := latestByModel[row.ModelID]
		if !exists {
			latestByModel[row.ModelID] = row
			continue
		}
		cmp, err := compareFirmwareVersionCode(row.VersionCode, current.VersionCode)
		if err != nil {
			return err
		}
		if cmp > 0 {
			latestByModel[row.ModelID] = row
		}
	}
	for _, modelID := range modelIDs {
		latest, exists := latestByModel[modelID]
		if !exists {
			continue
		}
		cmp, err := compareFirmwareVersionCode(versionCode, latest.VersionCode)
		if err != nil {
			return err
		}
		if cmp < 0 {
			modelName := strings.TrimSpace(latest.ModelName)
			if modelName != "" {
				return fmt.Errorf("版本号不能小于当前正式版最新版本(%s: %s)", modelName, latest.VersionCode)
			}
			return fmt.Errorf("版本号不能小于当前正式版最新版本(%s)", latest.VersionCode)
		}
	}
	return nil
}

func uniqueUintIDs(ids []uint) []uint {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

// CreateFirmwareVersion 创建固件版本
func (s *FirmwareVersionService) CreateFirmwareVersion(firmware *deviceModel.FirmwareVersion) error {
	if err := validateFirmwareVersionCode(firmware.VersionCode); err != nil {
		return err
	}
	modelIDs := uniqueUintIDs(firmware.ModelIDs)
	if len(modelIDs) == 0 {
		return errors.New("请选择设备型号")
	}
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if firmware.Status == "" || firmware.Status == "draft" {
			firmware.Status = "pending_test"
		}
		if !validFirmwareStatuses[firmware.Status] {
			return errors.New("开发状态不合法")
		}
		if firmware.PublishStatus == "" {
			firmware.PublishStatus = "unpublished"
		}
		if err := ensureFirmwareVersionNotLowerThanPublished(tx, firmware.VersionCode, modelIDs); err != nil {
			return err
		}
		if err := tx.Create(firmware).Error; err != nil {
			return err
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, "upload", "", firmware.Status, firmware.UploadedBy, "创建固件版本", firmware); err != nil {
			return err
		}
		for _, modelID := range modelIDs {
			if err := ensureDeviceModelExists(tx, modelID); err != nil {
				return err
			}
			if err := cleanupOrRejectDuplicateModelFirmwareRel(tx, modelID, firmware.ID, 0); err != nil {
				return err
			}
			if err := tx.Create(&deviceModel.ModelFirmwareRel{
				ModelID:       modelID,
				FirmwareID:    firmware.ID,
				IsSupported:   true,
				IsRecommended: false,
			}).Error; err != nil {
				return err
			}
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   "upload",
			Status:   firmware.Status,
			Operator: firmware.UploadedBy,
			Content:  firmware.ReleaseNote,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, firmware.NotifyTo); sendErr != nil {
			return fmt.Errorf("固件已新增，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// DeleteFirmwareVersion 删除固件版本
func (s *FirmwareVersionService) DeleteFirmwareVersion(id string) error {
	var firmware deviceModel.FirmwareVersion
	if err := global.GVA_DB.Select("id", "publish_status").Where("id = ?", id).First(&firmware).Error; err != nil {
		return err
	}
	if firmware.PublishStatus == "published" || firmware.PublishStatus == "voided" {
		return errors.New("已发布的版本不能删除，只能下架并保留历史")
	}
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
	var publishedCount int64
	if err := global.GVA_DB.Model(&deviceModel.FirmwareVersion{}).Where("id in ? AND publish_status in ?", ids.Ids, []string{"published", "voided"}).Count(&publishedCount).Error; err != nil {
		return err
	}
	if publishedCount > 0 {
		return errors.New("所选固件版本中存在已发布或已下架的数据，不能删除")
	}
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
	if err := validateFirmwareVersionCode(firmware.VersionCode); err != nil {
		return err
	}
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var current deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", firmware.ID).First(&current).Error; err != nil {
			return err
		}
		packageChanged := current.PackageURL != firmware.PackageURL ||
			current.PackageName != firmware.PackageName ||
			current.PackageFileID != firmware.PackageFileID ||
			current.Checksum != firmware.Checksum
		if packageChanged && (current.PublishStatus == "published" || current.PublishStatus == "voided") {
			return errors.New("已发布版本不能再更新安装包")
		}
		if packageChanged && !canReplaceFirmwarePackage(current.Status, current.PublishStatus) {
			return errors.New("当前阶段不允许替换安装包")
		}

		nextStatus := normalizeFirmwareStatus(current.Status)
		if packageChanged && current.Status == "test_failed" {
			nextStatus = "pending_test"
		}

		updates := map[string]interface{}{
			"version_code":    firmware.VersionCode,
			"version_name":    firmware.VersionName,
			"package_url":     firmware.PackageURL,
			"package_name":    firmware.PackageName,
			"package_file_id": firmware.PackageFileID,
			"checksum":        firmware.Checksum,
			"status":          nextStatus,
			"release_note":    firmware.ReleaseNote,
			"test_summary":    firmware.TestSummary,
			"uploaded_by":     firmware.UploadedBy,
			"uploaded_at":     firmware.UploadedAt,
			"remark":          firmware.Remark,
			"log_file_id":     firmware.LogFileID,
			"log_file_name":   firmware.LogFileName,
			"log_file_size":   firmware.LogFileSize,
			"log_uploaded_by": firmware.LogUploadedBy,
			"log_uploaded_at": firmware.LogUploadedAt,
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", firmware.ID).Updates(updates).Error; err != nil {
			return err
		}
		if packageChanged {
			content := strings.TrimSpace(firmware.OperationContent)
			if content == "" {
				content = "测试包更新"
			}
			if err := createFirmwareVersionLog(
				tx,
				current.ID,
				nil,
				"fix_upload",
				current.Status,
				nextStatus,
				firmware.UploadedBy,
				content,
			); err != nil {
				return err
			}
			mailPayload = buildFirmwareActionMailPayload(tx, current.ID, firmwareActionMailOptions{
				Action:   "fix_upload",
				Status:   nextStatus,
				Operator: firmware.UploadedBy,
				Content:  content,
			})
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, firmware.NotifyTo); sendErr != nil {
			return fmt.Errorf("更新包已完成，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
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
		db = db.Where("status = ?", normalizeFirmwareStatus(info.Status))
	}
	if info.PublishStatus != "" {
		db = db.Where("publish_status = ?", normalizePublishStatus(info.PublishStatus))
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
	req.Status = normalizeFirmwareStatus(req.Status)
	if !validFirmwareStatuses[req.Status] {
		return errors.New("开发状态不合法")
	}
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus == "published" || firmware.PublishStatus == "voided" {
			return errors.New("已进入发布域的版本不能再直接修改开发状态")
		}
		if !isAllowedFirmwareStatusTransition(firmware.Status, req.Status) {
			return errors.New("当前开发状态不允许执行该操作")
		}
		fromStatus := firmware.Status
		updates := map[string]interface{}{
			"status": req.Status,
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		action := convertStatusToAction(fromStatus, req.Status)
		operator := req.Operator
		if operator == "" {
			operator = firmware.UploadedBy
		}
		content := req.Content
		if content == "" {
			content = "更新固件状态"
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, action, fromStatus, req.Status, operator, content); err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   action,
			Status:   req.Status,
			Operator: operator,
			Content:  content,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, req.NotifyTo); sendErr != nil {
			return fmt.Errorf("状态已更新，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// PublishFirmwareVersion 发布固件版本
func (s *FirmwareVersionService) PublishFirmwareVersion(req deviceReq.PublishFirmwareVersionRequest) error {
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus == "voided" {
			return errors.New("已下架版本不能再次发布")
		}
		if firmware.PublishStatus == "published" {
			return errors.New("该版本已经发布")
		}
		switch normalizeFirmwareStatus(firmware.Status) {
		case "tested_pass", "pending_release":
		default:
			return errors.New("测试通过后才能发布")
		}
		now := time.Now()
		if err := tx.Model(&deviceModel.FirmwareVersion{}).
			Where("is_latest = ? AND publish_status = ?", true, "published").
			Updates(map[string]interface{}{"is_latest": false}).Error; err != nil {
			return err
		}
		operator := req.Operator
		if operator == "" {
			operator = firmware.UploadedBy
		}
		updates := map[string]interface{}{
			"publish_status": "published",
			"is_latest":      true,
			"published_by":   operator,
			"published_at":   &now,
		}
		releaseNote := strings.TrimSpace(req.ReleaseNote)
		if releaseNote == "" {
			return errors.New("请填写版本说明")
		}
		updates["release_note"] = releaseNote
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		content := req.Content
		if content == "" {
			content = "发布版本"
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, "publish", firmware.Status, "published", operator, content); err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   "publish",
			Status:   "published",
			Operator: operator,
			Content:  content,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, req.NotifyTo); sendErr != nil {
			return fmt.Errorf("发布已完成，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// SetFirmwareStable 设置稳定版本标记
func (s *FirmwareVersionService) SetFirmwareStable(req deviceReq.SetFirmwareStableRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus != "published" {
			return errors.New("只有已发布版本才能设置稳定版本")
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Update("is_stable", req.Stable).Error; err != nil {
			return err
		}
		operator := req.Operator
		if operator == "" {
			operator = firmware.PublishedBy
		}
		content := req.Content
		action := "mark_stable"
		if req.Stable {
			if content == "" {
				content = "标记为稳定版本"
			}
		} else {
			action = "unmark_stable"
			if content == "" {
				content = "取消稳定版本标记"
			}
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, action, firmware.PublishStatus, firmware.PublishStatus, operator, content)
	})
}

// VoidFirmwareVersion 下架已发布固件版本
func (s *FirmwareVersionService) VoidFirmwareVersion(req deviceReq.VoidFirmwareVersionRequest) error {
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus != "published" {
			return errors.New("只有已发布版本才能下架")
		}
		now := time.Now()
		operator := req.Operator
		if operator == "" {
			operator = firmware.PublishedBy
		}
		updates := map[string]interface{}{
			"publish_status": "voided",
			"is_latest":      false,
			"is_stable":      false,
			"voided_by":      operator,
			"voided_at":      &now,
			"void_reason":    req.VoidReason,
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Model(&deviceModel.ModelFirmwareRel{}).Where("firmware_id = ?", req.ID).Update("is_recommended", false).Error; err != nil {
			return err
		}
		content := req.Content
		if content == "" {
			content = req.VoidReason
		}
		if content == "" {
			content = "下架已发布版本"
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, "void_release", firmware.PublishStatus, "voided", operator, content); err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   "void_release",
			Status:   "voided",
			Operator: operator,
			Content:  content,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, req.NotifyTo); sendErr != nil {
			return fmt.Errorf("下架已完成，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// OnShelfFirmwareVersion 上架已下架固件版本
func (s *FirmwareVersionService) OnShelfFirmwareVersion(req deviceReq.OnShelfFirmwareVersionRequest) error {
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus != "voided" {
			return errors.New("只有已下架版本才能上架")
		}
		now := time.Now()
		if err := tx.Model(&deviceModel.FirmwareVersion{}).
			Where("is_latest = ? AND publish_status = ?", true, "published").
			Updates(map[string]interface{}{"is_latest": false}).Error; err != nil {
			return err
		}
		operator := req.Operator
		if operator == "" {
			operator = firmware.PublishedBy
		}
		updates := map[string]interface{}{
			"publish_status": "published",
			"is_latest":      true,
			"published_by":   operator,
			"published_at":   &now,
			"voided_by":      "",
			"voided_at":      nil,
			"void_reason":    "",
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		content := req.Content
		if content == "" {
			content = "上架已下架版本"
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, "on_shelf_release", "voided", "published", operator, content); err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   "on_shelf_release",
			Status:   "published",
			Operator: operator,
			Content:  content,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, req.NotifyTo); sendErr != nil {
			return fmt.Errorf("上架已完成，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// RemoveFirmwareVersion 移除已下架固件版本
func (s *FirmwareVersionService) RemoveFirmwareVersion(req deviceReq.RemoveFirmwareVersionRequest) error {
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus != "voided" {
			return errors.New("只有已下架版本才能移除")
		}
		operator := req.Operator
		if operator == "" {
			operator = firmware.VoidedBy
		}
		updates := map[string]interface{}{
			"publish_status": "removed",
			"is_latest":      false,
			"is_stable":      false,
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		content := req.Content
		if content == "" {
			content = "移除已下架版本"
		}
		if err := createFirmwareVersionLog(tx, firmware.ID, nil, "remove_release", "voided", "removed", operator, content); err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, firmware.ID, firmwareActionMailOptions{
			Action:   "remove_release",
			Status:   "removed",
			Operator: operator,
			Content:  content,
		})
		return nil
	})
	if txErr != nil {
		return txErr
	}
	if mailPayload != nil {
		if sendErr := sendFirmwareActionEmail(mailPayload, req.NotifyTo); sendErr != nil {
			return fmt.Errorf("移除已完成，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// DeleteFirmwarePackage 删除固件包
func (s *FirmwareVersionService) DeleteFirmwarePackage(req deviceReq.DeleteFirmwarePackageRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", req.ID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.ID > 0 {
			return errors.New("已创建版本不能删除安装包，只能重新上传替换")
		}
		fileRecord, err := resolveFirmwarePackageFile(tx, firmware)
		if err != nil {
			return err
		}
		if fileRecord.ID == 0 {
			return errors.New("未找到可删除的安装包")
		}
		if err := upload.NewOss().DeleteFile(fileRecord.Key); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"package_url":     "",
			"package_name":    "",
			"package_file_id": 0,
			"checksum":        "",
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		operator := req.Operator
		if operator == "" {
			operator = firmware.UploadedBy
		}
		content := req.Content
		if content == "" {
			content = "删除安装包"
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, "delete_package", firmware.Status, firmware.Status, operator, content, &firmware)
	})
}

type firmwareActionMailPayload struct {
	Subject string
	Body    string
}

type firmwareActionMailOptions struct {
	Action   string
	Status   string
	Operator string
	Content  string
}

func buildFirmwareActionMailPayload(tx *gorm.DB, firmwareID uint, opts firmwareActionMailOptions) *firmwareActionMailPayload {
	var firmware deviceModel.FirmwareVersion
	if err := tx.Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return nil
	}

	var rels []deviceModel.ModelFirmwareRel
	_ = tx.Preload("Model.Category").Where("firmware_id = ?", firmwareID).Find(&rels).Error

	modelNames := make([]string, 0, len(rels))
	categoryNames := make([]string, 0, len(rels))
	modelSeen := map[string]struct{}{}
	categorySeen := map[string]struct{}{}
	for _, rel := range rels {
		modelName := strings.TrimSpace(rel.Model.ModelName)
		categoryName := strings.TrimSpace(rel.Model.Category.Name)
		if modelName != "" {
			if _, ok := modelSeen[modelName]; !ok {
				modelSeen[modelName] = struct{}{}
				modelNames = append(modelNames, modelName)
			}
		}
		if categoryName != "" {
			if _, ok := categorySeen[categoryName]; !ok {
				categorySeen[categoryName] = struct{}{}
				categoryNames = append(categoryNames, categoryName)
			}
		}
	}

	versionLabel := strings.TrimSpace(firmware.VersionCode)
	if versionLabel == "" {
		versionLabel = strings.TrimSpace(firmware.VersionName)
	}
	if versionLabel == "" {
		versionLabel = "未知版本"
	}
	modelLabel := "-"
	if len(modelNames) > 0 {
		modelLabel = strings.Join(modelNames, "、")
	}
	categoryLabel := "-"
	if len(categoryNames) > 0 {
		categoryLabel = strings.Join(categoryNames, "、")
	}
	actionLabel := firmwareActionLabel(opts.Action)
	statusLabel := firmwareActionStatusLabel(opts.Status)
	subject := fmt.Sprintf("固件流程通知 - %s - %s - %s", modelLabel, versionLabel, actionLabel)

	escapeText := func(value string) string {
		value = strings.TrimSpace(value)
		if value == "" {
			return "-"
		}
		return strings.ReplaceAll(html.EscapeString(value), "\n", "<br/>")
	}

	var builder strings.Builder
	builder.WriteString(`<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; line-height: 1.7; color: #1f2937;">`)
	builder.WriteString(`<h2 style="margin: 0 0 16px; font-size: 20px;">固件流程通知</h2>`)
	builder.WriteString(`<table style="border-collapse: collapse; width: 100%; max-width: 760px;">`)
	writeRow := func(label, value string) {
		builder.WriteString(`<tr>`)
		builder.WriteString(`<td style="padding: 8px 12px; border: 1px solid #e5e7eb; background: #f9fafb; width: 140px; font-weight: 600;">`)
		builder.WriteString(escapeText(label))
		builder.WriteString(`</td>`)
		builder.WriteString(`<td style="padding: 8px 12px; border: 1px solid #e5e7eb;">`)
		builder.WriteString(value)
		builder.WriteString(`</td>`)
		builder.WriteString(`</tr>`)
	}
	writeRow("设备类别", escapeText(categoryLabel))
	writeRow("设备型号", escapeText(modelLabel))
	writeRow("版本号", escapeText(firmware.VersionCode))
	writeRow("版本名称", escapeText(firmware.VersionName))
	writeRow("当前动作", escapeText(actionLabel))
	writeRow("当前状态", escapeText(statusLabel))
	writeRow("操作人", escapeText(opts.Operator))
	writeRow("说明", escapeText(opts.Content))
	builder.WriteString(`</table>`)
	if opts.Action == "publish" && strings.TrimSpace(global.GVA_CONFIG.System.WebURL) != "" {
		webURL := strings.TrimSpace(global.GVA_CONFIG.System.WebURL)
		downloadURL := webURL + "/#/publicFirmwareDownload"
		if len(rels) > 0 {
			categoryID := rels[0].Model.CategoryID
			modelID := rels[0].ModelID
			downloadURL += fmt.Sprintf("?categoryId=%d&modelId=%d", categoryID, modelID)
		}
		builder.WriteString(`<div style="margin-top: 20px; padding: 14px 18px; border-radius: 12px; background: #eff6ff; border: 1px solid #bfdbfe;">`)
		builder.WriteString(`<div style="font-size: 14px; font-weight: 600; color: #1e40af; margin-bottom: 6px;">公开下载页</div>`)
		builder.WriteString(`<div style="font-size: 13px; color: #475569; margin-bottom: 8px;">可以前往公开下载页下载该版本：</div>`)
		builder.WriteString(`<a href="` + downloadURL + `" style="color: #2563eb; font-size: 13px; word-break: break-all;">` + downloadURL + `</a>`)
		builder.WriteString(`</div>`)
	}
	builder.WriteString(`<div style="margin-top: 16px; color: #6b7280; font-size: 12px;">该邮件由固件流程管理自动发送</div>`)
	builder.WriteString(`</div>`)

	return &firmwareActionMailPayload{
		Subject: subject,
		Body:    builder.String(),
	}
}

func sendFirmwareActionEmail(payload *firmwareActionMailPayload, notifyTo string) error {
	if payload == nil {
		return nil
	}
	recipient := normalizeFirmwareActionRecipients(notifyTo)
	if recipient == "" {
		return nil
	}
	return emailUtils.Email(recipient, payload.Subject, payload.Body)
}

func normalizeFirmwareActionRecipients(raw string) string {
	parts := strings.Split(raw, ",")
	recipients := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		recipients = append(recipients, part)
	}
	return strings.Join(recipients, ",")
}

func firmwareActionLabel(action string) string {
	switch action {
	case "upload":
		return "新增固件"
	case "fix_upload":
		return "更新固件包"
	case "start_testing":
		return "开始测试"
	case "reject_release":
		return "驳回到测试中"
	case "publish":
		return "发布版本"
	case "void_release":
		return "下架版本"
	case "on_shelf_release":
		return "上架版本"
	default:
		return "流程通知"
	}
}

func firmwareActionStatusLabel(status string) string {
	switch normalizeFirmwareStatus(status) {
	case "pending_test":
		return "待测试"
	case "testing":
		return "测试中"
	case "tested_pass":
		return "测试通过"
	case "test_failed":
		return "测试未通过"
	case "pending_release":
		return "待发布"
	default:
		if status == "published" {
			return "已发布"
		}
		return status
	}
}

func convertStatusToAction(fromStatus, toStatus string) string {
	switch toStatus {
	case "pending_test":
		if normalizeFirmwareStatus(fromStatus) == "test_failed" {
			return "fix_upload"
		}
		return "upload"
	case "testing":
		if normalizeFirmwareStatus(fromStatus) == "pending_release" ||
			normalizeFirmwareStatus(fromStatus) == "tested_pass" {
			return "reject_release"
		}
		return "start_testing"
	case "tested_pass":
		return "test_pass"
	case "test_failed":
		return "test_fail"
	case "pending_release":
		return "submit_release"
	default:
		return "upload"
	}
}

func canReplaceFirmwarePackage(status, publishStatus string) bool {
	if normalizePublishStatus(publishStatus) != "unpublished" {
		return false
	}
	switch normalizeFirmwareStatus(status) {
	case "pending_test", "test_failed":
		return true
	default:
		return false
	}
}

func normalizeFirmwareStatus(status string) string {
	switch status {
	case "", "draft":
		return "pending_test"
	case "failed":
		return "test_failed"
	default:
		return status
	}
}

func normalizePublishStatus(status string) string {
	if status == "" {
		return "unpublished"
	}
	return status
}

func isAllowedFirmwareStatusTransition(from, to string) bool {
	from = normalizeFirmwareStatus(from)
	to = normalizeFirmwareStatus(to)
	if from == to {
		return true
	}
	switch from {
	case "pending_test":
		return to == "testing"
	case "testing":
		return to == "tested_pass" || to == "test_failed"
	case "tested_pass":
		return to == "testing"
	case "test_failed":
		return to == "testing" || to == "pending_test"
	case "pending_release":
		return to == "testing"
	default:
		return false
	}
}

func createFirmwareVersionLog(tx *gorm.DB, firmwareID uint, modelID *uint, action, fromStatus, toStatus, operator, content string, snapshot ...*deviceModel.FirmwareVersion) error {
	pkgURL, pkgName, pkgFileID, pkgSize, checksum := resolveFirmwareLogPackageSnapshot(tx, firmwareID, snapshot...)
	now := time.Now()
	log := deviceModel.FirmwareVersionLog{
		FirmwareID:    firmwareID,
		ModelID:       modelID,
		Action:        action,
		FromStatus:    fromStatus,
		ToStatus:      toStatus,
		Operator:      operator,
		OperateAt:     &now,
		Content:       content,
		PackageURL:    pkgURL,
		PackageName:   pkgName,
		PackageFileID: pkgFileID,
		PackageSize:   pkgSize,
		Checksum:      checksum,
	}
	return tx.Create(&log).Error
}

func resolveFirmwareLogPackageSnapshot(tx *gorm.DB, firmwareID uint, snapshot ...*deviceModel.FirmwareVersion) (string, string, uint, int64, string) {
	if len(snapshot) > 0 && snapshot[0] != nil {
		firmware := snapshot[0]
		fileRecord, _ := resolveFirmwarePackageFile(tx, *firmware)
		return firmware.PackageURL, firmware.PackageName, firmware.PackageFileID, fileRecord.Size, firmware.Checksum
	}

	var firmware deviceModel.FirmwareVersion
	if err := tx.Select("package_url", "package_name", "package_file_id", "checksum").Where("id = ?", firmwareID).First(&firmware).Error; err == nil {
		fileRecord, _ := resolveFirmwarePackageFile(tx, firmware)
		return firmware.PackageURL, firmware.PackageName, firmware.PackageFileID, fileRecord.Size, firmware.Checksum
	}

	return "", "", 0, 0, ""
}

func resolveFirmwarePackageFile(tx *gorm.DB, firmware deviceModel.FirmwareVersion) (exampleModel.ExaFileUploadAndDownload, error) {
	if firmware.PackageFileID > 0 {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("id = ?", firmware.PackageFileID).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	if firmware.PackageURL != "" {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("url = ?", firmware.PackageURL).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	return exampleModel.ExaFileUploadAndDownload{}, errors.New("未找到安装包文件记录")
}

func resolveFirmwareLogPackageFile(tx *gorm.DB, log deviceModel.FirmwareVersionLog) (exampleModel.ExaFileUploadAndDownload, error) {
	if log.PackageFileID > 0 {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("id = ?", log.PackageFileID).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	if log.PackageURL != "" {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("url = ?", log.PackageURL).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	return exampleModel.ExaFileUploadAndDownload{}, errors.New("未找到安装包文件记录")
}

type firmwarePackageDownload struct {
	Reader      io.ReadCloser
	Size        int64
	ContentType string
	FileName    string
}

func (d *firmwarePackageDownload) Close() error {
	if d == nil || d.Reader == nil {
		return nil
	}
	return d.Reader.Close()
}

func (s *FirmwareVersionService) OpenFirmwarePackageDownload(firmwareID uint) (*firmwarePackageDownload, error) {
	var firmware deviceModel.FirmwareVersion
	if err := global.GVA_DB.Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return nil, err
	}
	return openFirmwarePackageDownload(global.GVA_DB, firmware)
}

func (s *FirmwareVersionService) OpenPublicFirmwarePackageDownload(firmwareID uint) (*firmwarePackageDownload, error) {
	var firmware deviceModel.FirmwareVersion
	if err := global.GVA_DB.Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return nil, err
	}
	if firmware.PublishStatus != "published" {
		return nil, errors.New("只有已发布固件支持公开下载")
	}
	return openFirmwarePackageDownload(global.GVA_DB, firmware)
}

func (s *FirmwareVersionService) OpenFirmwareLogPackageDownload(logID uint) (*firmwarePackageDownload, error) {
	var log deviceModel.FirmwareVersionLog
	if err := global.GVA_DB.Where("id = ?", logID).First(&log).Error; err != nil {
		return nil, err
	}
	return openFirmwareLogPackageDownload(global.GVA_DB, log)
}

func (s *FirmwareVersionService) OpenDeveloperLogDownload(firmwareID uint) (*firmwarePackageDownload, error) {
	var firmware deviceModel.FirmwareVersion
	if err := global.GVA_DB.Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return nil, err
	}
	return openDeveloperLogDownload(global.GVA_DB, firmware)
}

func (s *FirmwareVersionService) OpenPublicDeveloperLogDownload(firmwareID uint) (*firmwarePackageDownload, error) {
	var firmware deviceModel.FirmwareVersion
	if err := global.GVA_DB.Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return nil, err
	}
	if firmware.PublishStatus != "published" {
		return nil, errors.New("只有已发布固件的更新日志支持公开下载")
	}
	return openDeveloperLogDownload(global.GVA_DB, firmware)
}

func openFirmwarePackageDownload(tx *gorm.DB, firmware deviceModel.FirmwareVersion) (*firmwarePackageDownload, error) {
	fileRecord, err := resolveFirmwarePackageFile(tx, firmware)
	if err != nil {
		return nil, err
	}
	fileName := strings.TrimSpace(firmware.PackageName)
	if fileName == "" {
		fileName = strings.TrimSpace(fileRecord.Name)
	}
	return openFileRecordDownload(fileRecord, firmware.PackageURL, fileName)
}

func openFirmwareLogPackageDownload(tx *gorm.DB, log deviceModel.FirmwareVersionLog) (*firmwarePackageDownload, error) {
	fileRecord, err := resolveFirmwareLogPackageFile(tx, log)
	if err != nil {
		return nil, err
	}
	fileName := strings.TrimSpace(log.PackageName)
	if fileName == "" {
		fileName = strings.TrimSpace(fileRecord.Name)
	}
	return openFileRecordDownload(fileRecord, log.PackageURL, fileName)
}

func resolveDeveloperLogFile(tx *gorm.DB, firmware deviceModel.FirmwareVersion) (exampleModel.ExaFileUploadAndDownload, error) {
	if firmware.LogFileID > 0 {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("id = ?", firmware.LogFileID).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	return exampleModel.ExaFileUploadAndDownload{}, errors.New("未找到更新日志文件记录")
}

func openDeveloperLogDownload(tx *gorm.DB, firmware deviceModel.FirmwareVersion) (*firmwarePackageDownload, error) {
	fileRecord, err := resolveDeveloperLogFile(tx, firmware)
	if err != nil {
		return nil, err
	}
	fileName := strings.TrimSpace(firmware.LogFileName)
	if fileName == "" {
		fileName = strings.TrimSpace(fileRecord.Name)
	}
	return openFileRecordDownload(fileRecord, "", fileName)
}

func openFileRecordDownload(fileRecord exampleModel.ExaFileUploadAndDownload, fallbackURL, fallbackName string) (*firmwarePackageDownload, error) {
	switch global.GVA_CONFIG.System.OssType {
	case "local":
		return openLocalPackageDownload(fileRecord, fallbackName)
	case "minio":
		if download, err := openMinioPackageDownload(fileRecord, fallbackName); err == nil {
			return download, nil
		}
	}
	download, err := openHTTPPackageDownload(strings.TrimSpace(fileRecord.Url), strings.TrimSpace(fileRecord.Name))
	if err == nil {
		if strings.TrimSpace(fallbackName) != "" {
			download.FileName = strings.TrimSpace(fallbackName)
		}
		return download, nil
	}
	if strings.TrimSpace(fallbackURL) != "" {
		return openHTTPPackageDownload(strings.TrimSpace(fallbackURL), fallbackName)
	}
	return nil, err
}

func openLocalPackageDownload(fileRecord exampleModel.ExaFileUploadAndDownload, fallbackName string) (*firmwarePackageDownload, error) {
	key := strings.TrimSpace(fileRecord.Key)
	if key == "" {
		key = filepath.Base(strings.TrimSpace(fileRecord.Url))
	}
	if key == "" {
		return nil, errors.New("未找到本地文件路径")
	}
	fullPath := filepath.Join(global.GVA_CONFIG.Local.StorePath, key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	fileName := strings.TrimSpace(fallbackName)
	if fileName == "" {
		fileName = strings.TrimSpace(fileRecord.Name)
	}
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName)))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &firmwarePackageDownload{
		Reader:      file,
		Size:        info.Size(),
		ContentType: contentType,
		FileName:    fileName,
	}, nil
}

func openMinioPackageDownload(fileRecord exampleModel.ExaFileUploadAndDownload, fallbackName string) (*firmwarePackageDownload, error) {
	key := strings.TrimSpace(fileRecord.Key)
	if key == "" {
		return nil, errors.New("未找到 MinIO 文件键")
	}
	client, err := upload.GetMinio(
		global.GVA_CONFIG.Minio.Endpoint,
		global.GVA_CONFIG.Minio.AccessKeyId,
		global.GVA_CONFIG.Minio.AccessKeySecret,
		global.GVA_CONFIG.Minio.BucketName,
		global.GVA_CONFIG.Minio.UseSSL,
	)
	if err != nil {
		return nil, err
	}
	object, err := client.Client.GetObject(
		context.Background(),
		global.GVA_CONFIG.Minio.BucketName,
		key,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, err
	}
	fileName := strings.TrimSpace(fallbackName)
	if fileName == "" {
		fileName = strings.TrimSpace(fileRecord.Name)
	}
	contentType := strings.TrimSpace(stat.ContentType)
	if contentType == "" {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName)))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &firmwarePackageDownload{
		Reader:      object,
		Size:        stat.Size,
		ContentType: contentType,
		FileName:    fileName,
	}, nil
}

func openHTTPPackageDownload(rawURL, fallbackName string) (*firmwarePackageDownload, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, errors.New("未找到文件下载地址")
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		resp.Body.Close()
		return nil, fmt.Errorf("下载文件失败, HTTP状态码: %d", resp.StatusCode)
	}
	fileName := strings.TrimSpace(fallbackName)
	if fileName == "" {
		if parsedURL, parseErr := url.Parse(rawURL); parseErr == nil {
			fileName = filepath.Base(parsedURL.Path)
		}
	}
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName)))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &firmwarePackageDownload{
		Reader:      resp.Body,
		Size:        resp.ContentLength,
		ContentType: contentType,
		FileName:    fileName,
	}, nil
}

// GetPublicFirmwareDownloadPage 获取公开固件下载页数据
func (s *FirmwareVersionService) GetPublicFirmwareDownloadPage(categoryID, modelID uint) (deviceResp.PublicFirmwareDownloadPageResponse, error) {
	resp := deviceResp.PublicFirmwareDownloadPageResponse{}

	categoryStatus := 1
	categories, _, err := (&DeviceCategoryService{}).GetDeviceCategoryInfoList(deviceReq.DeviceCategorySearch{
		Status: &categoryStatus,
		PageInfo: commonReq.PageInfo{
			Page:     1,
			PageSize: 999,
		},
	})
	if err != nil {
		return resp, err
	}
	resp.Categories = categories
	if len(categories) == 0 {
		return resp, nil
	}

	selectedCategoryID := categoryID
	if selectedCategoryID == 0 || !containsCategoryID(categories, selectedCategoryID) {
		selectedCategoryID = categories[0].ID
	}

	modelStatus := 1
	models, _, err := (&DeviceModelService{}).GetDeviceModelInfoList(deviceReq.DeviceModelSearch{
		CategoryID: selectedCategoryID,
		Status:     &modelStatus,
		PageInfo: commonReq.PageInfo{
			Page:     1,
			PageSize: 999,
		},
	})
	if err != nil {
		return resp, err
	}
	resp.Models = models
	resp.SelectedCategoryID = selectedCategoryID
	if len(models) == 0 {
		return resp, nil
	}

	selectedModelID := modelID
	if selectedModelID == 0 || !containsModelID(models, selectedModelID) {
		selectedModelID = models[0].ID
	}
	resp.SelectedModelID = selectedModelID

	selectedModel := models[0]
	for _, model := range models {
		if model.ID == selectedModelID {
			selectedModel = model
			break
		}
	}

	var rels []deviceModel.ModelFirmwareRel
	if err := global.GVA_DB.Preload("Firmware").
		Where("model_id = ?", selectedModelID).
		Order("id desc").
		Find(&rels).Error; err != nil {
		return resp, err
	}

	items := make([]deviceResp.PublicFirmwareDownloadItem, 0, len(rels))
	relationMap := map[uint]deviceModel.ModelFirmwareRel{}
	for _, rel := range rels {
		if rel.Firmware.PublishStatus != "published" {
			continue
		}
		relationMap[rel.FirmwareID] = rel
		pkgURL, pkgName, _, pkgSize, checksum := resolveFirmwareLogPackageSnapshot(global.GVA_DB, rel.FirmwareID, &rel.Firmware)
		firmware := rel.Firmware
		if pkgURL != "" {
			firmware.PackageURL = pkgURL
		}
		if pkgName != "" {
			firmware.PackageName = pkgName
		}
		if checksum != "" {
			firmware.Checksum = checksum
		}
		items = append(items, deviceResp.PublicFirmwareDownloadItem{
			RelationID:    rel.ID,
			Category:      selectedModel.Category,
			Model:         selectedModel,
			Firmware:      firmware,
			IsRecommended: rel.IsRecommended,
			PackageSize:   pkgSize,
		})
	}

	var latestFirmwareID uint
	var latestFirmwareTime time.Time
	for i := range items {
		itemTime := publicFirmwareItemTime(items[i])
		if latestFirmwareID == 0 ||
			itemTime.After(latestFirmwareTime) ||
			(itemTime.Equal(latestFirmwareTime) && items[i].Firmware.ID > latestFirmwareID) {
			latestFirmwareID = items[i].Firmware.ID
			latestFirmwareTime = itemTime
		}
	}
	for i := range items {
		items[i].Firmware.IsLatest = items[i].Firmware.ID == latestFirmwareID
	}

	sort.SliceStable(items, func(i, j int) bool {
		rankI := publicFirmwareItemRank(items[i])
		rankJ := publicFirmwareItemRank(items[j])
		if rankI != rankJ {
			return rankI < rankJ
		}
		timeI := publicFirmwareItemTime(items[i])
		timeJ := publicFirmwareItemTime(items[j])
		if !timeI.Equal(timeJ) {
			return timeI.After(timeJ)
		}
		return items[i].Firmware.ID > items[j].Firmware.ID
	})

	var current *deviceResp.PublicFirmwareDownloadItem
	var stable *deviceResp.PublicFirmwareDownloadItem
	var latest *deviceResp.PublicFirmwareDownloadItem

	for i := range items {
		item := items[i]
		switch {
		case current == nil && item.IsRecommended:
			current = &items[i]
		case stable == nil && item.Firmware.IsStable:
			stable = &items[i]
		case latest == nil && item.Firmware.IsLatest:
			latest = &items[i]
		}
	}

	if current == nil && stable != nil {
		current = stable
	}
	if current == nil && latest != nil {
		current = latest
	}
	if current == nil && len(items) > 0 {
		current = &items[0]
	}

	selectedIDs := map[uint]struct{}{}
	if current != nil {
		selectedIDs[current.Firmware.ID] = struct{}{}
	}
	if stable != nil {
		selectedIDs[stable.Firmware.ID] = struct{}{}
	}
	if latest != nil {
		selectedIDs[latest.Firmware.ID] = struct{}{}
	}

	history := make([]deviceResp.PublicFirmwareDownloadItem, 0, len(items))
	for _, item := range items {
		if _, exists := selectedIDs[item.Firmware.ID]; exists {
			continue
		}
		history = append(history, item)
	}

	resp.Current = current
	resp.Stable = stable
	resp.Latest = latest
	resp.History = history
	resp.Packages = buildPublicFirmwareDownloadPackageItems(global.GVA_DB, selectedModel, relationMap, items)
	if current != nil {
		resp.PrimaryType = publicFirmwarePrimaryType(*current)
	}
	return resp, nil
}

func buildPublicFirmwareDownloadPackageItems(tx *gorm.DB, model deviceModel.DeviceModel, relationMap map[uint]deviceModel.ModelFirmwareRel, publishedItems []deviceResp.PublicFirmwareDownloadItem) []deviceResp.PublicFirmwareDownloadPackageItem {
	items := make([]deviceResp.PublicFirmwareDownloadPackageItem, 0, len(publishedItems))
	for _, published := range publishedItems {
		firmware := published.Firmware
		if firmware.ID == 0 || firmware.PublishStatus != "published" {
			continue
		}
		if firmware.PackageURL == "" && firmware.PackageName == "" && firmware.PackageFileID == 0 {
			continue
		}
		fileRecord, _ := resolveFirmwarePackageFile(tx, firmware)
		rel, ok := relationMap[firmware.ID]
		item := deviceResp.PublicFirmwareDownloadPackageItem{
			LogID:         firmware.ID,
			RelationID:    0,
			Category:      model.Category,
			Model:         model,
			Firmware:      firmware,
			Action:        "official",
			OperateAt:     firmware.PublishedAt,
			IsRecommended: ok && rel.IsRecommended,
			PackageSize:   fileRecord.Size,
		}
		if item.OperateAt == nil {
			item.OperateAt = firmware.UploadedAt
		}
		if ok {
			item.RelationID = rel.ID
		}
		items = append(items, item)
	}

	return items
}

func containsCategoryID(categories []deviceModel.DeviceCategory, id uint) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

func containsModelID(models []deviceModel.DeviceModel, id uint) bool {
	for _, model := range models {
		if model.ID == id {
			return true
		}
	}
	return false
}

func publicFirmwareItemRank(item deviceResp.PublicFirmwareDownloadItem) int {
	switch {
	case item.IsRecommended:
		return 0
	case item.Firmware.IsStable:
		return 1
	case item.Firmware.IsLatest:
		return 2
	default:
		return 3
	}
}

func publicFirmwareItemTime(item deviceResp.PublicFirmwareDownloadItem) time.Time {
	if item.Firmware.PublishedAt != nil {
		return *item.Firmware.PublishedAt
	}
	if item.Firmware.UploadedAt != nil {
		return *item.Firmware.UploadedAt
	}
	return time.Time{}
}

func publicFirmwarePrimaryType(item deviceResp.PublicFirmwareDownloadItem) string {
	switch {
	case item.IsRecommended:
		return "recommended"
	case item.Firmware.IsStable:
		return "stable"
	case item.Firmware.IsLatest:
		return "latest"
	default:
		return "published"
	}
}
