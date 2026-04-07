package device

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// DeviceCategory 设备类别表
type DeviceCategory struct {
	global.GVA_MODEL
	Name   string `json:"name" form:"name" gorm:"column:name;comment:类别名称;size:100;not null"`     // 类别名称
	Code   string `json:"code" form:"code" gorm:"column:code;comment:类别编码;size:100;not null"`     // 类别编码
	Sort   int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0;not null"`      // 排序
	Status int    `json:"status" form:"status" gorm:"column:status;comment:状态:1启用 2禁用;default:1"` // 状态
	Remark string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255"`          // 备注
}

// TableName 设备类别表
func (DeviceCategory) TableName() string {
	return "alpha_device_categories"
}
