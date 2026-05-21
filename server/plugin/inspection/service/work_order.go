package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"gorm.io/gorm"
)

var WorkOrder = new(workOrderSvc)

type workOrderSvc struct{}

func (s *workOrderSvc) StartInspection(req *request.StartInspection, inspectorID uint, inspectorName string) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.ProductionBatch{}).Where("id = ? AND status = 1", req.ID).Updates(map[string]interface{}{
		"status":          2,
		"inspector_id":    inspectorID,
		"inspector_name":  inspectorName,
		"inspection_date": &now,
	}).Error
}

func (s *workOrderSvc) StartRecheck(req *request.StartRecheck, inspectorID uint, inspectorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ? AND status = 3", req.ID).First(&batch).Error; err != nil {
			return err
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

func (s *workOrderSvc) AssignBatchTemplate(req *request.AssignBatchTemplate) error {
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
			"status":      1,
		}).Error; err != nil {
			return err
		}

		return tx.Model(&model.ProductionOrder{}).
			Where("id = ? AND status = 0", batch.ProductionOrderID).
			Update("status", 1).Error
	})
}

func (s *workOrderSvc) AssignOrderTemplate(req *request.AssignOrderTemplate) error {
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
			Updates(map[string]interface{}{
				"template_id": req.TemplateID,
				"status":      1,
			}).Error; err != nil {
			return err
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

func (s *workOrderSvc) SaveResults(req *request.SaveInspectionResult) error {
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
			if batch.Status == 3 && device.Status != "rechecking" {
				return fmt.Errorf("设备 %s 当前不是复检中，不能修改检测结果", device.SN)
			}
			var existing model.InspectionDeviceResult
			err := tx.Where("production_order_device_id = ? AND item_id = ?", dr.DeviceID, dr.ItemID).First(&existing).Error
			if err == nil {
				tx.Model(&existing).Updates(map[string]interface{}{
					"pass_result":   dr.PassResult,
					"number_result": dr.NumberResult,
					"remark":        dr.Remark,
				})
			} else {
				tx.Create(&model.InspectionDeviceResult{
					ProductionOrderDeviceID: dr.DeviceID,
					ItemID:                  dr.ItemID,
					PassResult:              dr.PassResult,
					NumberResult:            dr.NumberResult,
					Remark:                  dr.Remark,
				})
			}
		}
		for _, ds := range req.DeviceStatuses {
			var device model.ProductionOrderDevice
			if err := tx.Where("id = ? AND batch_id = ?", ds.DeviceID, req.BatchID).First(&device).Error; err != nil {
				return err
			}
			if batch.Status == 3 && device.Status != "rechecking" {
				return fmt.Errorf("设备 %s 当前不是复检中，不能修改检测结果", device.SN)
			}
			if err := updateDeviceStatusWithLog(tx, device, ds.Status, "保存检测结果", nil, ""); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *workOrderSvc) ReturnDevices(req *request.ReturnDevices, returnByID uint, returnByName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return err
		}
		if batch.Status != 3 {
			return errors.New("只有检测完成后的不合格设备才可以打回生产")
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
			if err := updateDeviceStatusWithLog(tx, device, "rework", req.Reason, returnByID, returnByName); err != nil {
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

func (s *workOrderSvc) CompleteInspection(req *request.CompleteInspection) error {
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

		return tx.Model(&model.ProductionBatch{}).Where("id = ?", req.ID).Update("status", 3).Error
	})
}

func (s *workOrderSvc) CompleteRecheck(req *request.CompleteRecheck) error {
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
			for _, templateItem := range templateItems {
				result, ok := resultMap[resultKey{DeviceID: device.ID, ItemID: templateItem.ItemID}]
				if !ok || !inspectionResultCompleted(templateItem.Item.ResultType, result) {
					return fmt.Errorf("设备 %s 的复检项 %s 未完成", device.SN, templateItem.Item.Name)
				}
			}
			if device.Status == "rechecking" {
				return fmt.Errorf("设备 %s 还没有复检判定，请先保存复检结果", device.SN)
			}
		}

		return nil
	})
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

type InspectionResultItem struct {
	ItemID       uint     `json:"itemID"`
	ItemName     string   `json:"itemName"`
	ResultType   string   `json:"resultType"`
	Unit         string   `json:"unit"`
	MinValue     *float64 `json:"minValue"`
	MaxValue     *float64 `json:"maxValue"`
	PassResult   *bool    `json:"passResult"`
	NumberResult *float64 `json:"numberResult"`
	Remark       string   `json:"remark"`
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
		deviceDB := global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("batch_id = ?", list[i].ID)
		deviceDB.Count(&count)
		deviceDB.Where("status = 'pass'").Count(&pass)
		deviceDB.Where("status = 'fail'").Count(&fail)
		var rework int64
		deviceDB.Where("status = 'rework'").Count(&rework)
		var recheck int64
		deviceDB.Where("status IN ?", []string{"pending_recheck", "rechecking"}).Count(&recheck)
		var rechecking int64
		deviceDB.Where("status = 'rechecking'").Count(&rechecking)
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
			"inspectorID":              batch.InspectorID,
			"inspectorName":            batch.InspectorName,
			"inspectionDate":           batch.InspectionDate,
		},
		"devices":       deviceInfos,
		"templateItems": templateItemInfos,
	}, nil
}
