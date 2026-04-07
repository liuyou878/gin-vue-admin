package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
)

// CreateFirmwareVersionLog 创建固件日志
func (s *FirmwareVersionLogService) CreateFirmwareVersionLog(log *deviceModel.FirmwareVersionLog) error {
	return global.GVA_DB.Create(log).Error
}

// GetFirmwareVersionLog 获取固件日志详情
func (s *FirmwareVersionLogService) GetFirmwareVersionLog(id string) (log deviceModel.FirmwareVersionLog, err error) {
	err = global.GVA_DB.Preload("Firmware").Preload("Model").Where("id = ?", id).First(&log).Error
	return
}

// GetFirmwareVersionLogInfoList 获取固件日志分页列表
func (s *FirmwareVersionLogService) GetFirmwareVersionLogInfoList(info deviceReq.FirmwareVersionLogSearch) (list []deviceModel.FirmwareVersionLog, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.FirmwareVersionLog{})
	if info.FirmwareID > 0 {
		db = db.Where("firmware_id = ?", info.FirmwareID)
	}
	if info.ModelID > 0 {
		db = db.Where("model_id = ?", info.ModelID)
	}
	if info.Action != "" {
		db = db.Where("action = ?", info.Action)
	}
	if len(info.OperateAtRange) == 2 {
		db = db.Where("operate_at BETWEEN ? AND ?", info.OperateAtRange[0], info.OperateAtRange[1])
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}
	err = db.Preload("Firmware").Preload("Model").Order("operate_at desc, id desc").Find(&list).Error
	return
}
