package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FirmwareVersionRouter struct{}

// InitFirmwareVersionRouter 初始化固件版本路由
func (r *FirmwareVersionRouter) InitFirmwareVersionRouter(Router *gin.RouterGroup) {
	firmwareVersionRouter := Router.Group("firmwareVersion").Use(middleware.OperationRecord())
	firmwareVersionRouterWithoutRecord := Router.Group("firmwareVersion")
	{
		firmwareVersionRouter.POST("createFirmwareVersion", firmwareVersionApi.CreateFirmwareVersion)             // 创建固件版本
		firmwareVersionRouter.DELETE("deleteFirmwareVersion", firmwareVersionApi.DeleteFirmwareVersion)           // 删除固件版本
		firmwareVersionRouter.DELETE("deleteFirmwareVersionByIds", firmwareVersionApi.DeleteFirmwareVersionByIds) // 批量删除固件版本
		firmwareVersionRouter.PUT("updateFirmwareVersion", firmwareVersionApi.UpdateFirmwareVersion)              // 更新固件版本
		firmwareVersionRouter.POST("changeFirmwareVersionStatus", firmwareVersionApi.ChangeFirmwareVersionStatus) // 更新固件版本状态
		firmwareVersionRouter.POST("publishFirmwareVersion", firmwareVersionApi.PublishFirmwareVersion)           // 发布固件版本
		firmwareVersionRouter.POST("setFirmwareStable", firmwareVersionApi.SetFirmwareStable)                     // 设置稳定版本
		firmwareVersionRouter.POST("voidFirmwareVersion", firmwareVersionApi.VoidFirmwareVersion)                 // 作废固件版本
		firmwareVersionRouter.POST("deleteFirmwarePackage", firmwareVersionApi.DeleteFirmwarePackage)             // 删除固件包
	}
	{
		firmwareVersionRouterWithoutRecord.GET("findFirmwareVersion", firmwareVersionApi.FindFirmwareVersion)       // 获取固件版本详情
		firmwareVersionRouterWithoutRecord.GET("getFirmwareVersionList", firmwareVersionApi.GetFirmwareVersionList) // 获取固件版本列表
	}
}
