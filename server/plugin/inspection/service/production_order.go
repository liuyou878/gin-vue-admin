package service

import (
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
			BatchNumber:        req.BatchNumber,
			Status:             1,
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
		if err := tx.Where("production_order_id = ?", id).Delete(&model.ProductionOrderDevice{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ProductionOrder{}, "id = ?", id).Error
	})
}

func (s *productionOrderSvc) UpdateProductionOrder(req *request.UpdateProductionOrder) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"mo_number":            req.MONumber,
			"template_id":         req.TemplateID,
			"product_name":         req.ProductName,
			"model":                req.Model,
			"firmware_version":     req.FirmwareVersion,
			"instrument_category":  req.InstrumentCategory,
			"batch_number":         req.BatchNumber,
			"remark":               req.Remark,
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		if err := tx.Model(&model.ProductionOrder{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Where("production_order_id = ?", req.ID).Delete(&model.ProductionOrderDevice{}).Error; err != nil {
			return err
		}
		for i, sn := range req.SNs {
			if sn == "" {
				continue
			}
			device := model.ProductionOrderDevice{
				ProductionOrderID: req.ID,
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

func (s *productionOrderSvc) FindProductionOrder(id string) (model.ProductionOrder, error) {
	var po model.ProductionOrder
	err := global.GVA_DB.Preload("Template").Where("id = ?", id).First(&po).Error
	if err != nil {
		return po, err
	}
	var devices []model.ProductionOrderDevice
	err = global.GVA_DB.Where("production_order_id = ?", id).Order("line_number asc").Find(&devices).Error
	po.Devices = devices
	po.DeviceCount = len(devices)
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
		var count int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ?", list[i].ID).Count(&count)
		list[i].DeviceCount = int(count)
		// Preload template name
		if list[i].TemplateID != nil {
			var tmpl model.InspectionTemplate
			if global.GVA_DB.Where("id = ?", *list[i].TemplateID).First(&tmpl).Error == nil {
				list[i].Template = &tmpl
			}
		}
	}
	return list, total, err
}
