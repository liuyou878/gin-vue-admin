package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DeviceModelRouter struct{}

// InitDeviceModelRouter 初始化设备型号路由
func (r *DeviceModelRouter) InitDeviceModelRouter(Router *gin.RouterGroup) {
	deviceModelRouter := Router.Group("deviceModel").Use(middleware.OperationRecord())
	deviceModelRouterWithoutRecord := Router.Group("deviceModel")
	{
		deviceModelRouter.POST("createDeviceModel", deviceModelApi.CreateDeviceModel)             // 创建设备型号
		deviceModelRouter.DELETE("deleteDeviceModel", deviceModelApi.DeleteDeviceModel)           // 删除设备型号
		deviceModelRouter.DELETE("deleteDeviceModelByIds", deviceModelApi.DeleteDeviceModelByIds) // 批量删除设备型号
		deviceModelRouter.PUT("updateDeviceModel", deviceModelApi.UpdateDeviceModel)              // 更新设备型号
	}
	{
		deviceModelRouterWithoutRecord.GET("findDeviceModel", deviceModelApi.FindDeviceModel)       // 获取设备型号详情
		deviceModelRouterWithoutRecord.GET("getDeviceModelList", deviceModelApi.GetDeviceModelList) // 获取设备型号列表
	}
}
