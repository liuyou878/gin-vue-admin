package service

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var WorkOrder = new(workOrderSvc)

type workOrderSvc struct{}

func (s *workOrderSvc) StartInspection(req *request.StartInspection, inspectorID uint, inspectorName string) error {
	now := time.Now()
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 1", req.ID).First(&batch).Error; err != nil {
			return err
		}
		if err := tx.Model(&batch).Updates(map[string]interface{}{
			"inspector_id":    inspectorID,
			"inspector_name":  inspectorName,
			"inspection_date": &now,
		}).Error; err != nil {
			return err
		}
		return updateBatchStatusWithLog(tx, batch, 2, "检测接收并开始检测", inspectorID, inspectorName, "")
	})
}

func (s *workOrderSvc) StartRecheck(req *request.StartRecheck, inspectorID uint, inspectorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 3", req.ID).First(&batch).Error; err != nil {
			return err
		}

		if req.DeviceID > 0 {
			var device model.ProductionOrderDevice
			if err := tx.Where("id = ? AND batch_id = ? AND status = ?", req.DeviceID, req.ID, "pending_recheck").First(&device).Error; err != nil {
				return errors.New("设备不是待复检状态，不能开始复检")
			}
			return updateDeviceStatusWithLog(tx, device, "rechecking", "开始复检", inspectorID, inspectorName)
		}

		var devices []model.ProductionOrderDevice
		if err := tx.Where("batch_id = ? AND status = ?", req.ID, "pending_recheck").Find(&devices).Error; err != nil {
			return err
		}
		if len(devices) == 0 {
			return errors.New("当前批次没有待复检设备")
		}

		for _, device := range devices {
			if err := updateDeviceStatusWithLog(tx, device, "rechecking", "开始复检", inspectorID, inspectorName); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *workOrderSvc) AssignBatchTemplate(req *request.AssignBatchTemplate, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.ID).First(&batch).Error; err != nil {
			return err
		}
		if batch.Status != 0 {
			return errors.New("仅允许为未进入检测流程的批次选择模板")
		}

		var tmpl model.InspectionTemplate
		if err := tx.Where("id = ?", req.TemplateID).First(&tmpl).Error; err != nil {
			return errors.New("检测模板不存在")
		}

		if err := tx.Model(&batch).Updates(map[string]interface{}{
			"template_id": req.TemplateID,
		}).Error; err != nil {
			return err
		}
		if err := updateBatchStatusWithLog(tx, batch, 1, "生产提交检测接收", operatorID, operatorName, ""); err != nil {
			return err
		}

		return tx.Model(&model.ProductionOrder{}).
			Where("id = ? AND status = 0", batch.ProductionOrderID).
			Update("status", 1).Error
	})
}

func (s *workOrderSvc) AssignOrderTemplate(req *request.AssignOrderTemplate, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order model.ProductionOrder
		if err := tx.Where("id = ?", req.ProductionOrderID).First(&order).Error; err != nil {
			return errors.New("生产订单不存在")
		}

		var tmpl model.InspectionTemplate
		if err := tx.Where("id = ?", req.TemplateID).First(&tmpl).Error; err != nil {
			return errors.New("检测模板不存在")
		}

		var batches []model.ProductionBatch
		if err := tx.Where("production_order_id = ?", req.ProductionOrderID).Find(&batches).Error; err != nil {
			return err
		}
		if len(batches) == 0 {
			return errors.New("当前生产订单没有批次，不能提交检测")
		}

		assignableIDs := make([]uint, 0, len(batches))
		for _, batch := range batches {
			if batch.Status == 0 {
				assignableIDs = append(assignableIDs, batch.ID)
			}
		}
		if len(assignableIDs) == 0 {
			return errors.New("当前生产订单没有未派检批次")
		}

		if err := tx.Model(&model.ProductionBatch{}).
			Where("id IN ?", assignableIDs).
			Update("template_id", req.TemplateID).Error; err != nil {
			return err
		}
		for _, batch := range batches {
			if batch.Status != 0 {
				continue
			}
			if err := updateBatchStatusWithLog(tx, batch, 1, "生产提交检测接收", operatorID, operatorName, ""); err != nil {
				return err
			}
		}

		updates := map[string]interface{}{
			"template_id": req.TemplateID,
		}
		if order.Status == 0 {
			updates["status"] = 1
		}
		if strings.TrimSpace(req.InstrumentCategory) != "" {
			updates["instrument_category"] = strings.TrimSpace(req.InstrumentCategory)
		}
		return tx.Model(&model.ProductionOrder{}).Where("id = ?", req.ProductionOrderID).Updates(updates).Error
	})
}

func (s *workOrderSvc) SaveResults(req *request.SaveInspectionResult, inspectorID uint, inspectorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return err
		}
		for _, dr := range req.DeviceResults {
			var device model.ProductionOrderDevice
			if err := tx.Where("id = ? AND batch_id = ?", dr.DeviceID, req.BatchID).First(&device).Error; err != nil {
				return err
			}
			if batch.Status >= 3 && device.Status != "rechecking" {
				return fmt.Errorf("设备 %s 当前不是复检中，不能修改检测结果", device.SN)
			}
			var existing model.InspectionDeviceResult
			err := tx.Where("production_order_device_id = ? AND item_id = ?", dr.DeviceID, dr.ItemID).First(&existing).Error
			if err == nil {
				updates := map[string]interface{}{
					"pass_result":   dr.PassResult,
					"number_result": dr.NumberResult,
					"remark":        dr.Remark,
				}
				if resultFieldsChanged(existing, dr.PassResult, dr.NumberResult, dr.Remark) ||
					(existing.InspectedAt == nil && resultHasValue(dr.PassResult, dr.NumberResult, dr.Remark)) {
					now := time.Now()
					updates["inspector_id"] = inspectorID
					updates["inspector_name"] = inspectorName
					updates["inspected_at"] = &now
				}
				if err := tx.Model(&existing).Updates(updates).Error; err != nil {
					return err
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				if !resultHasValue(dr.PassResult, dr.NumberResult, dr.Remark) {
					continue
				}
				now := time.Now()
				if err := tx.Create(&model.InspectionDeviceResult{
					ProductionOrderDeviceID: dr.DeviceID,
					ItemID:                  dr.ItemID,
					PassResult:              dr.PassResult,
					NumberResult:            dr.NumberResult,
					Remark:                  dr.Remark,
					InspectorID:             &inspectorID,
					InspectorName:           inspectorName,
					InspectedAt:             &now,
				}).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		for _, ds := range req.DeviceStatuses {
			var device model.ProductionOrderDevice
			if err := tx.Where("id = ? AND batch_id = ?", ds.DeviceID, req.BatchID).First(&device).Error; err != nil {
				return err
			}
			if batch.Status >= 3 && device.Status != "rechecking" {
				return fmt.Errorf("设备 %s 当前不是复检中，不能修改检测结果", device.SN)
			}
			if err := updateDeviceStatusForInspectionSave(tx, device, ds.Status, "保存检测结果"); err != nil {
				return err
			}
		}
		return nil
	})
}

func resultHasValue(passResult *bool, numberResult *float64, remark string) bool {
	return passResult != nil || numberResult != nil || strings.TrimSpace(remark) != ""
}

func resultFieldsChanged(existing model.InspectionDeviceResult, passResult *bool, numberResult *float64, remark string) bool {
	return !boolPtrEqual(existing.PassResult, passResult) ||
		!floatPtrEqual(existing.NumberResult, numberResult) ||
		existing.Remark != remark
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func floatPtrEqual(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (s *workOrderSvc) SaveSingleResult(req *request.SaveSingleInspectionResult, inspectorID uint, inspectorName string) (*model.InspectionDeviceResult, error) {
	var saved model.InspectionDeviceResult
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return err
		}
		var device model.ProductionOrderDevice
		if err := tx.Where("id = ? AND batch_id = ?", req.DeviceID, req.BatchID).First(&device).Error; err != nil {
			return err
		}
		if batch.Status >= 3 && device.Status != "rechecking" {
			return fmt.Errorf("设备 %s 当前不是复检中，不能修改检测结果", device.SN)
		}

		now := time.Now()
		var existing model.InspectionDeviceResult
		err := tx.Where("production_order_device_id = ? AND item_id = ?", req.DeviceID, req.ItemID).First(&existing).Error
		if err == nil {
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"pass_result":    req.PassResult,
				"number_result":  req.NumberResult,
				"remark":         req.Remark,
				"inspector_id":   inspectorID,
				"inspector_name": inspectorName,
				"inspected_at":   &now,
			}).Error; err != nil {
				return err
			}
			existing.PassResult = req.PassResult
			existing.NumberResult = req.NumberResult
			existing.Remark = req.Remark
			existing.InspectorID = &inspectorID
			existing.InspectorName = inspectorName
			existing.InspectedAt = &now
			saved = existing
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if !resultHasValue(req.PassResult, req.NumberResult, req.Remark) {
				return errors.New("没有可保存的检测结果")
			}
			saved = model.InspectionDeviceResult{
				ProductionOrderDeviceID: req.DeviceID,
				ItemID:                  req.ItemID,
				PassResult:              req.PassResult,
				NumberResult:            req.NumberResult,
				Remark:                  req.Remark,
				InspectorID:             &inspectorID,
				InspectorName:           inspectorName,
				InspectedAt:             &now,
			}
			if err := tx.Create(&saved).Error; err != nil {
				return err
			}
		} else {
			return err
		}

		if strings.TrimSpace(req.Status) != "" && !(batch.Status >= 3 && device.Status == "rechecking") {
			if err := updateDeviceStatusForInspectionSave(tx, device, req.Status, "保存单项检测结果"); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &saved, nil
}

func (s *workOrderSvc) ReturnDevices(req *request.ReturnDevices, returnByID uint, returnByName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return err
		}
		if batch.Status != 3 {
			return errors.New("只有待确认批次的不合格设备才可以打回生产")
		}

		now := time.Now()
		for _, deviceID := range req.DeviceIDs {
			var device model.ProductionOrderDevice
			if err := tx.Where("id = ? AND batch_id = ?", deviceID, req.BatchID).First(&device).Error; err != nil {
				return fmt.Errorf("设备不存在或不属于该批次: %d", deviceID)
			}
			if device.Status != "fail" {
				return fmt.Errorf("设备 %s 不是不合格状态，不能打回生产", device.SN)
			}
			if err := tx.Model(&device).Updates(map[string]interface{}{
				"return_reason":  req.Reason,
				"return_at":      &now,
				"return_by_id":   returnByID,
				"return_by_name": returnByName,
			}).Error; err != nil {
				return err
			}
			device.ReturnReason = req.Reason
			device.ReturnAt = &now
			device.ReturnByID = &returnByID
			device.ReturnByName = returnByName
			if err := updateDeviceStatusWithLog(tx, device, "returned", req.Reason, returnByID, returnByName); err != nil {
				return err
			}
		}
		return nil
	})
}

func updateDeviceStatusWithLog(tx *gorm.DB, device model.ProductionOrderDevice, nextStatus string, reason string, operatorID interface{}, operatorName string) error {
	if device.Status == nextStatus {
		return nil
	}
	updates := map[string]interface{}{"status": nextStatus}
	if err := tx.Model(&model.ProductionOrderDevice{}).Where("id = ?", device.ID).Updates(updates).Error; err != nil {
		return err
	}
	return updateDeviceStatusLogOnly(tx, device, nextStatus, reason, operatorID, operatorName)
}

func updateDeviceStatusForInspectionSave(tx *gorm.DB, device model.ProductionOrderDevice, nextStatus string, reason string) error {
	if device.Status == nextStatus {
		return nil
	}
	if strings.TrimSpace(nextStatus) == "pending" {
		return tx.Model(&model.ProductionOrderDevice{}).Where("id = ?", device.ID).Update("status", nextStatus).Error
	}
	return updateDeviceStatusWithLog(tx, device, nextStatus, reason, nil, "")
}

func updateDeviceStatusLogOnly(tx *gorm.DB, device model.ProductionOrderDevice, nextStatus string, reason string, operatorID interface{}, operatorName string) error {
	var opID *uint
	switch v := operatorID.(type) {
	case uint:
		opID = &v
	case *uint:
		opID = v
	}
	return tx.Create(&model.ProductionOrderDeviceStatusLog{
		ProductionOrderDeviceID: device.ID,
		FromStatus:              device.Status,
		ToStatus:                nextStatus,
		Reason:                  reason,
		OperatorID:              opID,
		OperatorName:            operatorName,
	}).Error
}

func updateBatchStatusWithLog(tx *gorm.DB, batch model.ProductionBatch, nextStatus int, action string, operatorID interface{}, operatorName string, reason string) error {
	if batch.Status == nextStatus {
		return nil
	}
	updates := map[string]interface{}{"status": nextStatus}
	if err := tx.Model(&model.ProductionBatch{}).Where("id = ?", batch.ID).Updates(updates).Error; err != nil {
		return err
	}
	var opID *uint
	switch v := operatorID.(type) {
	case uint:
		opID = &v
	case *uint:
		opID = v
	}
	return tx.Create(&model.ProductionBatchStatusLog{
		ProductionBatchID: batch.ID,
		FromStatus:        batch.Status,
		ToStatus:          nextStatus,
		Action:            action,
		Reason:            reason,
		OperatorID:        opID,
		OperatorName:      operatorName,
	}).Error
}

func (s *workOrderSvc) CompleteInspection(req *request.CompleteInspection, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 2", req.ID).First(&batch).Error; err != nil {
			return err
		}

		var devices []model.ProductionOrderDevice
		if err := tx.Where("batch_id = ?", req.ID).Find(&devices).Error; err != nil {
			return err
		}
		if len(devices) == 0 {
			return errors.New("批次下没有设备，不能完成检测")
		}

		var templateItems []model.InspectionTemplateItem
		if batch.TemplateID != nil {
			if err := tx.Preload("Item").Where("template_id = ?", *batch.TemplateID).Find(&templateItems).Error; err != nil {
				return err
			}
		}
		if len(templateItems) == 0 {
			return errors.New("批次没有检测模板或检测项，不能完成检测")
		}

		deviceIDs := make([]uint, 0, len(devices))
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.ID)
		}

		var results []model.InspectionDeviceResult
		if err := tx.Where("production_order_device_id IN ?", deviceIDs).Find(&results).Error; err != nil {
			return err
		}
		type resultKey struct {
			DeviceID uint
			ItemID   uint
		}
		resultMap := make(map[resultKey]model.InspectionDeviceResult, len(results))
		for _, result := range results {
			resultMap[resultKey{DeviceID: result.ProductionOrderDeviceID, ItemID: result.ItemID}] = result
		}

		for _, device := range devices {
			if device.Status == "rework" {
				continue
			}
			for _, templateItem := range templateItems {
				result, ok := resultMap[resultKey{DeviceID: device.ID, ItemID: templateItem.ItemID}]
				if !ok || !inspectionResultCompleted(templateItem.Item.ResultType, result) {
					return fmt.Errorf("设备 %s 的检测项 %s 未完成", device.SN, templateItem.Item.Name)
				}
			}
		}

		return updateBatchStatusWithLog(tx, batch, 3, "检测提交待确认", operatorID, operatorName, "")
	})
}

func (s *workOrderSvc) ConfirmInspectionComplete(req *request.ConfirmInspectionComplete, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 3", req.ID).First(&batch).Error; err != nil {
			return err
		}

		var pendingCount int64
		if err := tx.Model(&model.ProductionOrderDevice{}).
			Where("batch_id = ? AND status IN ?", req.ID, []string{"pending", "fail", "returned", "rework", "pending_recheck", "rechecking"}).
			Count(&pendingCount).Error; err != nil {
			return err
		}
		if pendingCount > 0 {
			return errors.New("还有未闭环设备，不能确认完成")
		}

		return updateBatchStatusWithLog(tx, batch, 4, "完成检测", operatorID, operatorName, "")
	})
}

func (s *workOrderSvc) CompleteRecheck(req *request.CompleteRecheck, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 3", req.ID).First(&batch).Error; err != nil {
			return err
		}

		var devices []model.ProductionOrderDevice
		if err := tx.Model(&model.ProductionOrderDevice{}).Where(
			`batch_id = ? AND (
				status = ?
				OR EXISTS (
					SELECT 1 FROM production_order_device_status_logs logs
					WHERE logs.production_order_device_id = production_order_devices.id
					AND logs.to_status = ?
				)
			)`,
			req.ID,
			"rechecking",
			"rechecking",
		).Find(&devices).Error; err != nil {
			return err
		}
		if len(devices) == 0 {
			return errors.New("当前批次没有复检中的设备")
		}

		var templateItems []model.InspectionTemplateItem
		if batch.TemplateID != nil {
			if err := tx.Preload("Item").Where("template_id = ?", *batch.TemplateID).Find(&templateItems).Error; err != nil {
				return err
			}
		}
		if len(templateItems) == 0 {
			return errors.New("批次没有检测模板或检测项，不能完成复检")
		}

		deviceIDs := make([]uint, 0, len(devices))
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.ID)
		}

		var results []model.InspectionDeviceResult
		if err := tx.Where("production_order_device_id IN ?", deviceIDs).Find(&results).Error; err != nil {
			return err
		}
		type resultKey struct {
			DeviceID uint
			ItemID   uint
		}
		resultMap := make(map[resultKey]model.InspectionDeviceResult, len(results))
		for _, result := range results {
			resultMap[resultKey{DeviceID: result.ProductionOrderDeviceID, ItemID: result.ItemID}] = result
		}

		for _, device := range devices {
			hasFail := false
			for _, templateItem := range templateItems {
				result, ok := resultMap[resultKey{DeviceID: device.ID, ItemID: templateItem.ItemID}]
				if !ok || !inspectionResultCompleted(templateItem.Item.ResultType, result) {
					return fmt.Errorf("设备 %s 的复检项 %s 未完成", device.SN, templateItem.Item.Name)
				}
				if resultIndicatesFail(templateItem.Item, result) {
					hasFail = true
				}
			}
			nextStatus := "pass"
			if hasFail {
				nextStatus = "fail"
			}
			if err := updateDeviceStatusWithLog(tx, device, nextStatus, "完成复检", operatorID, operatorName); err != nil {
				return err
			}
		}

		return nil
	})
}

func updateBatchStatusLogOnly(tx *gorm.DB, batch model.ProductionBatch, action string, operatorID interface{}, operatorName string, reason string) error {
	var opID *uint
	switch v := operatorID.(type) {
	case uint:
		opID = &v
	case *uint:
		opID = v
	}
	return tx.Create(&model.ProductionBatchStatusLog{
		ProductionBatchID: batch.ID,
		FromStatus:        batch.Status,
		ToStatus:          batch.Status,
		Action:            action,
		Reason:            reason,
		OperatorID:        opID,
		OperatorName:      operatorName,
	}).Error
}

func inspectionResultCompleted(resultType string, result model.InspectionDeviceResult) bool {
	switch resultType {
	case "number":
		return result.NumberResult != nil
	case "pass_fail":
		return result.PassResult != nil
	default:
		return result.PassResult != nil && result.NumberResult != nil
	}
}

func resultIndicatesFail(item model.InspectionItem, result model.InspectionDeviceResult) bool {
	switch item.ResultType {
	case "number":
		if result.NumberResult == nil {
			return false
		}
		return (item.MinValue != nil && *result.NumberResult < *item.MinValue) ||
			(item.MaxValue != nil && *result.NumberResult > *item.MaxValue)
	case "pass_fail":
		return result.PassResult != nil && !*result.PassResult
	default:
		passFailFailed := result.PassResult != nil && !*result.PassResult
		numberFailed := result.NumberResult != nil &&
			((item.MinValue != nil && *result.NumberResult < *item.MinValue) ||
				(item.MaxValue != nil && *result.NumberResult > *item.MaxValue))
		return passFailFailed || numberFailed
	}
}

type InspectionResultItem struct {
	ItemID        uint       `json:"itemID"`
	ItemName      string     `json:"itemName"`
	ResultType    string     `json:"resultType"`
	Unit          string     `json:"unit"`
	MinValue      *float64   `json:"minValue"`
	MaxValue      *float64   `json:"maxValue"`
	PassResult    *bool      `json:"passResult"`
	NumberResult  *float64   `json:"numberResult"`
	Remark        string     `json:"remark"`
	InspectorID   *uint      `json:"inspectorID"`
	InspectorName string     `json:"inspectorName"`
	InspectedAt   *time.Time `json:"inspectedAt"`
}

func (s *workOrderSvc) GetInspectionBatchList(search request.InspectionBatchSearch) ([]model.InspectionBatchListItem, int64, error) {
	if search.Page <= 0 {
		search.Page = 1
	}
	if search.PageSize <= 0 {
		search.PageSize = 30
	}

	db := global.GVA_DB.Table("production_batches AS pb").
		Select("pb.id, pb.production_order_id, po.mo_number, pb.batch_number, COALESCE(it.product_name, po.product_name) AS product_name, COALESCE(it.model, po.model) AS model, po.firmware_version, po.mainboard_firmware_version, po.pn_code, po.instrument_category, pb.status, pb.template_id, pb.inspector_id, pb.inspector_name, pb.inspection_date, pb.created_at").
		Joins("JOIN production_orders po ON po.id = pb.production_order_id").
		Joins("LEFT JOIN inspection_templates it ON it.id = pb.template_id")

	if search.MONumber != "" {
		db = db.Where("po.mo_number LIKE ?", "%"+search.MONumber+"%")
	}
	if search.Model != "" {
		db = db.Where("po.model LIKE ?", "%"+search.Model+"%")
	}
	if search.Status != nil {
		db = db.Where("pb.status = ?", *search.Status)
	}
	if search.DeviceStatus != "" {
		statuses := splitDeviceStatuses(search.DeviceStatus)
		db = db.Where(
			"EXISTS (SELECT 1 FROM production_order_devices pod WHERE pod.batch_id = pb.id AND pod.status IN ?)",
			statuses,
		)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.InspectionBatchListItem
	if err := db.Order("pb.id desc").Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	for i := range list {
		var count, pass, fail int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ?", list[i].ID).Count(&count)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ? AND status = ?", list[i].ID, "pass").Count(&pass)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ? AND status = ?", list[i].ID, "fail").Count(&fail)
		var rework int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ? AND status IN ?", list[i].ID, []string{"returned", "rework"}).Count(&rework)
		var recheck int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ? AND status IN ?", list[i].ID, []string{"pending_recheck", "rechecking"}).Count(&recheck)
		var rechecking int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ? AND status = ?", list[i].ID, "rechecking").Count(&rechecking)
		list[i].DeviceCount = int(count)
		list[i].PassCount = int(pass)
		list[i].FailCount = int(fail)
		list[i].ReworkCount = int(rework)
		list[i].RecheckCount = int(recheck)
		list[i].RecheckingCount = int(rechecking)
		if list[i].TemplateID != nil {
			var tmpl model.InspectionTemplate
			if global.GVA_DB.Where("id = ?", *list[i].TemplateID).First(&tmpl).Error == nil {
				list[i].Template = &tmpl
			}
		}
	}

	return list, total, nil
}

func splitDeviceStatuses(value string) []string {
	parts := strings.Split(value, ",")
	statuses := make([]string, 0, len(parts))
	for _, part := range parts {
		status := strings.TrimSpace(part)
		if status != "" {
			statuses = append(statuses, status)
		}
	}
	if len(statuses) == 0 {
		return []string{value}
	}
	return statuses
}

// GetInspectionDetailData returns full detail with devices+results for frontend
func (s *workOrderSvc) GetInspectionDetailData(batchID string) (map[string]interface{}, error) {
	var batch model.ProductionBatch
	err := global.GVA_DB.Preload("Template").Where("id = ?", batchID).First(&batch).Error
	if err != nil {
		return nil, err
	}

	var order model.ProductionOrder
	if err := global.GVA_DB.Where("id = ?", batch.ProductionOrderID).First(&order).Error; err != nil {
		return nil, err
	}

	var devices []model.ProductionOrderDevice
	global.GVA_DB.Where("batch_id = ?", batchID).Order("line_number asc").Find(&devices)

	var templateItems []model.InspectionTemplateItem
	if batch.TemplateID != nil {
		global.GVA_DB.Preload("Item").Where("template_id = ?", *batch.TemplateID).Order("sort asc").Find(&templateItems)
	}

	type ResultKey struct {
		DeviceID uint
		ItemID   uint
	}
	resultMap := make(map[ResultKey]model.InspectionDeviceResult)
	var allResults []model.InspectionDeviceResult
	deviceIDs := make([]uint, len(devices))
	for i, d := range devices {
		deviceIDs[i] = d.ID
	}
	if len(deviceIDs) > 0 {
		global.GVA_DB.Where("production_order_device_id IN ?", deviceIDs).Find(&allResults)
		for _, r := range allResults {
			resultMap[ResultKey{r.ProductionOrderDeviceID, r.ItemID}] = r
		}
	}

	type DeviceInfo struct {
		ID           uint                   `json:"ID"`
		SN           string                 `json:"sn"`
		LineNumber   int                    `json:"lineNumber"`
		Status       string                 `json:"status"`
		ReturnReason string                 `json:"returnReason"`
		ReturnAt     *time.Time             `json:"returnAt"`
		ReturnByName string                 `json:"returnByName"`
		Results      []InspectionResultItem `json:"results"`
	}
	type TemplateItemInfo struct {
		ItemID     uint     `json:"itemID"`
		ItemName   string   `json:"itemName"`
		ResultType string   `json:"resultType"`
		Unit       string   `json:"unit"`
		MinValue   *float64 `json:"minValue"`
		MaxValue   *float64 `json:"maxValue"`
		Sort       int      `json:"sort"`
	}

	deviceInfos := make([]DeviceInfo, len(devices))
	for i, d := range devices {
		deviceInfos[i] = DeviceInfo{
			ID:           d.ID,
			SN:           d.SN,
			LineNumber:   d.LineNumber,
			Status:       d.Status,
			ReturnReason: d.ReturnReason,
			ReturnAt:     d.ReturnAt,
			ReturnByName: d.ReturnByName,
			Results:      make([]InspectionResultItem, len(templateItems)),
		}
		for j, ti := range templateItems {
			key := ResultKey{d.ID, ti.ItemID}
			r := resultMap[key]
			deviceInfos[i].Results[j] = InspectionResultItem{
				ItemID: ti.ItemID, ItemName: ti.Item.Name, ResultType: ti.Item.ResultType,
				Unit: ti.Item.Unit, MinValue: ti.Item.MinValue, MaxValue: ti.Item.MaxValue,
				PassResult: r.PassResult, NumberResult: r.NumberResult, Remark: r.Remark,
				InspectorID: r.InspectorID, InspectorName: r.InspectorName, InspectedAt: r.InspectedAt,
			}
		}
	}

	templateItemInfos := make([]TemplateItemInfo, len(templateItems))
	for i, ti := range templateItems {
		templateItemInfos[i] = TemplateItemInfo{
			ItemID: ti.ItemID, ItemName: ti.Item.Name, ResultType: ti.Item.ResultType,
			Unit: ti.Item.Unit, MinValue: ti.Item.MinValue, MaxValue: ti.Item.MaxValue, Sort: ti.Sort,
		}
	}

	return map[string]interface{}{
		"order": map[string]interface{}{
			"ID":                batch.ID,
			"productionOrderID": batch.ProductionOrderID,
			"moNumber":          order.MONumber,
			"batchNumber":       batch.BatchNumber,
			"productName": func() string {
				if batch.Template != nil && batch.Template.ProductName != "" {
					return batch.Template.ProductName
				}
				return order.ProductName
			}(),
			"model":                    order.Model,
			"firmwareVersion":          order.FirmwareVersion,
			"mainboardFirmwareVersion": order.MainboardFirmwareVersion,
			"pnCode":                   order.PNCode,
			"instrumentCategory":       order.InstrumentCategory,
			"status":                   batch.Status,
			"templateID":               batch.TemplateID,
			"templateName": func() string {
				if batch.Template != nil {
					return batch.Template.Name
				}
				return ""
			}(),
			"inspectorID":    batch.InspectorID,
			"inspectorName":  batch.InspectorName,
			"inspectionDate": batch.InspectionDate,
		},
		"devices":       deviceInfos,
		"templateItems": templateItemInfos,
	}, nil
}

func (s *workOrderSvc) GetBatchStatusLogs(batchID string) ([]model.ProductionBatchStatusLog, error) {
	var logs []model.ProductionBatchStatusLog
	err := global.GVA_DB.Where("production_batch_id = ?", batchID).Order("id asc").Find(&logs).Error
	return logs, err
}

type FlowLogItem struct {
	ID           uint      `json:"ID"`
	Scope        string    `json:"scope"`
	ScopeLabel   string    `json:"scopeLabel"`
	BatchID      uint      `json:"batchID"`
	DeviceID     uint      `json:"deviceID"`
	DeviceSN     string    `json:"deviceSN"`
	Title        string    `json:"title"`
	FromStatus   string    `json:"fromStatus"`
	ToStatus     string    `json:"toStatus"`
	Action       string    `json:"action"`
	Reason       string    `json:"reason"`
	OperatorID   *uint     `json:"operatorID"`
	OperatorName string    `json:"operatorName"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

func (s *workOrderSvc) GetFlowLogs(batchID string, deviceID string) ([]FlowLogItem, error) {
	items := make([]FlowLogItem, 0)
	if strings.TrimSpace(deviceID) != "" {
		var device model.ProductionOrderDevice
		if err := global.GVA_DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
			return nil, err
		}
		var deviceLogs []model.ProductionOrderDeviceStatusLog
		if err := global.GVA_DB.Where("production_order_device_id = ?", deviceID).Order("id asc").Find(&deviceLogs).Error; err != nil {
			return nil, err
		}
		for _, log := range deviceLogs {
			action := deviceFlowAction(log.FromStatus, log.ToStatus, log.Reason)
			items = append(items, FlowLogItem{
				ID:           log.ID,
				Scope:        "device",
				ScopeLabel:   "设备",
				BatchID:      valueOrZero(device.BatchID),
				DeviceID:     device.ID,
				DeviceSN:     device.SN,
				Title:        action,
				FromStatus:   log.FromStatus,
				ToStatus:     log.ToStatus,
				Action:       action,
				Reason:       log.Reason,
				OperatorID:   log.OperatorID,
				OperatorName: log.OperatorName,
				CreatedAt:    log.CreatedAt,
			})
		}
		if strings.TrimSpace(batchID) == "" && device.BatchID != nil {
			batchID = fmt.Sprintf("%d", *device.BatchID)
		}
	}

	if strings.TrimSpace(batchID) != "" {
		var batch model.ProductionBatch
		batchFound := global.GVA_DB.Where("id = ?", batchID).First(&batch).Error == nil
		var batchLogs []model.ProductionBatchStatusLog
		batchLogErr := global.GVA_DB.Where("production_batch_id = ?", batchID).Order("id asc").Find(&batchLogs).Error
		if batchLogErr == nil && len(batchLogs) > 0 {
			for _, log := range batchLogs {
				if isDeviceRecheckBatchLog(log.Action) {
					continue
				}
				items = append(items, FlowLogItem{
					ID:           log.ID,
					Scope:        "batch",
					ScopeLabel:   "批次",
					BatchID:      log.ProductionBatchID,
					Title:        firstNonEmpty(log.Action, "批次流转"),
					FromStatus:   fmt.Sprintf("%d", log.FromStatus),
					ToStatus:     fmt.Sprintf("%d", log.ToStatus),
					Action:       log.Action,
					Reason:       log.Reason,
					OperatorID:   log.OperatorID,
					OperatorName: log.OperatorName,
					CreatedAt:    log.CreatedAt,
				})
			}
		} else if batchFound {
			items = append(items, syntheticBatchFlowLogs(batch)...)
		} else if batchLogErr != nil {
			return nil, batchLogErr
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func isDeviceRecheckBatchLog(action string) bool {
	switch strings.TrimSpace(action) {
	case "开始复检", "完成复检":
		return true
	default:
		return false
	}
}

func deviceFlowAction(fromStatus string, toStatus string, reason string) string {
	fromStatus = strings.TrimSpace(fromStatus)
	toStatus = strings.TrimSpace(toStatus)
	reason = strings.TrimSpace(reason)

	switch {
	case toStatus == "returned":
		return "退回生产"
	case fromStatus == "returned" && toStatus == "rework":
		return "生产接收返工"
	case fromStatus == "rework" && toStatus == "pending_recheck":
		return "返工完成"
	case toStatus == "rechecking":
		return "开始复检"
	case fromStatus == "rechecking" && (toStatus == "pass" || toStatus == "fail"):
		return "完成复检"
	case toStatus == "pass":
		return "判定合格"
	case toStatus == "fail":
		return "判定不合格"
	case reason == "保存检测结果" || reason == "保存单项检测结果":
		return reason
	default:
		return "设备状态变更"
	}
}

func syntheticBatchFlowLogs(batch model.ProductionBatch) []FlowLogItem {
	logTime := batch.UpdatedAt
	if logTime.IsZero() {
		logTime = batch.CreatedAt
	}
	if batch.Status <= 0 {
		return []FlowLogItem{
			{
				ID:         batch.ID*10 + 1,
				Scope:      "batch",
				ScopeLabel: "批次",
				BatchID:    batch.ID,
				Title:      "生产分批",
				FromStatus: "0",
				ToStatus:   "0",
				Action:     "生产分批",
				CreatedAt:  logTime,
			},
		}
	}
	logs := []FlowLogItem{
		{
			ID:           batch.ID*10 + 1,
			Scope:        "batch",
			ScopeLabel:   "批次",
			BatchID:      batch.ID,
			Title:        "生产提交检测接收",
			FromStatus:   "0",
			ToStatus:     "1",
			Action:       "生产提交检测接收",
			Reason:       "历史数据无原始派检日志，系统自动补显",
			OperatorName: "",
			CreatedAt:    logTime,
		},
	}
	if batch.Status >= 2 {
		inspectionTime := logTime
		if batch.InspectionDate != nil {
			inspectionTime = *batch.InspectionDate
		}
		logs = append(logs, FlowLogItem{
			ID:           batch.ID*10 + 2,
			Scope:        "batch",
			ScopeLabel:   "批次",
			BatchID:      batch.ID,
			Title:        "检测接收并开始检测",
			FromStatus:   "1",
			ToStatus:     "2",
			Action:       "检测接收并开始检测",
			Reason:       "历史数据无原始接收日志，系统自动补显",
			OperatorID:   batch.InspectorID,
			OperatorName: batch.InspectorName,
			CreatedAt:    inspectionTime,
		})
	}
	if batch.Status >= 3 {
		logs = append(logs, FlowLogItem{
			ID:         batch.ID*10 + 3,
			Scope:      "batch",
			ScopeLabel: "批次",
			BatchID:    batch.ID,
			Title:      "检测提交待确认",
			FromStatus: "2",
			ToStatus:   "3",
			Action:     "检测提交待确认",
			Reason:     "历史数据无原始待确认日志，系统自动补显",
			CreatedAt:  logTime,
		})
	}
	if batch.Status >= 4 {
		logs = append(logs, FlowLogItem{
			ID:         batch.ID*10 + 4,
			Scope:      "batch",
			ScopeLabel: "批次",
			BatchID:    batch.ID,
			Title:      "完成检测",
			FromStatus: "3",
			ToStatus:   "4",
			Action:     "完成检测",
			Reason:     "历史数据无原始完成日志，系统自动补显",
			CreatedAt:  logTime,
		})
	}
	return logs
}

func valueOrZero(value *uint) uint {
	if value == nil {
		return 0
	}
	return *value
}

func (s *workOrderSvc) ExportInspectionExcel(batchID string) (*bytes.Buffer, string, error) {
	var batch model.ProductionBatch
	if err := global.GVA_DB.Preload("Template").Where("id = ?", batchID).First(&batch).Error; err != nil {
		return nil, "", err
	}

	var order model.ProductionOrder
	if err := global.GVA_DB.Where("id = ?", batch.ProductionOrderID).First(&order).Error; err != nil {
		return nil, "", err
	}

	var devices []model.ProductionOrderDevice
	if err := global.GVA_DB.Where("batch_id = ?", batchID).Order("line_number asc, id asc").Find(&devices).Error; err != nil {
		return nil, "", err
	}

	var templateItems []model.InspectionTemplateItem
	if batch.TemplateID != nil {
		if err := global.GVA_DB.Preload("Item").Where("template_id = ?", *batch.TemplateID).Order("sort asc").Find(&templateItems).Error; err != nil {
			return nil, "", err
		}
	}
	if len(templateItems) == 0 {
		return nil, "", errors.New("批次没有检测模板或检测项，不能导出")
	}

	deviceIDs := make([]uint, 0, len(devices))
	for _, device := range devices {
		deviceIDs = append(deviceIDs, device.ID)
	}
	type exportResultKey struct {
		DeviceID uint
		ItemID   uint
	}
	resultMap := make(map[exportResultKey]model.InspectionDeviceResult)
	if len(deviceIDs) > 0 {
		var results []model.InspectionDeviceResult
		if err := global.GVA_DB.Where("production_order_device_id IN ?", deviceIDs).Find(&results).Error; err != nil {
			return nil, "", err
		}
		for _, result := range results {
			resultMap[exportResultKey{DeviceID: result.ProductionOrderDeviceID, ItemID: result.ItemID}] = result
		}
	}

	f := excelize.NewFile()
	sheet := "检测工单"
	defaultSheet := f.GetSheetName(0)
	f.SetSheetName(defaultSheet, sheet)

	lastColIndex := 5 + len(templateItems)
	if lastColIndex < 8 {
		lastColIndex = 8
	}
	lastCol, _ := excelize.ColumnNumberToName(lastColIndex)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    thinBorders(),
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F3F4F6"}, Pattern: 1},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    thinBorders(),
	})
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border:    thinBorders(),
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F9FAFB"}, Pattern: 1},
	})
	footerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
	})
	_ = f.MergeCell(sheet, "A1", lastCol+"1")
	_ = f.SetCellValue(sheet, "A1", "GNSS接收机产品检测工单")
	_ = f.SetCellStyle(sheet, "A1", lastCol+"1", titleStyle)
	_ = f.SetRowHeight(sheet, 1, 28)

	productName := order.ProductName
	modelName := order.Model
	if batch.Template != nil {
		if batch.Template.ProductName != "" {
			productName = batch.Template.ProductName
		}
		if batch.Template.Model != "" {
			modelName = batch.Template.Model
		}
	}

	infoRows := [][]string{
		{fmt.Sprintf("生产订单号：%s", order.MONumber), fmt.Sprintf("批次号：%s", batch.BatchNumber), fmt.Sprintf("业务类型：%s", businessCategoryLabel(order.InstrumentCategory))},
		{fmt.Sprintf("产品名称：%s", productName), fmt.Sprintf("模板型号：%s", modelName), fmt.Sprintf("PN码：%s", order.PNCode)},
		{fmt.Sprintf("固件版本：%s", order.FirmwareVersion), fmt.Sprintf("主板固件版本：%s", order.MainboardFirmwareVersion), fmt.Sprintf("检测员：%s", batch.InspectorName)},
	}
	infoSegments := splitExcelSegments(lastColIndex, 3)
	for i, row := range infoRows {
		r := i + 3
		for j, text := range row {
			startCol, _ := excelize.ColumnNumberToName(infoSegments[j][0])
			endCol, _ := excelize.ColumnNumberToName(infoSegments[j][1])
			startCell := fmt.Sprintf("%s%d", startCol, r)
			endCell := fmt.Sprintf("%s%d", endCol, r)
			if startCell != endCell {
				_ = f.MergeCell(sheet, startCell, endCell)
			}
			_ = f.SetCellValue(sheet, startCell, text)
		}
		_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", r), fmt.Sprintf("%s%d", lastCol, r), infoStyle)
		_ = f.SetRowHeight(sheet, r, 22)
	}

	headerRow := 7
	_ = f.SetCellValue(sheet, "A7", "序号")
	_ = f.SetCellValue(sheet, "B7", "机身码(SN)")
	_ = f.SetCellValue(sheet, "C7", "检测结果")
	for i, item := range templateItems {
		col, _ := excelize.ColumnNumberToName(i + 4)
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, headerRow), item.Item.Name)
	}
	remarkCol, _ := excelize.ColumnNumberToName(len(templateItems) + 4)
	signCol, _ := excelize.ColumnNumberToName(len(templateItems) + 5)
	_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", remarkCol, headerRow), "备注")
	_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", signCol, headerRow), "签名")
	_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", headerRow), fmt.Sprintf("%s%d", signCol, headerRow), headerStyle)
	_ = f.SetRowHeight(sheet, headerRow, 36)

	minRows := len(devices)
	if minRows < 8 {
		minRows = 8
	}
	for i := 0; i < minRows; i++ {
		row := headerRow + 1 + i
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		if i < len(devices) {
			device := devices[i]
			_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), device.SN)
			_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), deviceInspectionResultLabel(device.Status))
			remarks := make([]string, 0)
			for j, item := range templateItems {
				result, ok := resultMap[exportResultKey{DeviceID: device.ID, ItemID: item.ItemID}]
				if !ok {
					continue
				}
				col, _ := excelize.ColumnNumberToName(j + 4)
				_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), exportResultText(item.Item.ResultType, result))
				if strings.TrimSpace(result.Remark) != "" {
					remarks = append(remarks, fmt.Sprintf("%s：%s", item.Item.Name, strings.TrimSpace(result.Remark)))
				}
			}
			if len(remarks) > 0 {
				_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", remarkCol, row), strings.Join(remarks, "；"))
			}
		}
		_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", signCol, row), cellStyle)
		_ = f.SetRowHeight(sheet, row, 28)
	}

	footerRow := headerRow + minRows + 2
	_ = f.MergeCell(sheet, fmt.Sprintf("A%d", footerRow), fmt.Sprintf("%s%d", signCol, footerRow))
	_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", footerRow), "仪器检测依据：企业标准：Q/440112014000AEFCHKJ 1-2025\n计量检定规程参照：JJG 1200-2023\n抽样检测参照：GB/T 2828.1-2021")
	_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", footerRow), fmt.Sprintf("%s%d", signCol, footerRow), footerStyle)
	_ = f.SetRowHeight(sheet, footerRow, 48)

	_ = f.SetColWidth(sheet, "A", "A", 6)
	_ = f.SetColWidth(sheet, "B", "B", 18)
	_ = f.SetColWidth(sheet, "C", "C", 12)
	for i := 0; i < len(templateItems); i++ {
		col, _ := excelize.ColumnNumberToName(i + 4)
		_ = f.SetColWidth(sheet, col, col, 10)
	}
	_ = f.SetColWidth(sheet, remarkCol, signCol, 14)
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      3,
		YSplit:      7,
		TopLeftCell: "D8",
		ActivePane:  "bottomRight",
	})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("%s-%s-检测工单.xlsx", order.MONumber, batch.BatchNumber)
	return buf, filename, nil
}

func thinBorders() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}

func businessCategoryLabel(value string) string {
	switch value {
	case "online":
		return "线上"
	case "offline":
		return "线下"
	case "foreign_trade":
		return "外贸"
	case "custom":
		return "定制款"
	default:
		return value
	}
}

func exportResultText(resultType string, result model.InspectionDeviceResult) string {
	parts := make([]string, 0, 2)
	if resultType != "number" {
		if result.PassResult != nil {
			if *result.PassResult {
				parts = append(parts, "√")
			} else {
				parts = append(parts, "×")
			}
		}
	}
	if resultType != "pass_fail" && result.NumberResult != nil {
		parts = append(parts, fmt.Sprintf("%g", *result.NumberResult))
	}
	return strings.Join(parts, " ")
}

func deviceInspectionResultLabel(status string) string {
	switch status {
	case "pass":
		return "合格"
	case "fail":
		return "不合格"
	case "pending":
		return "未完成"
	case "returned":
		return "待生产接收"
	case "rework":
		return "返工中"
	case "pending_recheck":
		return "待复检"
	case "rechecking":
		return "复检中"
	default:
		return ""
	}
}

func splitExcelSegments(totalCols int, parts int) [][2]int {
	if parts <= 0 {
		return nil
	}
	segments := make([][2]int, 0, parts)
	base := totalCols / parts
	remainder := totalCols % parts
	start := 1
	for i := 0; i < parts; i++ {
		width := base
		if i < remainder {
			width++
		}
		if width < 1 {
			width = 1
		}
		end := start + width - 1
		if end > totalCols {
			end = totalCols
		}
		segments = append(segments, [2]int{start, end})
		start = end + 1
	}
	return segments
}
