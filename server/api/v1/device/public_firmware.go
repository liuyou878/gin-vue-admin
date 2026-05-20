package device

import (
	"net/http"

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

// DownloadFirmwarePackage 公开下载已发布固件包
func (a *PublicFirmwareDownloadApi) DownloadFirmwarePackage(c *gin.Context) {
	var req deviceReq.DownloadFirmwarePackageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.FirmwareID == 0 {
		response.FailWithMessage("固件版本ID不能为空", c)
		return
	}
	download, err := firmwareVersionService.OpenPublicFirmwarePackageDownload(req.FirmwareID)
	if err != nil {
		global.GVA_LOG.Error("公开下载固件包失败!", zap.Error(err))
		response.FailWithMessage("公开下载固件包失败:"+err.Error(), c)
		return
	}
	defer download.Close()

	c.Header("Content-Disposition", buildAttachmentDisposition(download.FileName))
	c.Header("Content-Type", download.ContentType)
	if download.Size > 0 {
		c.Header("Content-Length", int64ToString(download.Size))
	}
	c.DataFromReader(http.StatusOK, download.Size, download.ContentType, download.Reader, nil)
}

// DownloadDeveloperLog 公开下载固件更新日志
func (a *PublicFirmwareDownloadApi) DownloadDeveloperLog(c *gin.Context) {
	var req deviceReq.DownloadDeveloperLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.FirmwareID == 0 {
		response.FailWithMessage("固件版本ID不能为空", c)
		return
	}
	download, err := firmwareVersionService.OpenPublicDeveloperLogDownload(req.FirmwareID)
	if err != nil {
		global.GVA_LOG.Error("公开下载更新日志失败!", zap.Error(err))
		response.FailWithMessage("公开下载更新日志失败:"+err.Error(), c)
		return
	}
	defer download.Close()

	c.Header("Content-Disposition", buildAttachmentDisposition(download.FileName))
	c.Header("Content-Type", download.ContentType)
	if download.Size > 0 {
		c.Header("Content-Length", int64ToString(download.Size))
	}
	c.DataFromReader(http.StatusOK, download.Size, download.ContentType, download.Reader, nil)
}
