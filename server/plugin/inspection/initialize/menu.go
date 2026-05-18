package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	parent := model.SysBaseMenu{
		ParentId: 0,
		Path:     "inspection",
		Name:     "inspection",
		Hidden:   false,
		Component: "view/about/index.vue",
		Sort:     7,
		Meta: model.Meta{
			Title: "检测管理",
			Icon:  "monitor",
		},
	}
	child1 := model.SysBaseMenu{
		ParentId:  0,
		Path:      "inspectionItem",
		Name:      "inspectionItem",
		Hidden:    false,
		Component: "plugin/inspection/view/inspection_item.vue",
		Sort:      1,
		Meta: model.Meta{
			Title: "检测项管理",
			Icon:  "list",
		},
	}
	child2 := model.SysBaseMenu{
		ParentId:  0,
		Path:      "inspectionTemplate",
		Name:      "inspectionTemplate",
		Hidden:    false,
		Component: "plugin/inspection/view/template.vue",
		Sort:      2,
		Meta: model.Meta{
			Title: "检测模板",
			Icon:  "documentation",
		},
	}
	utils.RegisterMenus(parent, child1, child2)
}
