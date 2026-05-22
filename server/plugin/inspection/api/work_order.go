package api

import (
	"net/url"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var WorkOrder = new(workOrderApi)

type workOrderApi struct{}

// StartInspection 接收并开始检测
// @Tags     WorkOrder
// @Summary  接收并开始检测（状态:待检测接收→检测中）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.StartInspection true "订单ID"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/startInspection [post]
func (a *workOrderApi) StartInspection(c *gin.Context) {
	var req request.StartInspection
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.StartInspection(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已接收并开始检测", c)
}

// StartRecheck 开始复检
// @Tags     WorkOrder
// @Summary  开始复检（待复检→复检中）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.StartRecheck true "批次ID"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/startRecheck [post]
func (a *workOrderApi) StartRecheck(c *gin.Context) {
	var req request.StartRecheck
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.StartRecheck(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已开始复检", c)
}

// AssignBatchTemplate 为批次选择模板并提交检测接收
// @Tags     WorkOrder
// @Summary  为批次选择模板并提交检测接收
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AssignBatchTemplate true "批次模板信息"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/assignBatchTemplate [post]
func (a *workOrderApi) AssignBatchTemplate(c *gin.Context) {
	var req request.AssignBatchTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.AssignBatchTemplate(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已提交检测接收", c)
}

// AssignOrderTemplate 为生产订单选择模板并提交未派检批次
// @Tags     WorkOrder
// @Summary  为生产订单选择模板并提交未派检批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AssignOrderTemplate true "生产订单模板信息"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/assignOrderTemplate [post]
func (a *workOrderApi) AssignOrderTemplate(c *gin.Context) {
	var req request.AssignOrderTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.AssignOrderTemplate(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已提交检测接收", c)
}

// SaveResults 保存检测结果
// @Tags     WorkOrder
// @Summary  保存检测结果
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.SaveInspectionResult true "检测结果"
// @Success  200 {object} response.Response{msg=string} "保存成功"
// @Router   /workOrder/saveResults [post]
func (a *workOrderApi) SaveResults(c *gin.Context) {
	var req request.SaveInspectionResult
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.SaveResults(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// SaveSingleResult 保存单项检测结果
// @Tags     WorkOrder
// @Summary  保存单项检测结果并记录检测人
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.SaveSingleInspectionResult true "单项检测结果"
// @Success  200 {object} response.Response{data=object,msg=string} "保存成功"
// @Router   /workOrder/saveSingleResult [post]
func (a *workOrderApi) SaveSingleResult(c *gin.Context) {
	var req request.SaveSingleInspectionResult
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	result, err := serviceWorkOrder.SaveSingleResult(&req, claims.BaseClaims.ID, claims.NickName)
	if err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "保存成功", c)
}

// CompleteInspection 提交检测待确认
// @Tags     WorkOrder
// @Summary  提交检测待确认（状态:检测中→待确认）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CompleteInspection true "订单ID"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/completeInspection [post]
func (a *workOrderApi) CompleteInspection(c *gin.Context) {
	var req request.CompleteInspection
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.CompleteInspection(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已提交待确认", c)
}

// ConfirmInspectionComplete 确认完成检测
// @Tags     WorkOrder
// @Summary  确认完成检测（状态:待确认→已完成）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ConfirmInspectionComplete true "批次ID"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/confirmInspectionComplete [post]
func (a *workOrderApi) ConfirmInspectionComplete(c *gin.Context) {
	var req request.ConfirmInspectionComplete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.ConfirmInspectionComplete(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已确认完成", c)
}

// CompleteRecheck 完成复检
// @Tags     WorkOrder
// @Summary  完成复检
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CompleteRecheck true "批次ID"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/completeRecheck [post]
func (a *workOrderApi) CompleteRecheck(c *gin.Context) {
	var req request.CompleteRecheck
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.CompleteRecheck(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("复检完成", c)
}

// ReturnDevices 设备打回生产
// @Tags     WorkOrder
// @Summary  设备打回生产
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ReturnDevices true "打回信息"
// @Success  200 {object} response.Response{msg=string} "操作成功"
// @Router   /workOrder/returnDevices [post]
func (a *workOrderApi) ReturnDevices(c *gin.Context) {
	var req request.ReturnDevices
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceWorkOrder.ReturnDevices(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已打回生产", c)
}

// GetInspectionDetail 获取检测详情
// @Tags     WorkOrder
// @Summary  获取检测详情（含设备+检测项+已有结果）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "生产订单ID"
// @Success  200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router   /workOrder/getInspectionDetail [get]
func (a *workOrderApi) GetInspectionDetail(c *gin.Context) {
	id := c.Query("id")
	data, err := serviceWorkOrder.GetInspectionDetailData(id)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// GetBatchStatusLogs 查询批次流转日志
// @Tags     WorkOrder
// @Summary  查询批次流转日志
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    batchID query string true "批次ID"
// @Success  200 {object} response.Response{data=[]model.ProductionBatchStatusLog,msg=string} "查询成功"
// @Router   /workOrder/getBatchStatusLogs [get]
func (a *workOrderApi) GetBatchStatusLogs(c *gin.Context) {
	batchID := c.Query("batchID")
	logs, err := serviceWorkOrder.GetBatchStatusLogs(batchID)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(logs, c)
}

// GetFlowLogs 查询融合流转日志
// @Tags     WorkOrder
// @Summary  查询融合流转日志
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    batchID query string false "批次ID"
// @Param    deviceID query string false "设备ID"
// @Success  200 {object} response.Response{data=[]service.FlowLogItem,msg=string} "查询成功"
// @Router   /workOrder/getFlowLogs [get]
func (a *workOrderApi) GetFlowLogs(c *gin.Context) {
	batchID := c.Query("batchID")
	deviceID := c.Query("deviceID")
	logs, err := serviceWorkOrder.GetFlowLogs(batchID, deviceID)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(logs, c)
}

// ExportInspectionExcel 导出检测工单Excel
// @Tags     WorkOrder
// @Summary  导出检测工单Excel
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/octet-stream
// @Param    id query string true "批次ID"
// @Success  200 {file} file "检测工单Excel"
// @Router   /workOrder/exportInspectionExcel [get]
func (a *workOrderApi) ExportInspectionExcel(c *gin.Context) {
	id := c.Query("id")
	buf, filename, err := serviceWorkOrder.ExportInspectionExcel(id)
	if err != nil {
		response.FailWithMessage("导出失败: "+err.Error(), c)
		return
	}
	escaped := url.QueryEscape(filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+escaped)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// GetInspectionBatchList 获取检测批次列表
// @Tags     WorkOrder
// @Summary  获取检测批次列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Param    moNumber query string false "MO号"
// @Param    model query string false "型号"
// @Param    status query int false "状态"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router   /workOrder/getInspectionBatchList [get]
func (a *workOrderApi) GetInspectionBatchList(c *gin.Context) {
	var search request.InspectionBatchSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceWorkOrder.GetInspectionBatchList(search)
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
