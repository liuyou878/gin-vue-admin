package device

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	DeviceCategoryRouter
	DeviceModelRouter
	FirmwareVersionRouter
	ModelFirmwareRelRouter
	FirmwareTagRouter
	FirmwareVersionLogRouter
}

var (
	deviceCategoryApi     = api.ApiGroupApp.DeviceApiGroup.DeviceCategoryApi
	deviceModelApi        = api.ApiGroupApp.DeviceApiGroup.DeviceModelApi
	firmwareVersionApi    = api.ApiGroupApp.DeviceApiGroup.FirmwareVersionApi
	modelFirmwareRelApi   = api.ApiGroupApp.DeviceApiGroup.ModelFirmwareRelApi
	firmwareTagApi        = api.ApiGroupApp.DeviceApiGroup.FirmwareTagApi
	firmwareVersionLogApi = api.ApiGroupApp.DeviceApiGroup.FirmwareVersionLogApi
)
