package device

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DeviceCategoryApi
	DeviceModelApi
	FirmwareVersionApi
	ModelFirmwareRelApi
	FirmwareTagApi
	FirmwareVersionLogApi
	PublicFirmwareDownloadApi
}

var (
	deviceCategoryService        = service.ServiceGroupApp.DeviceServiceGroup.DeviceCategoryService
	deviceModelService           = service.ServiceGroupApp.DeviceServiceGroup.DeviceModelService
	firmwareVersionService       = service.ServiceGroupApp.DeviceServiceGroup.FirmwareVersionService
	modelFirmwareRelService      = service.ServiceGroupApp.DeviceServiceGroup.ModelFirmwareRelService
	firmwareTagService           = service.ServiceGroupApp.DeviceServiceGroup.FirmwareTagService
	firmwareVersionLogService    = service.ServiceGroupApp.DeviceServiceGroup.FirmwareVersionLogService
	firmwareVersionTagRelService = service.ServiceGroupApp.DeviceServiceGroup.FirmwareVersionTagRelService
)
