package service

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
)

var InspectionItem = new(inspectionItem)

type inspectionItem struct{}

func (s *inspectionItem) CreateItem(item *model.InspectionItem) error {
	return global.GVA_DB.Create(item).Error
}

func (s *inspectionItem) DeleteItem(id string) error {
	var count int64
	if err := global.GVA_DB.Table("inspection_template_items").Where("item_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该检测项被模板引用，无法删除")
	}
	return global.GVA_DB.Delete(&model.InspectionItem{}, "id = ?", id).Error
}

func (s *inspectionItem) DeleteItemByIds(ids []string) error {
	var count int64
	if err := global.GVA_DB.Table("inspection_template_items").Where("item_id in ?", ids).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("部分检测项被模板引用，无法批量删除")
	}
	return global.GVA_DB.Delete(&model.InspectionItem{}, "id in ?", ids).Error
}

func (s *inspectionItem) UpdateItem(item *model.InspectionItem) error {
	return global.GVA_DB.Model(&model.InspectionItem{}).Where("id = ?", item.ID).Updates(item).Error
}

func (s *inspectionItem) FindItem(id string) (model.InspectionItem, error) {
	var item model.InspectionItem
	err := global.GVA_DB.Where("id = ?", id).First(&item).Error
	return item, err
}

func (s *inspectionItem) GetItemList(search request.InspectionItemSearch) (list []model.InspectionItem, total int64, err error) {
	db := global.GVA_DB.Model(&model.InspectionItem{})
	if search.Name != "" {
		db = db.Where("name LIKE ?", "%"+search.Name+"%")
	}
	if search.ResultType != "" {
		db = db.Where("result_type = ?", search.ResultType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Order("id desc").Find(&list).Error
	return list, total, err
}
