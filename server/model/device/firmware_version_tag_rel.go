package device

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// FirmwareVersionTagRel 固件版本标签关系表
type FirmwareVersionTagRel struct {
	global.GVA_MODEL
	FirmwareID uint            `json:"firmwareId" form:"firmwareId" gorm:"column:firmware_id;comment:固件版本ID;not null"` // 固件版本ID
	TagID      uint            `json:"tagId" form:"tagId" gorm:"column:tag_id;comment:标签ID;not null"`                  // 标签ID
	Firmware   FirmwareVersion `json:"firmware" gorm:"foreignKey:FirmwareID"`                                          // 固件信息
	Tag        FirmwareTag     `json:"tag" gorm:"foreignKey:TagID"`                                                    // 标签信息
}

// TableName 固件版本标签关系表
func (FirmwareVersionTagRel) TableName() string {
	return "alpha_firmware_version_tag_rels"
}
