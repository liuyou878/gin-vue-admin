package device

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	emailUtils "github.com/flipped-aurora/gin-vue-admin/server/plugin/email/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var validModelTestResults = map[string]bool{
	"pending":     true,
	"testing":     true,
	"passed":      true,
	"failed":      true,
	"tested_pass": true,
	"test_failed": true,
	"test_passed": true,
}

// CreateModelFirmwareRel 创建型号固件关系
func (s *ModelFirmwareRelService) CreateModelFirmwareRel(rel *deviceModel.ModelFirmwareRel) error {
	var mailPayload *firmwareActionMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		var firmware deviceModel.FirmwareVersion
		if err := tx.Where("id = ?", rel.FirmwareID).First(&firmware).Error; err != nil {
			return err
		}
		mailPayload = buildFirmwareActionMailPayload(tx, rel.FirmwareID, firmwareActionMailOptions{
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
		if sendErr := sendFirmwareActionEmail(mailPayload, rel.NotifyTo); sendErr != nil {
			return fmt.Errorf("固件已新增，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

// DeleteModelFirmwareRel 删除型号固件关系
func (s *ModelFirmwareRelService) DeleteModelFirmwareRel(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var rel deviceModel.ModelFirmwareRel
		if err := tx.Where("id = ?", id).First(&rel).Error; err != nil {
			return err
		}
		var firmware deviceModel.FirmwareVersion
		if err := tx.Select("publish_status").Where("id = ?", rel.FirmwareID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus == "published" || firmware.PublishStatus == "voided" {
			return errors.New("已发布版本不能删除关联，只能保留历史或作废")
		}
		return tx.Unscoped().Delete(&deviceModel.ModelFirmwareRel{}, "id = ?", id).Error
	})
}

// DeleteModelFirmwareRelByIds 批量删除型号固件关系
func (s *ModelFirmwareRelService) DeleteModelFirmwareRelByIds(ids commonReq.IdsReq) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var rels []deviceModel.ModelFirmwareRel
		if err := tx.Where("id in ?", ids.Ids).Find(&rels).Error; err != nil {
			return err
		}
		var firmwareIDs []uint
		for _, rel := range rels {
			firmwareIDs = append(firmwareIDs, rel.FirmwareID)
		}
		if len(firmwareIDs) > 0 {
			var count int64
			if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id in ? AND publish_status in ?", firmwareIDs, []string{"published", "voided"}).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("所选关联中存在已发布或已作废版本，不能删除")
			}
		}
		return tx.Unscoped().Delete(&[]deviceModel.ModelFirmwareRel{}, "id in ?", ids.Ids).Error
	})
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
		var firmware deviceModel.FirmwareVersion
		if err := tx.Select("status", "publish_status").Where("id = ?", rel.FirmwareID).First(&firmware).Error; err != nil {
			return err
		}
		if firmware.PublishStatus != "published" {
			return errors.New("只有已发布版本才能设为当前发布")
		}
		firmwareStatus, err := getFirmwareStatusForLog(tx, rel.FirmwareID)
		if err != nil {
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
		return createFirmwareVersionLog(tx, rel.FirmwareID, &modelID, "set_recommended", "", firmwareStatus, req.Operator, content)
	})
}

// SetModelFirmwareTestResult 设置型号固件测试结果
func (s *ModelFirmwareRelService) SetModelFirmwareTestResult(req deviceReq.SetModelFirmwareTestResultRequest) error {
	result := normalizeModelTestResult(req.TestResult)
	if !validModelTestResults[result] {
		global.GVA_LOG.Warn("收到非法的型号固件测试结果", zap.String("testResult", req.TestResult))
		return fmt.Errorf("测试结果不合法: %q", req.TestResult)
	}
	var mailPayload *modelFirmwareTestResultMailPayload
	txErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var rel deviceModel.ModelFirmwareRel
		if err := tx.Preload("Model.Category").Preload("Firmware").Where("id = ?", req.ID).First(&rel).Error; err != nil {
			return err
		}
		testedAt := req.TestedAt
		if testedAt == nil {
			now := time.Now()
			testedAt = &now
		}
		currentTester := req.Tester
		if currentTester == "" {
			currentTester = req.Operator
		}
		if currentTester == "" {
			currentTester = rel.Tester
		}
		if currentTester == "" {
			currentTester = rel.Firmware.UploadedBy
		}
		updates := map[string]interface{}{
			"test_result": result,
			"tester":      currentTester,
			"tested_at":   testedAt,
		}
		if err := tx.Model(&deviceModel.ModelFirmwareRel{}).Where("id = ?", rel.ID).Updates(updates).Error; err != nil {
			return err
		}
		firmwareStatus := result
		if result == "tested_pass" {
			firmwareStatus = "pending_release"
		}
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", rel.FirmwareID).Update("status", firmwareStatus).Error; err != nil {
			return err
		}
		modelID := rel.ModelID
		content := req.Content
		if content == "" {
			content = "更新型号固件测试结果"
		}
		if err := createFirmwareVersionLog(tx, rel.FirmwareID, &modelID, convertTestResultToAction(result), normalizeModelTestResult(rel.TestResult), firmwareStatus, req.Operator, content); err != nil {
			return err
		}
		mailPayload = buildModelFirmwareTestResultMailPayload(rel, result, currentTester, *testedAt, req.Operator, content)
		return nil
	})
	if txErr != nil {
		return txErr
	}
	var sendErr error
	if mailPayload != nil {
		sendErr = sendModelFirmwareTestResultEmail(mailPayload, req.NotifyTo)
		if sendErr != nil {
			global.GVA_LOG.Error("发送固件测试结果邮件失败", zap.Error(sendErr))
			return fmt.Errorf("测试结果已保存，但邮件发送失败: %w", sendErr)
		}
	}
	return nil
}

func convertTestResultToAction(testResult string) string {
	testResult = normalizeModelTestResult(testResult)
	switch testResult {
	case "testing":
		return "start_testing"
	case "tested_pass":
		return "test_pass"
	case "test_failed":
		return "test_fail"
	default:
		return "bind_model"
	}
}

func normalizeModelTestResult(testResult string) string {
	testResult = strings.TrimSpace(strings.ToLower(testResult))
	if testResult == "" {
		return "pending"
	}
	switch testResult {
	case "test_passed", "passed", "tested_pass":
		return "tested_pass"
	case "failed", "test_failed":
		return "test_failed"
	case "testing":
		return "testing"
	case "pending":
		return "pending"
	default:
		return testResult
	}
}

type modelFirmwareTestResultMailPayload struct {
	Subject string
	Body    string
}

func buildModelFirmwareTestResultMailPayload(
	rel deviceModel.ModelFirmwareRel,
	result string,
	tester string,
	testedAt time.Time,
	operator string,
	content string,
) *modelFirmwareTestResultMailPayload {
	firmware := rel.Firmware
	model := rel.Model
	resultLabel := map[string]string{
		"testing":     "测试中",
		"tested_pass": "测试通过",
		"test_failed": "测试不通过",
		"pending":     "待测试",
	}[result]
	if resultLabel == "" {
		resultLabel = result
	}
	versionLabel := firmware.VersionCode
	if versionLabel == "" {
		versionLabel = firmware.VersionName
	}
	if versionLabel == "" {
		versionLabel = "未知版本"
	}
	modelLabel := model.ModelName
	if modelLabel == "" {
		modelLabel = "未知型号"
	}
	subject := fmt.Sprintf("固件测试结果通知 - %s - %s", modelLabel, versionLabel)
	if resultLabel != "" {
		subject = fmt.Sprintf("%s - %s", subject, resultLabel)
	}
	body := buildModelFirmwareTestResultMailBody(
		firmware,
		model,
		resultLabel,
		tester,
		testedAt,
		operator,
		content,
	)
	return &modelFirmwareTestResultMailPayload{
		Subject: subject,
		Body:    body,
	}
}

func buildModelFirmwareTestResultMailBody(
	firmware deviceModel.FirmwareVersion,
	model deviceModel.DeviceModel,
	resultLabel string,
	tester string,
	testedAt time.Time,
	operator string,
	content string,
) string {
	categoryName := ""
	if model.Category.ID > 0 {
		categoryName = model.Category.Name
	}
	testedAtText := testedAt.Format("2006-01-02 15:04:05")
	if testedAt.IsZero() {
		testedAtText = "-"
	}
	escapeText := func(value string) string {
		value = strings.TrimSpace(value)
		if value == "" {
			return "-"
		}
		return strings.ReplaceAll(html.EscapeString(value), "\n", "<br/>")
	}
	var builder strings.Builder
	builder.WriteString(`<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; line-height: 1.7; color: #1f2937;">`)
	builder.WriteString(`<h2 style="margin: 0 0 16px; font-size: 20px;">固件测试结果通知</h2>`)
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
	writeRow("设备类别", escapeText(categoryName))
	writeRow("设备型号", escapeText(model.ModelName))
	writeRow("版本号", escapeText(firmware.VersionCode))
	writeRow("版本名称", escapeText(firmware.VersionName))
	writeRow("测试结果", escapeText(resultLabel))
	writeRow("测试人", escapeText(tester))
	writeRow("测试时间", escapeText(testedAtText))
	writeRow("操作人", escapeText(operator))
	writeRow("测试说明", escapeText(content))
	builder.WriteString(`</table>`)
	builder.WriteString(`<div style="margin-top: 16px; color: #6b7280; font-size: 12px;">`)
	builder.WriteString(`该邮件由固件流程管理自动发送`)
	builder.WriteString(`</div>`)
	builder.WriteString(`</div>`)
	return builder.String()
}

func sendModelFirmwareTestResultEmail(payload *modelFirmwareTestResultMailPayload, notifyTo string) error {
	recipient := normalizeEmailRecipients(notifyTo)
	if recipient == "" {
		return nil
	}
	return emailUtils.Email(recipient, payload.Subject, payload.Body)
}

func normalizeEmailRecipients(raw string) string {
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

func getFirmwareStatusForLog(tx *gorm.DB, firmwareID uint) (string, error) {
	var firmware deviceModel.FirmwareVersion
	if err := tx.Select("status").Where("id = ?", firmwareID).First(&firmware).Error; err != nil {
		return "", err
	}
	if firmware.Status == "draft" {
		return "pending_test", nil
	}
	return firmware.Status, nil
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
