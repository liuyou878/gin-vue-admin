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

type FirmwareTagApi struct{}

// CreateFirmwareTag 创建固件标签
func (a *FirmwareTagApi) CreateFirmwareTag(c *gin.Context) {
	var tag deviceModel.FirmwareTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareTagService.CreateFirmwareTag(&tag); err != nil {
		global.GVA_LOG.Error("创建固件标签失败!", zap.Error(err))
		response.FailWithMessage("创建固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteFirmwareTag 删除固件标签
func (a *FirmwareTagApi) DeleteFirmwareTag(c *gin.Context) {
	id := c.Query("ID")
	if err := firmwareTagService.DeleteFirmwareTag(id); err != nil {
		global.GVA_LOG.Error("删除固件标签失败!", zap.Error(err))
		response.FailWithMessage("删除固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteFirmwareTagByIds 批量删除固件标签
func (a *FirmwareTagApi) DeleteFirmwareTagByIds(c *gin.Context) {
	var ids commonReq.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareTagService.DeleteFirmwareTagByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除固件标签失败!", zap.Error(err))
		response.FailWithMessage("批量删除固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateFirmwareTag 更新固件标签
func (a *FirmwareTagApi) UpdateFirmwareTag(c *gin.Context) {
	var tag deviceModel.FirmwareTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareTagService.UpdateFirmwareTag(tag); err != nil {
		global.GVA_LOG.Error("更新固件标签失败!", zap.Error(err))
		response.FailWithMessage("更新固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindFirmwareTag 获取固件标签详情
func (a *FirmwareTagApi) FindFirmwareTag(c *gin.Context) {
	id := c.Query("ID")
	tag, err := firmwareTagService.GetFirmwareTag(id)
	if err != nil {
		global.GVA_LOG.Error("查询固件标签失败!", zap.Error(err))
		response.FailWithMessage("查询固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithData(tag, c)
}

// GetFirmwareTagList 获取固件标签列表
func (a *FirmwareTagApi) GetFirmwareTagList(c *gin.Context) {
	var pageInfo deviceReq.FirmwareTagSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := firmwareTagService.GetFirmwareTagInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取固件标签列表失败!", zap.Error(err))
		response.FailWithMessage("获取固件标签列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetFirmwareTags 设置固件版本标签
func (a *FirmwareTagApi) SetFirmwareTags(c *gin.Context) {
	var req deviceReq.SetFirmwareTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := firmwareVersionTagRelService.SetFirmwareTags(req); err != nil {
		global.GVA_LOG.Error("设置固件标签失败!", zap.Error(err))
		response.FailWithMessage("设置固件标签失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}
