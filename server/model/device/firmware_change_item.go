package device

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// FirmwareChangeItem 固件变更项表
type FirmwareChangeItem struct {
	global.GVA_MODEL
	FirmwareID uint            `json:"firmwareId" form:"firmwareId" gorm:"column:firmware_id;comment:固件版本ID;not null"`                                     // 固件版本ID
	ChangeType string          `json:"changeType" form:"changeType" gorm:"column:change_type;comment:变更类型:feature/fix/optimize/breaking;size:32;not null"` // 变更类型
	Title      string          `json:"title" form:"title" gorm:"column:title;comment:变更标题;size:150;not null"`                                              // 变更标题
	Content    string          `json:"content" form:"content" gorm:"column:content;comment:变更内容;type:text"`                                                // 变更内容
	Sort       int             `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0;not null"`                                                  // 排序
	Firmware   FirmwareVersion `json:"firmware" gorm:"foreignKey:FirmwareID"`                                                                              // 固件信息
}

// TableName 固件变更项表
func (FirmwareChangeItem) TableName() string {
	return "alpha_firmware_change_items"
}
