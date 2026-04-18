package device

import (
	"errors"
	"sort"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	deviceResp "github.com/flipped-aurora/gin-vue-admin/server/model/device/response"
	exampleModel "github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
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
}

// CreateFirmwareVersion 创建固件版本
func (s *FirmwareVersionService) CreateFirmwareVersion(firmware *deviceModel.FirmwareVersion) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if firmware.Status == "" || firmware.Status == "draft" {
			firmware.Status = "pending_test"
		}
		if !validFirmwareStatuses[firmware.Status] {
			return errors.New("开发状态不合法")
		}
		if firmware.PublishStatus == "" {
			firmware.PublishStatus = "unpublished"
		}
		if err := tx.Create(firmware).Error; err != nil {
			return err
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, "upload", "", firmware.Status, firmware.UploadedBy, "创建固件版本", firmware)
	})
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
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", firmware.ID).Updates(updates).Error; err != nil {
			return err
		}
		if packageChanged {
			content := "更新固件包"
			if current.Status == "test_failed" {
				content = "已修复，重新上传固件包"
			}
			return createFirmwareVersionLog(
				tx,
				current.ID,
				nil,
				"fix_upload",
				current.Status,
				nextStatus,
				firmware.UploadedBy,
				content,
			)
		}
		return nil
	})
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
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		return createFirmwareVersionLog(tx, firmware.ID, nil, action, fromStatus, req.Status, operator, content)
	})
}

// PublishFirmwareVersion 发布固件版本
func (s *FirmwareVersionService) PublishFirmwareVersion(req deviceReq.PublishFirmwareVersionRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		if !req.Direct {
			switch normalizeFirmwareStatus(firmware.Status) {
			case "tested_pass", "pending_release":
			default:
				return errors.New("测试通过后才能发布，或使用直接发布")
			}
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
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		content := req.Content
		if content == "" {
			if req.Direct {
				content = "直接发布版本"
			} else {
				content = "发布版本"
			}
		}
		return createFirmwareVersionLog(tx, firmware.ID, nil, "publish", firmware.Status, "published", operator, content)
	})
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
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		return createFirmwareVersionLog(tx, firmware.ID, nil, "void_release", firmware.PublishStatus, "voided", operator, content)
	})
}

// OnShelfFirmwareVersion 上架已下架固件版本
func (s *FirmwareVersionService) OnShelfFirmwareVersion(req deviceReq.OnShelfFirmwareVersionRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		return createFirmwareVersionLog(tx, firmware.ID, nil, "on_shelf_release", "voided", "published", operator, content)
	})
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
	if firmware.PackageName != "" {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("name = ?", firmware.PackageName).First(&fileRecord).Error; err == nil {
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
	if log.PackageName != "" {
		var fileRecord exampleModel.ExaFileUploadAndDownload
		if err := tx.Where("name = ?", log.PackageName).First(&fileRecord).Error; err == nil {
			return fileRecord, nil
		}
	}
	return exampleModel.ExaFileUploadAndDownload{}, errors.New("未找到安装包文件记录")
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
	firmwareIDs := make([]uint, 0, len(publishedItems))
	for _, item := range publishedItems {
		if item.Firmware.ID == 0 {
			continue
		}
		firmwareIDs = append(firmwareIDs, item.Firmware.ID)
	}
	if len(firmwareIDs) == 0 {
		return nil
	}

	var logs []deviceModel.FirmwareVersionLog
	if err := tx.Preload("Firmware").
		Where("firmware_id IN ? AND action IN ? AND package_url <> ''", firmwareIDs, []string{"upload", "fix_upload"}).
		Order("operate_at desc, id desc").
		Find(&logs).Error; err != nil {
		return nil
	}

	items := make([]deviceResp.PublicFirmwareDownloadPackageItem, 0, len(logs))
	for _, log := range logs {
		if _, err := resolveFirmwareLogPackageFile(tx, log); err != nil {
			continue
		}
		firmware := log.Firmware
		firmware.PackageURL = log.PackageURL
		firmware.PackageName = log.PackageName
		firmware.PackageFileID = log.PackageFileID
		firmware.Checksum = log.Checksum
		rel, ok := relationMap[log.FirmwareID]
		item := deviceResp.PublicFirmwareDownloadPackageItem{
			LogID:         log.ID,
			RelationID:    0,
			Category:      model.Category,
			Model:         model,
			Firmware:      firmware,
			Action:        log.Action,
			OperateAt:     log.OperateAt,
			IsRecommended: ok && rel.IsRecommended,
			PackageSize:   log.PackageSize,
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
