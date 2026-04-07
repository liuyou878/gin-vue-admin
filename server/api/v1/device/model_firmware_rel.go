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

type ModelFirmwareRelApi struct{}

// CreateModelFirmwareRel 创建型号固件关系
func (a *ModelFirmwareRelApi) CreateModelFirmwareRel(c *gin.Context) {
	var rel deviceModel.ModelFirmwareRel
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := modelFirmwareRelService.CreateModelFirmwareRel(&rel); err != nil {
		global.GVA_LOG.Error("创建型号固件关系失败!", zap.Error(err))
		response.FailWithMessage("创建型号固件关系失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteModelFirmwareRel 删除型号固件关系
func (a *ModelFirmwareRelApi) DeleteModelFirmwareRel(c *gin.Context) {
	id := c.Query("ID")
	if err := modelFirmwareRelService.DeleteModelFirmwareRel(id); err != nil {
		global.GVA_LOG.Error("删除型号固件关系失败!", zap.Error(err))
		response.FailWithMessage("删除型号固件关系失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteModelFirmwareRelByIds 批量删除型号固件关系
func (a *ModelFirmwareRelApi) DeleteModelFirmwareRelByIds(c *gin.Context) {
	var ids commonReq.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := modelFirmwareRelService.DeleteModelFirmwareRelByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除型号固件关系失败!", zap.Error(err))
		response.FailWithMessage("批量删除型号固件关系失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateModelFirmwareRel 更新型号固件关系
func (a *ModelFirmwareRelApi) UpdateModelFirmwareRel(c *gin.Context) {
	var rel deviceModel.ModelFirmwareRel
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := modelFirmwareRelService.UpdateModelFirmwareRel(rel); err != nil {
		global.GVA_LOG.Error("更新型号固件关系失败!", zap.Error(err))
		response.FailWithMessage("更新型号固件关系失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindModelFirmwareRel 获取型号固件关系详情
func (a *ModelFirmwareRelApi) FindModelFirmwareRel(c *gin.Context) {
	id := c.Query("ID")
	rel, err := modelFirmwareRelService.GetModelFirmwareRel(id)
	if err != nil {
		global.GVA_LOG.Error("查询型号固件关系失败!", zap.Error(err))
		response.FailWithMessage("查询型号固件关系失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rel, c)
}

// GetModelFirmwareRelList 获取型号固件关系列表
func (a *ModelFirmwareRelApi) GetModelFirmwareRelList(c *gin.Context) {
	var pageInfo deviceReq.ModelFirmwareRelSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := modelFirmwareRelService.GetModelFirmwareRelInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取型号固件关系列表失败!", zap.Error(err))
		response.FailWithMessage("获取型号固件关系列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetModelFirmwareRecommended 设置推荐版本
func (a *ModelFirmwareRelApi) SetModelFirmwareRecommended(c *gin.Context) {
	var req deviceReq.SetModelFirmwareRecommendedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := modelFirmwareRelService.SetModelFirmwareRecommended(req); err != nil {
		global.GVA_LOG.Error("设置推荐版本失败!", zap.Error(err))
		response.FailWithMessage("设置推荐版本失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// SetModelFirmwareTestResult 设置测试结果
func (a *ModelFirmwareRelApi) SetModelFirmwareTestResult(c *gin.Context) {
	var req deviceReq.SetModelFirmwareTestResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := modelFirmwareRelService.SetModelFirmwareTestResult(req); err != nil {
		global.GVA_LOG.Error("设置测试结果失败!", zap.Error(err))
		response.FailWithMessage("设置测试结果失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}
