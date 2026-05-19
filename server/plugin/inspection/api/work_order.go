package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var WorkOrder = new(workOrderApi)

type workOrderApi struct{}

// StartInspection 开始检测
// @Tags     WorkOrder
// @Summary  开始检测（状态:待检测→检测中）
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
	response.OkWithMessage("已开始检测", c)
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
	if err := serviceWorkOrder.SaveResults(&req); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// CompleteInspection 完成检测
// @Tags     WorkOrder
// @Summary  完成检测（状态:检测中→已完成）
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
	if err := serviceWorkOrder.CompleteInspection(&req); err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("检测完成", c)
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
