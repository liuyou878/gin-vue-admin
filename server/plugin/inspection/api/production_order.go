package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/gin-gonic/gin"
)

var ProductionOrder = new(productionOrderApi)

type productionOrderApi struct{}

// CreateProductionOrder 创建生产订单
// @Tags     ProductionOrder
// @Summary  创建生产订单（含 SN 列表）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateProductionOrder true "生产订单信息"
// @Success  200 {object} response.Response{msg=string} "创建成功"
// @Router   /productionOrder/createProductionOrder [post]
func (a *productionOrderApi) CreateProductionOrder(c *gin.Context) {
	var req request.CreateProductionOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceProductionOrder.CreateProductionOrder(&req); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteProductionOrder 删除生产订单
// @Tags     ProductionOrder
// @Summary  删除生产订单
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "订单ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /productionOrder/deleteProductionOrder [delete]
func (a *productionOrderApi) DeleteProductionOrder(c *gin.Context) {
	id := c.Query("id")
	if err := serviceProductionOrder.DeleteProductionOrder(id); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateProductionOrder 更新生产订单
// @Tags     ProductionOrder
// @Summary  更新生产订单
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateProductionOrder true "生产订单信息"
// @Success  200 {object} response.Response{msg=string} "更新成功"
// @Router   /productionOrder/updateProductionOrder [put]
func (a *productionOrderApi) UpdateProductionOrder(c *gin.Context) {
	var req request.UpdateProductionOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceProductionOrder.UpdateProductionOrder(&req); err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindProductionOrder 查询生产订单详情
// @Tags     ProductionOrder
// @Summary  查询生产订单详情（含 SN 列表）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "订单ID"
// @Success  200 {object} response.Response{data=model.ProductionOrder,msg=string} "查询成功"
// @Router   /productionOrder/findProductionOrder [get]
func (a *productionOrderApi) FindProductionOrder(c *gin.Context) {
	id := c.Query("id")
	po, err := serviceProductionOrder.FindProductionOrder(id)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(po, c)
}

// GetProductionOrderList 获取生产订单列表
// @Tags     ProductionOrder
// @Summary  获取生产订单列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Param    moNumber query string false "MO号"
// @Param    model query string false "型号"
// @Param    instrumentCategory query string false "仪器类别"
// @Param    status query int false "状态"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router   /productionOrder/getProductionOrderList [get]
func (a *productionOrderApi) GetProductionOrderList(c *gin.Context) {
	var search request.ProductionOrderSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceProductionOrder.GetProductionOrderList(search)
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
