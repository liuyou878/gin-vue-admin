package device

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
)

// CreateDeviceModel 创建设备型号
func (s *DeviceModelService) CreateDeviceModel(model *deviceModel.DeviceModel) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.DeviceCategory{}).Where("id = ?", model.CategoryID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("设备类别不存在")
	}
	return global.GVA_DB.Create(model).Error
}

// DeleteDeviceModel 删除设备型号
func (s *DeviceModelService) DeleteDeviceModel(id string) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.ModelFirmwareRel{}).Where("model_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该设备型号已关联固件版本，不能删除")
	}
	return global.GVA_DB.Delete(&deviceModel.DeviceModel{}, "id = ?", id).Error
}

// DeleteDeviceModelByIds 批量删除设备型号
func (s *DeviceModelService) DeleteDeviceModelByIds(ids commonReq.IdsReq) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.ModelFirmwareRel{}).Where("model_id in ?", ids.Ids).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("所选设备型号中已有关联固件版本，不能删除")
	}
	return global.GVA_DB.Delete(&[]deviceModel.DeviceModel{}, "id in ?", ids.Ids).Error
}

// UpdateDeviceModel 更新设备型号
func (s *DeviceModelService) UpdateDeviceModel(model deviceModel.DeviceModel) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.DeviceCategory{}).Where("id = ?", model.CategoryID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("设备类别不存在")
	}
	return global.GVA_DB.Model(&deviceModel.DeviceModel{}).
		Where("id = ?", model.ID).
		Select("category_id", "model_code", "model_name", "sort", "series_name", "generation", "status", "remark").
		Updates(&model).Error
}

// GetDeviceModel 获取设备型号详情
func (s *DeviceModelService) GetDeviceModel(id string) (model deviceModel.DeviceModel, err error) {
	err = global.GVA_DB.Preload("Category").Where("id = ?", id).First(&model).Error
	return
}

// GetDeviceModelInfoList 获取设备型号分页列表
func (s *DeviceModelService) GetDeviceModelInfoList(info deviceReq.DeviceModelSearch) (list []deviceModel.DeviceModel, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.DeviceModel{})
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	if info.ModelCode != "" {
		db = db.Where("model_code LIKE ?", "%"+info.ModelCode+"%")
	}
	if info.ModelName != "" {
		db = db.Where("model_name LIKE ?", "%"+info.ModelName+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}
	err = db.Preload("Category").Order("sort asc, created_at desc, id desc").Find(&list).Error
	return
}
