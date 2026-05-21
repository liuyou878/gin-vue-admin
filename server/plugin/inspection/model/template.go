package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type InspectionTemplate struct {
	global.GVA_MODEL
	Name             string `json:"name" gorm:"column:name;comment:模板名称;size:100;not null"`
	ProductName      string `json:"productName" gorm:"column:product_name;comment:产品名称;size:100"`
	Model            string `json:"model" gorm:"column:model;comment:适用型号;size:100"`
	Status           int    `json:"status" gorm:"column:status;comment:状态(1=启用,2=停用);default:1"`
	TemplateItems    []InspectionTemplateItem `json:"templateItems" gorm:"foreignKey:TemplateID"`
	ItemCount        int    `json:"itemCount" gorm:"-"`
}

func (InspectionTemplate) TableName() string {
	return "inspection_templates"
}

type InspectionTemplateItem struct {
	global.GVA_MODEL
	TemplateID uint            `json:"templateID" gorm:"column:template_id;index;not null;comment:模板ID"`
	ItemID     uint            `json:"itemID" gorm:"column:item_id;not null;comment:检测项ID"`
	Sort       int             `json:"sort" gorm:"column:sort;not null;default:0;comment:排序号"`
	Item       InspectionItem  `json:"item" gorm:"foreignKey:ItemID"`
}

func (InspectionTemplateItem) TableName() string {
	return "inspection_template_items"
}
