package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
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
	ensureMenuButton(ctx, "inspectionItem", "delete", "删除")
	ensureMenuButton(ctx, "inspectionTemplate", "delete", "删除")
	ensureMenuButton(ctx, "productionOrder", "delete", "删除")
	ensureMenuButton(ctx, "submittedDevice", "delete", "删除")
}

func ensureMenuButton(ctx context.Context, menuName string, btnName string, desc string) {
	if global.GVA_DB == nil {
		return
	}
	var menu model.SysBaseMenu
	if err := global.GVA_DB.WithContext(ctx).Where("name = ?", menuName).First(&menu).Error; err != nil {
		return
	}
	btn := model.SysBaseMenuBtn{
		Name:          btnName,
		Desc:          desc,
		SysBaseMenuID: menu.ID,
	}
	_ = global.GVA_DB.WithContext(ctx).
		Where("sys_base_menu_id = ? AND name = ?", menu.ID, btnName).
		FirstOrCreate(&btn).Error
}
