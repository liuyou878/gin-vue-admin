package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FirmwareTagRouter struct{}

// InitFirmwareTagRouter 初始化固件标签路由
func (r *FirmwareTagRouter) InitFirmwareTagRouter(Router *gin.RouterGroup) {
	firmwareTagRouter := Router.Group("firmwareTag").Use(middleware.OperationRecord())
	firmwareTagRouterWithoutRecord := Router.Group("firmwareTag")
	{
		firmwareTagRouter.POST("createFirmwareTag", firmwareTagApi.CreateFirmwareTag)             // 新建固件标签
		firmwareTagRouter.DELETE("deleteFirmwareTag", firmwareTagApi.DeleteFirmwareTag)           // 删除固件标签
		firmwareTagRouter.DELETE("deleteFirmwareTagByIds", firmwareTagApi.DeleteFirmwareTagByIds) // 批量删除固件标签
		firmwareTagRouter.PUT("updateFirmwareTag", firmwareTagApi.UpdateFirmwareTag)              // 更新固件标签
		firmwareTagRouter.POST("setFirmwareTags", firmwareTagApi.SetFirmwareTags)                 // 设置固件版本标签
	}
	{
		firmwareTagRouterWithoutRecord.GET("findFirmwareTag", firmwareTagApi.FindFirmwareTag)       // 获取固件标签详情
		firmwareTagRouterWithoutRecord.GET("getFirmwareTagList", firmwareTagApi.GetFirmwareTagList) // 获取固件标签列表
	}
}
