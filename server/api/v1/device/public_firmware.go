package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PublicFirmwareDownloadApi struct{}

// GetPublicFirmwareDownloadPage 获取公开固件下载页
func (a *PublicFirmwareDownloadApi) GetPublicFirmwareDownloadPage(c *gin.Context) {
	var req deviceReq.PublicFirmwareDownloadSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageData, err := firmwareVersionService.GetPublicFirmwareDownloadPage(req.CategoryID, req.ModelID)
	if err != nil {
		global.GVA_LOG.Error("获取公开固件下载页失败!", zap.Error(err))
		response.FailWithMessage("获取公开固件下载页失败:"+err.Error(), c)
		return
	}
	response.OkWithData(pageData, c)
}
