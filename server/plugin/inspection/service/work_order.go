package service

import (
	"errors"
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

func (s *workOrderSvc) SaveResults(req *request.SaveInspectionResult) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, ds := range req.DeviceStatuses {
			tx.Model(&model.ProductionOrderDevice{}).Where("id = ?", ds.DeviceID).Update("status", ds.Status)
		}
		for _, dr := range req.DeviceResults {
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
		return nil
	})
}

func (s *workOrderSvc) CompleteInspection(req *request.CompleteInspection) error {
	return global.GVA_DB.Model(&model.ProductionBatch{}).Where("id = ? AND status = 2", req.ID).Update("status", 3).Error
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
		list[i].DeviceCount = int(count)
		list[i].PassCount = int(pass)
		list[i].FailCount = int(fail)
		if list[i].TemplateID != nil {
			var tmpl model.InspectionTemplate
			if global.GVA_DB.Where("id = ?", *list[i].TemplateID).First(&tmpl).Error == nil {
				list[i].Template = &tmpl
			}
		}
	}

	return list, total, nil
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
		ID         uint                   `json:"ID"`
		SN         string                 `json:"sn"`
		LineNumber int                    `json:"lineNumber"`
		Status     string                 `json:"status"`
		Results    []InspectionResultItem `json:"results"`
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
		deviceInfos[i] = DeviceInfo{ID: d.ID, SN: d.SN, LineNumber: d.LineNumber, Status: d.Status, Results: make([]InspectionResultItem, len(templateItems))}
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
			"ID":                 batch.ID,
			"productionOrderID":  batch.ProductionOrderID,
			"moNumber":           order.MONumber,
			"batchNumber":        batch.BatchNumber,
			"productName":        func() string {
				if batch.Template != nil && batch.Template.ProductName != "" {
					return batch.Template.ProductName
				}
				return order.ProductName
			}(),
			"model":              order.Model,
			"firmwareVersion":    order.FirmwareVersion,
			"mainboardFirmwareVersion": order.MainboardFirmwareVersion,
			"pnCode":             order.PNCode,
			"instrumentCategory": order.InstrumentCategory,
			"status":             batch.Status,
			"templateID":         batch.TemplateID,
			"inspectorID":        batch.InspectorID,
			"inspectorName":      batch.InspectorName,
			"inspectionDate":     batch.InspectionDate,
		},
		"devices":       deviceInfos,
		"templateItems": templateItemInfos,
	}, nil
}
