package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
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
		{Path: "/inspectionTemplate/deleteTemplate", Description: "删除检测模板", ApiGroup: "检测模板", Method: "DELETE"},
		{Path: "/inspectionTemplate/updateTemplate", Description: "更新检测模板", ApiGroup: "检测模板", Method: "PUT"},
		{Path: "/inspectionTemplate/findTemplate", Description: "查询模板详情", ApiGroup: "检测模板", Method: "GET"},
		{Path: "/inspectionTemplate/getTemplateList", Description: "获取模板列表", ApiGroup: "检测模板", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
