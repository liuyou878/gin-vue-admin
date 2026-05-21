package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"gorm.io/gorm"
)

var Template = new(templateSvc)

type templateSvc struct{}

func (s *templateSvc) CreateTemplate(req *request.CreateTemplate) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		tmpl := model.InspectionTemplate{
			Name:            req.Name,
			ProductName:     req.ProductName,
			Model:           req.Model,
			Status:          1,
		}
		if err := tx.Create(&tmpl).Error; err != nil {
			return err
		}
		for _, it := range req.Items {
			ti := model.InspectionTemplateItem{
				TemplateID: tmpl.ID,
				ItemID:     it.ItemID,
				Sort:       it.Sort,
			}
			if err := tx.Create(&ti).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *templateSvc) DeleteTemplate(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("template_id = ?", id).Delete(&model.InspectionTemplateItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.InspectionTemplate{}, "id = ?", id).Error
	})
}

func (s *templateSvc) UpdateTemplate(req *request.UpdateTemplate) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"name":         req.Name,
			"product_name": req.ProductName,
			"model":        req.Model,
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		if err := tx.Model(&model.InspectionTemplate{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Where("template_id = ?", req.ID).Delete(&model.InspectionTemplateItem{}).Error; err != nil {
			return err
		}
		for _, it := range req.Items {
			ti := model.InspectionTemplateItem{
				TemplateID: req.ID,
				ItemID:     it.ItemID,
				Sort:       it.Sort,
			}
			if err := tx.Create(&ti).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *templateSvc) FindTemplate(id string) (model.InspectionTemplate, error) {
	var tmpl model.InspectionTemplate
	err := global.GVA_DB.Where("id = ?", id).First(&tmpl).Error
	if err != nil {
		return tmpl, err
	}
	var items []model.InspectionTemplateItem
	err = global.GVA_DB.Preload("Item").Where("template_id = ?", id).Order("sort asc").Find(&items).Error
	tmpl.TemplateItems = items
	tmpl.ItemCount = len(items)
	return tmpl, err
}

func (s *templateSvc) GetTemplateList(search request.TemplateSearch) (list []model.InspectionTemplate, total int64, err error) {
	db := global.GVA_DB.Model(&model.InspectionTemplate{})
	if search.Name != "" {
		db = db.Where("name LIKE ?", "%"+search.Name+"%")
	}
	if search.Model != "" {
		db = db.Where("model LIKE ?", "%"+search.Model+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Order("id desc").Find(&list).Error
	// Populate item counts
	for i := range list {
		var count int64
		global.GVA_DB.Model(&model.InspectionTemplateItem{}).Where("template_id = ?", list[i].ID).Count(&count)
		list[i].ItemCount = int(count)
	}
	return list, total, err
}
