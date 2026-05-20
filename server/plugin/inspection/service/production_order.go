package service

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"gorm.io/gorm"
)

var ProductionOrder = new(productionOrderSvc)

type productionOrderSvc struct{}

func (s *productionOrderSvc) CreateProductionOrder(req *request.CreateProductionOrder) error {
	now := time.Now()
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		po := model.ProductionOrder{
			MONumber:           req.MONumber,
			TemplateID:         req.TemplateID,
			ProductName:        req.ProductName,
			Model:              req.Model,
			FirmwareVersion:    req.FirmwareVersion,
			InstrumentCategory: req.InstrumentCategory,
			Status:             0,
			SubmitDate:         &now,
			Remark:             req.Remark,
		}
		if err := tx.Create(&po).Error; err != nil {
			return err
		}
		for i, sn := range req.SNs {
			if sn == "" {
				continue
			}
			device := model.ProductionOrderDevice{
				ProductionOrderID: po.ID,
				SN:                sn,
				LineNumber:        i + 1,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *productionOrderSvc) DeleteProductionOrder(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var po model.ProductionOrder
		if err := tx.Where("id = ?", id).First(&po).Error; err != nil {
			return err
		}
		if po.Status >= 1 {
			return errors.New("已确认的订单不允许删除")
		}
		if err := tx.Where("production_order_id = ?", id).Delete(&model.ProductionBatch{}).Error; err != nil {
			return err
		}
		return tx.Where("production_order_id = ?", id).Delete(&model.ProductionOrderDevice{}).Error
	})
}

func (s *productionOrderSvc) UpdateProductionOrder(req *request.UpdateProductionOrder) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"mo_number":           req.MONumber,
			"template_id":         req.TemplateID,
			"product_name":        req.ProductName,
			"model":               req.Model,
			"firmware_version":    req.FirmwareVersion,
			"instrument_category": req.InstrumentCategory,
			"remark":              req.Remark,
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		return tx.Model(&model.ProductionOrder{}).Where("id = ?", req.ID).Updates(updates).Error
	})
}

func (s *productionOrderSvc) FindProductionOrder(id string) (model.ProductionOrder, error) {
	var po model.ProductionOrder
	err := global.GVA_DB.Preload("Template").Preload("Batches.Devices").Where("id = ?", id).First(&po).Error
	if err != nil {
		return po, err
	}
	var unbatched []model.ProductionOrderDevice
	global.GVA_DB.Where("production_order_id = ? AND batch_id IS NULL", id).Order("line_number asc").Find(&unbatched)
	po.DeviceCount = len(unbatched)
	for _, b := range po.Batches {
		po.DeviceCount += len(b.Devices)
	}
	var pass, fail int64
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = 'pass'", id).Count(&pass)
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = 'fail'", id).Count(&fail)
	po.PassCount = int(pass)
	po.FailCount = int(fail)
	return po, err
}

func (s *productionOrderSvc) GetProductionOrderList(search request.ProductionOrderSearch) (list []model.ProductionOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.ProductionOrder{})
	if search.MONumber != "" {
		db = db.Where("mo_number LIKE ?", "%"+search.MONumber+"%")
	}
	if search.Model != "" {
		db = db.Where("model LIKE ?", "%"+search.Model+"%")
	}
	if search.InstrumentCategory != "" {
		db = db.Where("instrument_category = ?", search.InstrumentCategory)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Order("id desc").Find(&list).Error
	for i := range list {
		var count, pass, fail int64
		deviceDB := global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ?", list[i].ID)
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
	return list, total, err
}

// SubmitDeviceData 生产工具提交全量数据
func (s *productionOrderSvc) SubmitDeviceData(req *request.SubmitDeviceData, submitterID uint, submitterName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// Find or create production order
		var po model.ProductionOrder
		err := tx.Where("mo_number = ?", req.MONumber).First(&po).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				now := time.Now()
				po = model.ProductionOrder{
					MONumber:           req.MONumber,
					ProductName:        req.DeviceType,
					Model:              req.DeviceType,
					InstrumentCategory: req.InstrumentCategory,
					Status:             0,
					SubmitterID:        &submitterID,
					SubmitterName:      submitterName,
					SubmitDate:         &now,
				}
				if err := tx.Create(&po).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Create batch if provided
		var batchID *uint
		if req.BatchNumber != "" {
			var batch model.ProductionBatch
			err := tx.Where("production_order_id = ? AND batch_number = ?", po.ID, req.BatchNumber).First(&batch).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				batch = model.ProductionBatch{
					ProductionOrderID: po.ID,
					BatchNumber:       req.BatchNumber,
					Status:            0,
				}
				if err := tx.Create(&batch).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			batchID = &batch.ID
		}

		// Insert devices
		for i, sn := range req.SNs {
			if sn == "" {
				continue
			}
			var existing model.ProductionOrderDevice
			if err := tx.Where("sn = ?", sn).First(&existing).Error; err == nil {
				continue // Skip duplicates
			}
			device := model.ProductionOrderDevice{
				ProductionOrderID: po.ID,
				BatchID:           batchID,
				SN:                sn,
				Model:             req.DeviceType,
				PNCode:            req.PNCode,
				FirmwareVersion:   req.DeviceInfo, // will be parsed from deviceInfo JSON if needed
				TimeLicense:       req.TimeLicense,
				RegionLicense:     req.RegionLicense,
				NtripCode:         req.NtripCode,
				LineNumber:        i + 1,
				DeviceInfo:        req.DeviceInfo,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// AssignBatch 分配设备到批次
func (s *productionOrderSvc) AssignBatch(req *request.AssignBatch) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return errors.New("批次不存在")
		}
		for _, sn := range req.SNs {
			var device model.ProductionOrderDevice
			if err := tx.Where("sn = ?", sn).First(&device).Error; err != nil {
				return errors.New("SN不存在: " + sn)
			}
			if device.ProductionOrderID != batch.ProductionOrderID {
				return errors.New("SN不属于该生产号: " + sn)
			}
			tx.Model(&device).Update("batch_id", req.BatchID)
		}
		return nil
	})
}
