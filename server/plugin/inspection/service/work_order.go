package service

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"gorm.io/gorm"
)

var WorkOrder = new(workOrderSvc)

type workOrderSvc struct{}

func (s *workOrderSvc) StartInspection(req *request.StartInspection, inspectorID uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.ProductionOrder{}).Where("id = ? AND status = 1", req.ID).Updates(map[string]interface{}{
		"status":          2,
		"inspector_id":    inspectorID,
		"inspection_date": &now,
	}).Error
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
	return global.GVA_DB.Model(&model.ProductionOrder{}).Where("id = ? AND status = 2", req.ID).Update("status", 3).Error
}

// GetInspectionDetail returns the full inspection data for a production order
func (s *workOrderSvc) GetInspectionDetail(orderID string) (model.ProductionOrder, error) {
	var po model.ProductionOrder
	err := global.GVA_DB.Preload("Template").Where("id = ?", orderID).First(&po).Error
	if err != nil {
		return po, err
	}

	var devices []model.ProductionOrderDevice
	err = global.GVA_DB.Where("production_order_id = ?", orderID).Order("line_number asc").Find(&devices).Error
	if err != nil {
		return po, err
	}

	// Load template items if template exists
	var templateItems []model.InspectionTemplateItem
	if po.TemplateID != nil {
		err = global.GVA_DB.Preload("Item").Where("template_id = ?", *po.TemplateID).Order("sort asc").Find(&templateItems).Error
		if err != nil {
			return po, err
		}
	}

	// Load results for each device
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

	// Build device list with results
	type DeviceWithResults struct {
		model.ProductionOrderDevice
		Results []InspectionResultItem `json:"results"`
	}
	devicesWithResults := make([]DeviceWithResults, len(devices))
	for i, d := range devices {
		devicesWithResults[i].ProductionOrderDevice = d
		devicesWithResults[i].Results = make([]InspectionResultItem, len(templateItems))
		for j, ti := range templateItems {
			key := ResultKey{d.ID, ti.ItemID}
			r := resultMap[key]
			devicesWithResults[i].Results[j] = InspectionResultItem{
				ItemID:       ti.ItemID,
				ItemName:     ti.Item.Name,
				ResultType:   ti.Item.ResultType,
				Unit:         ti.Item.Unit,
				MinValue:     ti.Item.MinValue,
				MaxValue:     ti.Item.MaxValue,
				PassResult:   r.PassResult,
				NumberResult: r.NumberResult,
				Remark:       r.Remark,
			}
		}
	}

	po.Devices = devices
	po.DeviceCount = len(devices)

	// Store template items and device results as extra data for JSON serialization
	// We'll use a map to pass extra data
	return po, err
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

// GetInspectionDetailData returns full detail with devices+results for frontend
func (s *workOrderSvc) GetInspectionDetailData(orderID string) (map[string]interface{}, error) {
	var po model.ProductionOrder
	err := global.GVA_DB.Preload("Template").Where("id = ?", orderID).First(&po).Error
	if err != nil {
		return nil, err
	}

	var devices []model.ProductionOrderDevice
	global.GVA_DB.Where("production_order_id = ?", orderID).Order("line_number asc").Find(&devices)

	var templateItems []model.InspectionTemplateItem
	if po.TemplateID != nil {
		global.GVA_DB.Preload("Item").Where("template_id = ?", *po.TemplateID).Order("sort asc").Find(&templateItems)
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
		"order":         po,
		"devices":       deviceInfos,
		"templateItems": templateItemInfos,
	}, nil
}
