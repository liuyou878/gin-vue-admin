package initialize

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	pluginUtils "github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	serverUtils "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/inspectionItem/createItem", Description: "新增检测项", ApiGroup: "检测项管理", Method: "POST"},
		{Path: "/inspectionItem/deleteItem", Description: "删除检测项", ApiGroup: "检测项管理", Method: "DELETE"},
		{Path: "/inspectionItem/deleteItemByIds", Description: "批量删除检测项", ApiGroup: "检测项管理", Method: "DELETE"},
		{Path: "/inspectionItem/updateItem", Description: "更新检测项", ApiGroup: "检测项管理", Method: "PUT"},
		{Path: "/inspectionItem/findItem", Description: "根据ID获取检测项", ApiGroup: "检测项管理", Method: "GET"},
		{Path: "/inspectionItem/getItemList", Description: "获取检测项列表", ApiGroup: "检测项管理", Method: "GET"},
		{Path: "/inspectionTemplate/createTemplate", Description: "新增检测模板", ApiGroup: "检测模板", Method: "POST"},
		{Path: "/inspectionTemplate/copyTemplate", Description: "复制检测模板", ApiGroup: "检测模板", Method: "POST"},
		{Path: "/inspectionTemplate/deleteTemplate", Description: "删除检测模板", ApiGroup: "检测模板", Method: "DELETE"},
		{Path: "/inspectionTemplate/updateTemplate", Description: "更新检测模板", ApiGroup: "检测模板", Method: "PUT"},
		{Path: "/inspectionTemplate/findTemplate", Description: "查询模板详情", ApiGroup: "检测模板", Method: "GET"},
		{Path: "/inspectionTemplate/getTemplateList", Description: "获取模板列表", ApiGroup: "检测模板", Method: "GET"},
		{Path: "/productionOrder/submitDeviceData", Description: "生产工具提交全量数据", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/createProductionOrder", Description: "创建生产订单", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/confirmReworkDone", Description: "生产确认返工完成", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/deleteProductionOrder", Description: "删除生产订单", ApiGroup: "生产订单", Method: "DELETE"},
		{Path: "/productionOrder/forceDeleteProductionOrder", Description: "强制删除生产订单", ApiGroup: "生产订单", Method: "DELETE"},
		{Path: "/productionOrder/updateProductionOrder", Description: "更新生产订单", ApiGroup: "生产订单", Method: "PUT"},
		{Path: "/productionOrder/assignBatch", Description: "分配SN到批次", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/scanAssignBatch", Description: "扫码加入生产批次", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/createBatch", Description: "创建生产批次", ApiGroup: "生产订单", Method: "POST"},
		{Path: "/productionOrder/findProductionOrder", Description: "查询生产订单详情", ApiGroup: "生产订单", Method: "GET"},
		{Path: "/productionOrder/getProductionOrderList", Description: "获取生产订单列表", ApiGroup: "生产订单", Method: "GET"},
		{Path: "/productionOrder/getSubmittedDeviceList", Description: "获取生产工具提交设备列表", ApiGroup: "生产订单", Method: "GET"},
		{Path: "/productionOrder/findSubmittedDevice", Description: "查询生产工具提交设备详情", ApiGroup: "生产订单", Method: "GET"},
		{Path: "/productionOrder/getDeviceStatusLogs", Description: "查询设备状态日志", ApiGroup: "生产订单", Method: "GET"},
		{Path: "/productionOrder/deleteSubmittedDevice", Description: "删除生产工具提交设备记录", ApiGroup: "生产订单", Method: "DELETE"},
		{Path: "/workOrder/assignBatchTemplate", Description: "为批次选择模板并生成待检测工单", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/assignOrderTemplate", Description: "为生产订单选择模板并提交未派检批次", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/startInspection", Description: "开始检测", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/startRecheck", Description: "开始复检", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/saveResults", Description: "保存检测结果", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/saveSingleResult", Description: "保存单项检测结果", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/completeInspection", Description: "完成检测", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/completeRecheck", Description: "完成复检", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/returnDevices", Description: "设备打回生产", ApiGroup: "检测工单", Method: "POST"},
		{Path: "/workOrder/getInspectionBatchList", Description: "获取检测批次列表", ApiGroup: "检测工单", Method: "GET"},
		{Path: "/workOrder/getInspectionDetail", Description: "获取检测详情", ApiGroup: "检测工单", Method: "GET"},
		{Path: "/workOrder/exportInspectionExcel", Description: "导出检测工单Excel", ApiGroup: "检测工单", Method: "GET"},
	}
	pluginUtils.RegisterApis(entities...)
	grantSaveSingleResultPermission()
}

func grantSaveSingleResultPermission() {
	if global.GVA_DB == nil {
		return
	}
	var rules []adapter.CasbinRule
	if err := global.GVA_DB.
		Where("ptype = ? AND v1 = ? AND v2 = ?", "p", "/workOrder/saveResults", "POST").
		Find(&rules).Error; err != nil {
		return
	}
	for _, rule := range rules {
		newRule := adapter.CasbinRule{
			Ptype: "p",
			V0:    rule.V0,
			V1:    "/workOrder/saveSingleResult",
			V2:    "POST",
		}
		_ = global.GVA_DB.
			Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", newRule.Ptype, newRule.V0, newRule.V1, newRule.V2).
			FirstOrCreate(&newRule).Error
	}
	_ = serverUtils.GetCasbin().LoadPolicy()
}
