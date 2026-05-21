package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	parent := model.SysBaseMenu{
		ParentId: 0, Path: "inspection", Name: "inspection", Hidden: false,
		Component: "view/routerHolder.vue", Sort: 7,
		Meta: model.Meta{Title: "检测管理", Icon: "monitor"},
	}
	child1 := model.SysBaseMenu{
		ParentId: 0, Path: "inspectionWorkOrder", Name: "inspectionWorkOrder", Hidden: false,
		Component: "plugin/inspection/view/work_order_redirect.vue", Sort: 1,
		Meta: model.Meta{Title: "检测工单", Icon: "document"},
	}
	child2 := model.SysBaseMenu{
		ParentId: 0, Path: "inspectionItem", Name: "inspectionItem", Hidden: false,
		Component: "plugin/inspection/view/inspection_item.vue", Sort: 2,
		Meta: model.Meta{Title: "检测项管理", Icon: "list"},
	}
	child3 := model.SysBaseMenu{
		ParentId: 0, Path: "inspectionTemplate", Name: "inspectionTemplate", Hidden: false,
		Component: "plugin/inspection/view/template.vue", Sort: 3,
		Meta: model.Meta{Title: "检测模板", Icon: "documentation"},
	}
	child4 := model.SysBaseMenu{
		ParentId: 0, Path: "productionOrder", Name: "productionOrder", Hidden: false,
		Component: "plugin/inspection/view/production_order.vue", Sort: 4,
		Meta: model.Meta{Title: "生产订单", Icon: "tickets"},
	}
	child5 := model.SysBaseMenu{
		ParentId: 0, Path: "submittedDevice", Name: "submittedDevice", Hidden: false,
		Component: "plugin/inspection/view/submitted_device.vue", Sort: 5,
		Meta: model.Meta{Title: "生产提交数据", Icon: "data-board"},
	}
	utils.RegisterMenus(parent, child1, child2, child3, child4, child5)
}
