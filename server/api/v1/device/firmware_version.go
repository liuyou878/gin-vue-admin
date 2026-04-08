package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FirmwareVersionApi struct{}

// CreateFirmwareVersion 创建固件版本
func (a *FirmwareVersionApi) CreateFirmwareVersion(c *gin.Context) {
	var firmware deviceModel.FirmwareVersion
	if err := c.ShouldBindJSON(&firmware); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.CreateFirmwareVersion(&firmware); err != nil {
		global.GVA_LOG.Error("创建固件版本失败!", zap.Error(err))
		response.FailWithMessage("创建固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteFirmwareVersion 删除固件版本
func (a *FirmwareVersionApi) DeleteFirmwareVersion(c *gin.Context) {
	id := c.Query("ID")
	if err := firmwareVersionService.DeleteFirmwareVersion(id); err != nil {
		global.GVA_LOG.Error("删除固件版本失败!", zap.Error(err))
		response.FailWithMessage("删除固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteFirmwareVersionByIds 批量删除固件版本
func (a *FirmwareVersionApi) DeleteFirmwareVersionByIds(c *gin.Context) {
	var ids commonReq.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.DeleteFirmwareVersionByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除固件版本失败!", zap.Error(err))
		response.FailWithMessage("批量删除固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateFirmwareVersion 更新固件版本
func (a *FirmwareVersionApi) UpdateFirmwareVersion(c *gin.Context) {
	var firmware deviceModel.FirmwareVersion
	if err := c.ShouldBindJSON(&firmware); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.UpdateFirmwareVersion(firmware); err != nil {
		global.GVA_LOG.Error("更新固件版本失败!", zap.Error(err))
		response.FailWithMessage("更新固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindFirmwareVersion 获取固件版本详情
func (a *FirmwareVersionApi) FindFirmwareVersion(c *gin.Context) {
	id := c.Query("ID")
	firmware, err := firmwareVersionService.GetFirmwareVersion(id)
	if err != nil {
		global.GVA_LOG.Error("查询固件版本失败!", zap.Error(err))
		response.FailWithMessage("查询固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithData(firmware, c)
}

// GetFirmwareVersionList 获取固件版本列表
func (a *FirmwareVersionApi) GetFirmwareVersionList(c *gin.Context) {
	var pageInfo deviceReq.FirmwareVersionSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := firmwareVersionService.GetFirmwareVersionInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取固件版本列表失败!", zap.Error(err))
		response.FailWithMessage("获取固件版本列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ChangeFirmwareVersionStatus 更新固件版本状态
func (a *FirmwareVersionApi) ChangeFirmwareVersionStatus(c *gin.Context) {
	var req deviceReq.ChangeFirmwareVersionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.ChangeFirmwareVersionStatus(req); err != nil {
		global.GVA_LOG.Error("更新固件版本状态失败!", zap.Error(err))
		response.FailWithMessage("更新固件版本状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("状态更新成功", c)
}

// PublishFirmwareVersion 发布固件版本
func (a *FirmwareVersionApi) PublishFirmwareVersion(c *gin.Context) {
	var req deviceReq.PublishFirmwareVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.PublishFirmwareVersion(req); err != nil {
		global.GVA_LOG.Error("发布固件版本失败!", zap.Error(err))
		response.FailWithMessage("发布固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("发布成功", c)
}

// SetFirmwareStable 标记稳定版本
func (a *FirmwareVersionApi) SetFirmwareStable(c *gin.Context) {
	var req deviceReq.SetFirmwareStableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.SetFirmwareStable(req); err != nil {
		global.GVA_LOG.Error("设置稳定版本失败!", zap.Error(err))
		response.FailWithMessage("设置稳定版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// VoidFirmwareVersion 作废固件版本
func (a *FirmwareVersionApi) VoidFirmwareVersion(c *gin.Context) {
	var req deviceReq.VoidFirmwareVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.VoidFirmwareVersion(req); err != nil {
		global.GVA_LOG.Error("作废固件版本失败!", zap.Error(err))
		response.FailWithMessage("作废固件版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("作废成功", c)
}

// DeleteFirmwarePackage 删除固件包
func (a *FirmwareVersionApi) DeleteFirmwarePackage(c *gin.Context) {
	var req deviceReq.DeleteFirmwarePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionService.DeleteFirmwarePackage(req); err != nil {
		global.GVA_LOG.Error("删除固件包失败!", zap.Error(err))
		response.FailWithMessage("删除固件包失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
