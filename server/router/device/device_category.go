package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DeviceCategoryRouter struct{}

// InitDeviceCategoryRouter 初始化设备类别路由
func (r *DeviceCategoryRouter) InitDeviceCategoryRouter(Router *gin.RouterGroup) {
	deviceCategoryRouter := Router.Group("deviceCategory").Use(middleware.OperationRecord())
	deviceCategoryRouterWithoutRecord := Router.Group("deviceCategory")
	{
		deviceCategoryRouter.POST("createDeviceCategory", deviceCategoryApi.CreateDeviceCategory)             // 创建设备类别
		deviceCategoryRouter.DELETE("deleteDeviceCategory", deviceCategoryApi.DeleteDeviceCategory)           // 删除设备类别
		deviceCategoryRouter.DELETE("deleteDeviceCategoryByIds", deviceCategoryApi.DeleteDeviceCategoryByIds) // 批量删除设备类别
		deviceCategoryRouter.PUT("updateDeviceCategory", deviceCategoryApi.UpdateDeviceCategory)              // 更新设备类别
	}
	{
		deviceCategoryRouterWithoutRecord.GET("findDeviceCategory", deviceCategoryApi.FindDeviceCategory)       // 获取设备类别详情
		deviceCategoryRouterWithoutRecord.GET("getDeviceCategoryList", deviceCategoryApi.GetDeviceCategoryList) // 获取设备类别列表
	}
}
