package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type InspectionDeviceResult struct {
	global.GVA_MODEL
	ProductionOrderDeviceID uint            `json:"productionOrderDeviceID" gorm:"column:production_order_device_id;index;not null;comment:生产订单设备ID"`
	ItemID                  uint            `json:"itemID" gorm:"column:item_id;not null;comment:检测项ID"`
	PassResult              *bool           `json:"passResult" gorm:"column:pass_result;comment:通过/不通过"`
	NumberResult            *float64        `json:"numberResult" gorm:"column:number_result;comment:数值结果"`
	Remark                  string          `json:"remark" gorm:"column:remark;size:500;comment:单项备注"`
	Item                    *InspectionItem `json:"item" gorm:"foreignKey:ItemID"`
}

func (InspectionDeviceResult) TableName() string {
	return "inspection_device_results"
}
