package device

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// FirmwareTag 固件标签表
type FirmwareTag struct {
	global.GVA_MODEL
	TagCode  string `json:"tagCode" form:"tagCode" gorm:"column:tag_code;comment:标签编码;size:50;not null"`  // 标签编码
	TagName  string `json:"tagName" form:"tagName" gorm:"column:tag_name;comment:标签名称;size:100;not null"` // 标签名称
	TagColor string `json:"tagColor" form:"tagColor" gorm:"column:tag_color;comment:标签颜色;size:20"`        // 标签颜色
	Sort     int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0;not null"`            // 排序
	Status   int    `json:"status" form:"status" gorm:"column:status;comment:状态:1启用 2禁用;default:1"`       // 状态
	Remark   string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255"`                // 备注
}

// TableName 固件标签表
func (FirmwareTag) TableName() string {
	return "alpha_firmware_tags"
}
