package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var ProductionOrder = new(productionOrderApi)

type productionOrderApi struct{}

// SubmitDeviceData 生产工具提交全量数据
// @Tags     ProductionOrder
// @Summary  生产工具提交全量数据(SN+注册码+GETALL)
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.SubmitDeviceData true "提交数据"
// @Success  200 {object} response.Response{msg=string} "提交成功"
// @Router   /productionOrder/submitDeviceData [post]
func (a *productionOrderApi) SubmitDeviceData(c *gin.Context) {
	var req request.SubmitDeviceData
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.SubmitDeviceData(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("提交失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("提交成功", c)
}

// CreateProductionOrder 创建生产订单
// @Tags     ProductionOrder
// @Summary  创建生产订单
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

// ForceDeleteProductionOrder 强制删除生产订单
// @Tags     ProductionOrder
// @Summary  强制删除生产订单（用于清理已确认/测试数据）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "订单ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /productionOrder/forceDeleteProductionOrder [delete]
func (a *productionOrderApi) ForceDeleteProductionOrder(c *gin.Context) {
	id := c.Query("id")
	if err := serviceProductionOrder.ForceDeleteProductionOrder(id); err != nil {
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
// @Summary  查询生产订单详情(含批次+设备)
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

// GetSubmittedDeviceList 获取生产工具提交设备列表
// @Tags     ProductionOrder
// @Summary  获取生产工具提交设备列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Param    moNumber query string false "MO号"
// @Param    batchNumber query string false "批次号"
// @Param    sn query string false "SN"
// @Param    model query string false "型号"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router   /productionOrder/getSubmittedDeviceList [get]
func (a *productionOrderApi) GetSubmittedDeviceList(c *gin.Context) {
	var search request.SubmittedDeviceSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceProductionOrder.GetSubmittedDeviceList(search)
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

// FindSubmittedDevice 查询生产工具提交设备详情
// @Tags     ProductionOrder
// @Summary  查询生产工具提交设备详情（含GETALL原文）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "设备ID"
// @Success  200 {object} response.Response{data=model.ProductionOrderDevice,msg=string} "查询成功"
// @Router   /productionOrder/findSubmittedDevice [get]
func (a *productionOrderApi) FindSubmittedDevice(c *gin.Context) {
	id := c.Query("id")
	device, err := serviceProductionOrder.FindSubmittedDevice(id)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(device, c)
}

// DeleteSubmittedDevice 删除生产工具提交设备记录
// @Tags     ProductionOrder
// @Summary  删除生产工具提交设备记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "设备ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /productionOrder/deleteSubmittedDevice [delete]
func (a *productionOrderApi) DeleteSubmittedDevice(c *gin.Context) {
	id := c.Query("id")
	if err := serviceProductionOrder.DeleteSubmittedDevice(id); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ConfirmReworkDone 生产确认返工完成
// @Tags     ProductionOrder
// @Summary  生产确认返工完成（返工中→待复检）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ConfirmReworkDone true "返工确认信息"
// @Success  200 {object} response.Response{msg=string} "确认成功"
// @Router   /productionOrder/confirmReworkDone [post]
func (a *productionOrderApi) ConfirmReworkDone(c *gin.Context) {
	var req request.ConfirmReworkDone
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.ConfirmReworkDone(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("确认失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已进入待复检", c)
}

// ConfirmReworkReceived 生产确认接收返工
// @Tags     ProductionOrder
// @Summary  生产确认接收返工（待生产接收→返工中）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ConfirmReworkReceived true "返工接收信息"
// @Success  200 {object} response.Response{msg=string} "确认成功"
// @Router   /productionOrder/confirmReworkReceived [post]
func (a *productionOrderApi) ConfirmReworkReceived(c *gin.Context) {
	var req request.ConfirmReworkReceived
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.ConfirmReworkReceived(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("确认失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已进入返工中", c)
}

// GetDeviceStatusLogs 查询设备状态日志
// @Tags     ProductionOrder
// @Summary  查询设备状态日志
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    deviceID query string true "设备ID"
// @Success  200 {object} response.Response{data=[]model.ProductionOrderDeviceStatusLog,msg=string} "查询成功"
// @Router   /productionOrder/getDeviceStatusLogs [get]
func (a *productionOrderApi) GetDeviceStatusLogs(c *gin.Context) {
	deviceID := c.Query("deviceID")
	logs, err := serviceProductionOrder.GetDeviceStatusLogs(deviceID)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(logs, c)
}

// AssignBatch 分配设备到批次
// @Tags     ProductionOrder
// @Summary  分配SN到批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AssignBatch true "分配信息"
// @Success  200 {object} response.Response{msg=string} "分配成功"
// @Router   /productionOrder/assignBatch [post]
func (a *productionOrderApi) AssignBatch(c *gin.Context) {
	var req request.AssignBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.AssignBatch(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("分配失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("分配成功", c)
}

// ScanAssignBatch 扫码加入生产批次
// @Tags     ProductionOrder
// @Summary  扫码加入生产批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ScanAssignBatch true "扫码分批信息"
// @Success  200 {object} response.Response{msg=string} "分批成功"
// @Router   /productionOrder/scanAssignBatch [post]
func (a *productionOrderApi) ScanAssignBatch(c *gin.Context) {
	var req request.ScanAssignBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.ScanAssignBatch(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("分批失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("分批成功", c)
}

// CreateBatch 创建批次
// @Tags     ProductionOrder
// @Summary  创建生产批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateBatch true "批次信息"
// @Success  200 {object} response.Response{msg=string} "创建成功"
// @Router   /productionOrder/createBatch [post]
func (a *productionOrderApi) CreateBatch(c *gin.Context) {
	var req request.CreateBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	batch := &model.ProductionBatch{
		ProductionOrderID: req.ProductionOrderID,
		BatchNumber:       req.BatchNumber,
		Status:            0,
	}
	if err := global.GVA_DB.Create(batch).Error; err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// AddDevicesToBatch 添加设备到批次
// @Tags     ProductionOrder
// @Summary  添加设备到批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AddDevicesToBatch true "批次ID和SN列表"
// @Success  200 {object} response.Response{msg=string} "添加成功"
// @Router   /productionOrder/addDevicesToBatch [post]
func (a *productionOrderApi) AddDevicesToBatch(c *gin.Context) {
	var req request.AddDevicesToBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.AddDevicesToBatch(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("添加失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// RemoveDeviceFromBatch 从批次移除设备
// @Tags     ProductionOrder
// @Summary  从批次移除设备
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.RemoveDeviceFromBatch true "设备ID"
// @Success  200 {object} response.Response{msg=string} "移除成功"
// @Router   /productionOrder/removeDeviceFromBatch [post]
func (a *productionOrderApi) RemoveDeviceFromBatch(c *gin.Context) {
	var req request.RemoveDeviceFromBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.RemoveDeviceFromBatch(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("移除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("移除成功", c)
}

// DeleteEmptyBatch 删除空批次
// @Tags     ProductionOrder
// @Summary  删除未派检且没有设备的生产批次
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.DeleteEmptyBatch true "批次ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /productionOrder/deleteEmptyBatch [post]
func (a *productionOrderApi) DeleteEmptyBatch(c *gin.Context) {
	var req request.DeleteEmptyBatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceProductionOrder.DeleteEmptyBatch(&req); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateBatchNumber 修改批次号
// @Tags     ProductionOrder
// @Summary  修改未完成批次的批次号
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateBatchNumber true "批次号"
// @Success  200 {object} response.Response{msg=string} "修改成功"
// @Router   /productionOrder/updateBatchNumber [post]
func (a *productionOrderApi) UpdateBatchNumber(c *gin.Context) {
	var req request.UpdateBatchNumber
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims, _ := utils.GetClaims(c)
	if err := serviceProductionOrder.UpdateBatchNumber(&req, claims.BaseClaims.ID, claims.NickName); err != nil {
		response.FailWithMessage("修改失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}
