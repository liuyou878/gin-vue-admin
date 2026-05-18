package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type InspectionItem struct {
	global.GVA_MODEL
	Name       string   `json:"name" gorm:"column:name;comment:检测项目名称;size:100;not null"`
	ResultType string   `json:"resultType" gorm:"column:result_type;comment:结果类型(pass_fail/number/both);size:20;not null"`
	Unit       string   `json:"unit" gorm:"column:unit;comment:单位;size:20"`
	MinValue   *float64 `json:"minValue" gorm:"column:min_value;comment:合格范围下限"`
	MaxValue   *float64 `json:"maxValue" gorm:"column:max_value;comment:合格范围上限"`
	Remark     string   `json:"remark" gorm:"column:remark;comment:备注说明;size:500"`
}

func (InspectionItem) TableName() string {
	return "inspection_items"
}
