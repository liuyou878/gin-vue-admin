package device

import "github.com/gin-gonic/gin"

type FirmwareVersionLogRouter struct{}

// InitFirmwareVersionLogRouter 初始化固件日志路由
func (r *FirmwareVersionLogRouter) InitFirmwareVersionLogRouter(PrivateRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	firmwareVersionLogRouter := PrivateRouter.Group("firmwareVersionLog")
	publicFirmwareVersionLogRouter := PublicRouter.Group("firmwareVersionLog")
	{
		firmwareVersionLogRouter.GET("findFirmwareVersionLog", firmwareVersionLogApi.FindFirmwareVersionLog)       // 获取固件日志详情
		firmwareVersionLogRouter.GET("getFirmwareVersionLogList", firmwareVersionLogApi.GetFirmwareVersionLogList) // 获取固件日志列表
	}
	{
		publicFirmwareVersionLogRouter.GET("downloadFirmwareLogPackage", firmwareVersionLogApi.DownloadFirmwareLogPackage) // 下载固件日志包
	}
}
