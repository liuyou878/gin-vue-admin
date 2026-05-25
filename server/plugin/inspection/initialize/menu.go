package initialize

import (
	"context"
	"fmt"

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
		Component: "plugin/inspection/view/work_order.vue", Sort: 1,
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
	child6 := model.SysBaseMenu{
		ParentId: 0, Path: "productionMockSubmit", Name: "productionMockSubmit", Hidden: false,
		Component: "plugin/inspection/view/production_mock_submit.vue", Sort: 6,
		Meta: model.Meta{Title: "模拟生产提交", Icon: "cpu"},
	}
	child7 := model.SysBaseMenu{
		ParentId: 0, Path: "batchImportDevice", Name: "batchImportDevice", Hidden: false,
		Component: "plugin/inspection/view/batch_import.vue", Sort: 7,
		Meta: model.Meta{Title: "批量导入设备", Icon: "upload"},
	}
	utils.RegisterMenus(parent, child1, child2, child3, child4, child5, child6, child7)
	syncWorkOrderMenu(ctx)
	ensureMenuButton(ctx, "inspectionItem", "delete", "删除")
	ensureMenuButton(ctx, "inspectionTemplate", "delete", "删除")
	ensureMenuButton(ctx, "productionOrder", "delete", "删除")
	ensureMenuButton(ctx, "submittedDevice", "delete", "删除")
}

func syncWorkOrderMenu(ctx context.Context) {
	if global.GVA_DB == nil {
		return
	}
	var keep model.SysBaseMenu
	if err := global.GVA_DB.WithContext(ctx).
		Where("name = ?", "inspectionWorkOrder").
		First(&keep).Error; err == nil {
		var duplicate model.SysBaseMenu
		if err := global.GVA_DB.WithContext(ctx).
			Where("name = ?", "InspectWorkOrder").
			First(&duplicate).Error; err == nil {
			mergeAndDeleteDuplicateWorkOrderMenu(ctx, keep.ID, duplicate.ID)
		}
	}
	updates := map[string]interface{}{
		"path":         "inspectionWorkOrder",
		"component":    "plugin/inspection/view/work_order.vue",
		"default_menu": false,
	}
	_ = global.GVA_DB.WithContext(ctx).
		Model(&model.SysBaseMenu{}).
		Where("name = ?", "inspectionWorkOrder").
		Updates(updates).Error
}

func mergeAndDeleteDuplicateWorkOrderMenu(ctx context.Context, keepID uint, duplicateID uint) {
	var rels []model.SysAuthorityMenu
	duplicateIDString := fmt.Sprintf("%d", duplicateID)
	keepIDString := fmt.Sprintf("%d", keepID)
	_ = global.GVA_DB.WithContext(ctx).
		Table("sys_authority_menus").
		Where("sys_base_menu_id = ?", duplicateIDString).
		Find(&rels).Error
	for _, rel := range rels {
		_ = global.GVA_DB.WithContext(ctx).
			Model(&model.SysAuthorityMenu{}).
			Where("sys_base_menu_id = ? AND sys_authority_authority_id = ?", keepIDString, rel.AuthorityId).
			FirstOrCreate(&model.SysAuthorityMenu{
				MenuId:      keepIDString,
				AuthorityId: rel.AuthorityId,
			}).Error
	}
	_ = global.GVA_DB.WithContext(ctx).Where("sys_base_menu_id = ?", duplicateIDString).Delete(&model.SysAuthorityMenu{}).Error
	_ = global.GVA_DB.WithContext(ctx).Where("sys_base_menu_id = ?", duplicateID).Delete(&model.SysBaseMenuParameter{}).Error
	_ = global.GVA_DB.WithContext(ctx).Where("sys_base_menu_id = ?", duplicateID).Delete(&model.SysBaseMenuBtn{}).Error
	_ = global.GVA_DB.WithContext(ctx).Where("sys_menu_id = ?", duplicateID).Delete(&model.SysAuthorityBtn{}).Error
	_ = global.GVA_DB.WithContext(ctx).Delete(&model.SysBaseMenu{}, duplicateID).Error
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
