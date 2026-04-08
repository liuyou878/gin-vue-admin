package device

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FirmwareVersionLog 固件版本日志表
type FirmwareVersionLog struct {
	global.GVA_MODEL
	FirmwareID uint            `json:"firmwareId" form:"firmwareId" gorm:"column:firmware_id;comment:固件版本ID;not null"`                                                                                                                                                                   // 固件版本ID
	ModelID    *uint           `json:"modelId" form:"modelId" gorm:"column:model_id;comment:型号ID"`                                                                                                                                                                                       // 型号ID
	Action     string          `json:"action" form:"action" gorm:"column:action;comment:动作:upload/delete_package/bind_model/start_testing/test_pass/test_fail/fix_upload/submit_release/reject_release/publish/mark_stable/unmark_stable/void_release/set_recommended;size:32;not null"` // 动作
	FromStatus string          `json:"fromStatus" form:"fromStatus" gorm:"column:from_status;comment:原状态;size:32"`                                                                                                                                                                       // 原状态
	ToStatus   string          `json:"toStatus" form:"toStatus" gorm:"column:to_status;comment:目标状态;size:32"`                                                                                                                                                                            // 目标状态
	Operator   string          `json:"operator" form:"operator" gorm:"column:operator;comment:操作人;size:100"`                                                                                                                                                                             // 操作人
	OperateAt  *time.Time      `json:"operateAt" form:"operateAt" gorm:"column:operate_at;comment:操作时间"`                                                                                                                                                                                 // 操作时间
	Content    string          `json:"content" form:"content" gorm:"column:content;comment:日志内容;type:text"`                                                                                                                                                                              // 日志内容
	Firmware   FirmwareVersion `json:"firmware" gorm:"foreignKey:FirmwareID"`                                                                                                                                                                                                            // 固件信息
	Model      *DeviceModel    `json:"model" gorm:"foreignKey:ModelID"`                                                                                                                                                                                                                  // 型号信息
}

// TableName 固件版本日志表
func (FirmwareVersionLog) TableName() string {
	return "alpha_firmware_version_logs"
}
