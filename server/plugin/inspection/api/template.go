package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/gin-gonic/gin"
)

var Template = new(templateApi)

type templateApi struct{}

// CreateTemplate 新增检测模板
// @Tags     InspectionTemplate
// @Summary  新增检测模板
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateTemplate true "模板信息"
// @Success  200 {object} response.Response{msg=string} "创建成功"
// @Router   /inspectionTemplate/createTemplate [post]
func (a *templateApi) CreateTemplate(c *gin.Context) {
	var req request.CreateTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceTemplate.CreateTemplate(&req); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteTemplate 删除检测模板
// @Tags     InspectionTemplate
// @Summary  删除检测模板
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "模板ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /inspectionTemplate/deleteTemplate [delete]
func (a *templateApi) DeleteTemplate(c *gin.Context) {
	id := c.Query("id")
	if err := serviceTemplate.DeleteTemplate(id); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateTemplate 更新检测模板
// @Tags     InspectionTemplate
// @Summary  更新检测模板
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateTemplate true "模板信息"
// @Success  200 {object} response.Response{msg=string} "更新成功"
// @Router   /inspectionTemplate/updateTemplate [put]
func (a *templateApi) UpdateTemplate(c *gin.Context) {
	var req request.UpdateTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceTemplate.UpdateTemplate(&req); err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// CopyTemplate 复制检测模板
// @Tags     InspectionTemplate
// @Summary  复制检测模板
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CopyTemplate true "复制模板信息"
// @Success  200 {object} response.Response{msg=string} "复制成功"
// @Router   /inspectionTemplate/copyTemplate [post]
func (a *templateApi) CopyTemplate(c *gin.Context) {
	var req request.CopyTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceTemplate.CopyTemplate(&req); err != nil {
		response.FailWithMessage("复制失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("复制成功", c)
}

// FindTemplate 查询模板详情
// @Tags     InspectionTemplate
// @Summary  查询模板详情(含检测项列表)
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "模板ID"
// @Success  200 {object} response.Response{data=model.InspectionTemplate,msg=string} "查询成功"
// @Router   /inspectionTemplate/findTemplate [get]
func (a *templateApi) FindTemplate(c *gin.Context) {
	id := c.Query("id")
	tmpl, err := serviceTemplate.FindTemplate(id)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(tmpl, c)
}

// GetTemplateList 获取模板列表
// @Tags     InspectionTemplate
// @Summary  获取模板列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Param    name query string false "模板名称"
// @Param    model query string false "型号"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router   /inspectionTemplate/getTemplateList [get]
func (a *templateApi) GetTemplateList(c *gin.Context) {
	var search request.TemplateSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceTemplate.GetTemplateList(search)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "查询成功", c)
}
