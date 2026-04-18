package device

import "github.com/gin-gonic/gin"

type PublicFirmwareDownloadRouter struct{}

// InitPublicFirmwareDownloadRouter 初始化公开固件下载路由
func (r *PublicFirmwareDownloadRouter) InitPublicFirmwareDownloadRouter(Router *gin.RouterGroup) {
	publicFirmwareRouter := Router.Group("firmwarePublic")
	{
		publicFirmwareRouter.GET(
			"getPublicFirmwareDownloadPage",
			publicFirmwareApi.GetPublicFirmwareDownloadPage,
		)
	}
}
