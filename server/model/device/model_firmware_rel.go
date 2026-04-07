package device

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ModelFirmwareRel 型号固件关系表
type ModelFirmwareRel struct {
	global.GVA_MODEL
	ModelID       uint            `json:"modelId" form:"modelId" gorm:"column:model_id;comment:型号ID;not null"`                                                       // 型号ID
	FirmwareID    uint            `json:"firmwareId" form:"firmwareId" gorm:"column:firmware_id;comment:固件版本ID;not null"`                                            // 固件版本ID
	IsSupported   bool            `json:"isSupported" form:"isSupported" gorm:"column:is_supported;comment:是否支持;default:true"`                                       // 是否支持
	IsRecommended bool            `json:"isRecommended" form:"isRecommended" gorm:"column:is_recommended;comment:是否推荐版本;default:false"`                              // 是否推荐版本
	TestResult    string          `json:"testResult" form:"testResult" gorm:"column:test_result;comment:测试结果:pending/testing/passed/failed;size:32;default:pending"` // 测试结果
	TestedAt      *time.Time      `json:"testedAt" form:"testedAt" gorm:"column:tested_at;comment:测试时间"`                                                             // 测试时间
	Tester        string          `json:"tester" form:"tester" gorm:"column:tester;comment:测试人;size:100"`                                                            // 测试人
	Remark        string          `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255"`                                                             // 备注
	Model         DeviceModel     `json:"model" gorm:"foreignKey:ModelID"`                                                                                           // 型号信息
	Firmware      FirmwareVersion `json:"firmware" gorm:"foreignKey:FirmwareID"`                                                                                     // 固件信息
}

// TableName 型号固件关系表
func (ModelFirmwareRel) TableName() string {
	return "alpha_model_firmware_rels"
}
