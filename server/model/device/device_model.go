package device

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// DeviceModel 设备型号表
type DeviceModel struct {
	global.GVA_MODEL
	CategoryID uint           `json:"categoryId" form:"categoryId" gorm:"column:category_id;comment:设备类别ID;not null"` // 设备类别ID
	ModelCode  string         `json:"modelCode" form:"modelCode" gorm:"column:model_code;comment:型号编码;size:100;not null"`
	ModelName  string         `json:"modelName" form:"modelName" gorm:"column:model_name;comment:型号名称;size:150;not null"`
	Sort       int            `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0;not null"`            // 排序
	SeriesName string         `json:"seriesName" form:"seriesName" gorm:"column:series_name;comment:系列名称;size:100"` // 系列名称
	Generation string         `json:"generation" form:"generation" gorm:"column:generation;comment:代际;size:50"`     // 代际
	Status     int            `json:"status" form:"status" gorm:"column:status;comment:状态:1启用 2禁用;default:1"`       // 状态
	Remark     string         `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255"`                // 备注
	Category   DeviceCategory `json:"category" gorm:"foreignKey:CategoryID"`                                        // 设备类别
}

// TableName 设备型号表
func (DeviceModel) TableName() string {
	return "alpha_device_models"
}
