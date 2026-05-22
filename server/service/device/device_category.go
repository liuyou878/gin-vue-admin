package device

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
)

// CreateDeviceCategory 创建设备类别
func (s *DeviceCategoryService) CreateDeviceCategory(category *deviceModel.DeviceCategory) error {
	return global.GVA_DB.Create(category).Error
}

// DeleteDeviceCategory 删除设备类别
func (s *DeviceCategoryService) DeleteDeviceCategory(id string) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.DeviceModel{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该设备类别下仍有关联型号，不能删除")
	}
	return global.GVA_DB.Delete(&deviceModel.DeviceCategory{}, "id = ?", id).Error
}

// DeleteDeviceCategoryByIds 批量删除设备类别
func (s *DeviceCategoryService) DeleteDeviceCategoryByIds(ids commonReq.IdsReq) error {
	var count int64
	if err := global.GVA_DB.Model(&deviceModel.DeviceModel{}).Where("category_id in ?", ids.Ids).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("所选设备类别中仍有关联型号，不能删除")
	}
	return global.GVA_DB.Delete(&[]deviceModel.DeviceCategory{}, "id in ?", ids.Ids).Error
}

// UpdateDeviceCategory 更新设备类别
func (s *DeviceCategoryService) UpdateDeviceCategory(category deviceModel.DeviceCategory) error {
	return global.GVA_DB.Model(&deviceModel.DeviceCategory{}).
		Where("id = ?", category.ID).
		Select("name", "code", "sort", "status", "remark").
		Updates(&category).Error
}

// GetDeviceCategory 获取设备类别详情
func (s *DeviceCategoryService) GetDeviceCategory(id string) (category deviceModel.DeviceCategory, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&category).Error
	return
}

// GetDeviceCategoryInfoList 获取设备类别分页列表
func (s *DeviceCategoryService) GetDeviceCategoryInfoList(info deviceReq.DeviceCategorySearch) (list []deviceModel.DeviceCategory, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.DeviceCategory{})
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Code != "" {
		db = db.Where("code LIKE ?", "%"+info.Code+"%")
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
	err = db.Order("sort asc, created_at desc, id desc").Find(&list).Error
	return
}
