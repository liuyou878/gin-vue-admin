package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ProductionOrder 生产订单（生产号）
type ProductionOrder struct {
	global.GVA_MODEL
	MONumber           string               `json:"moNumber" gorm:"column:mo_number;uniqueIndex;size:100;not null;comment:生产订单号"`
	TemplateID         *uint                `json:"templateID" gorm:"column:template_id;comment:关联检测模板ID"`
	Template           *InspectionTemplate  `json:"template" gorm:"foreignKey:TemplateID"`
	ProductName        string               `json:"productName" gorm:"column:product_name;size:100;comment:产品名称"`
	Model              string               `json:"model" gorm:"column:model;size:100;comment:型号"`
	FirmwareVersion    string               `json:"firmwareVersion" gorm:"column:firmware_version;size:100;comment:主板固件版本"`
	InstrumentCategory string               `json:"instrumentCategory" gorm:"column:instrument_category;size:50;comment:仪器类别"`
	Status             int                  `json:"status" gorm:"column:status;default:0;comment:状态(0=待确认,1=待检测,2=检测中,3=已完成)"`
	SubmitterID        *uint                `json:"submitterID" gorm:"column:submitter_id;comment:提交人ID"`
	SubmitterName      string               `json:"submitterName" gorm:"column:submitter_name;size:100;comment:提交人姓名"`
	InspectorID        *uint                `json:"inspectorID" gorm:"column:inspector_id;comment:检测人ID"`
	InspectorName      string               `json:"inspectorName" gorm:"column:inspector_name;size:100;comment:检测人姓名"`
	InspectionDate     *time.Time           `json:"inspectionDate" gorm:"column:inspection_date;comment:检测日期"`
	SubmitDate         *time.Time           `json:"submitDate" gorm:"column:submit_date;comment:提交日期"`
	Remark             string               `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	Batches            []ProductionBatch    `json:"batches" gorm:"foreignKey:ProductionOrderID"`
	DeviceCount        int                  `json:"deviceCount" gorm:"-"`
	PassCount          int                  `json:"passCount" gorm:"-"`
	FailCount          int                  `json:"failCount" gorm:"-"`
}

func (ProductionOrder) TableName() string {
	return "production_orders"
}

// ProductionBatch 生产批次
type ProductionBatch struct {
	global.GVA_MODEL
	ProductionOrderID uint                   `json:"productionOrderID" gorm:"column:production_order_id;index;not null;comment:生产订单ID"`
	BatchNumber       string                 `json:"batchNumber" gorm:"column:batch_number;size:100;not null;comment:批次号"`
	Status            int                    `json:"status" gorm:"column:status;default:0;comment:状态(0=待确认,1=待检测,2=检测中,3=已完成)"`
	Devices           []ProductionOrderDevice `json:"devices" gorm:"foreignKey:BatchID"`
	DeviceCount       int                    `json:"deviceCount" gorm:"-"`
}

func (ProductionBatch) TableName() string {
	return "production_batches"
}

// ProductionOrderDevice 生产订单设备（序列号）
type ProductionOrderDevice struct {
	global.GVA_MODEL
	ProductionOrderID uint    `json:"productionOrderID" gorm:"column:production_order_id;index;not null;comment:生产订单ID"`
	BatchID           *uint   `json:"batchID" gorm:"column:batch_id;index;comment:批次ID(可空)"`
	Batch             *ProductionBatch `json:"batch" gorm:"foreignKey:BatchID"`
	SN                string  `json:"sn" gorm:"column:sn;size:100;uniqueIndex;not null;comment:机身码"`
	Model             string  `json:"model" gorm:"column:model;size:100;comment:设备型号"`
	PNCode            string  `json:"pnCode" gorm:"column:pn_code;size:20;comment:PN码"`
	FirmwareVersion   string  `json:"firmwareVersion" gorm:"column:firmware_version;size:100;comment:固件版本"`
	TimeLicense       string  `json:"timeLicense" gorm:"column:time_license;size:200;comment:时间注册码"`
	RegionLicense     string  `json:"regionLicense" gorm:"column:region_license;size:200;comment:围栏注册码"`
	NtripCode         string  `json:"ntripCode" gorm:"column:ntrip_code;size:200;comment:Ntrip码"`
	LineNumber        int     `json:"lineNumber" gorm:"column:line_number;default:0;comment:行号"`
	Status            string  `json:"status" gorm:"column:status;size:20;default:pending;comment:设备检测状态(pending/pass/fail)"`
	DeviceInfo        string  `json:"deviceInfo" gorm:"column:device_info;type:text;comment:GETALL完整信息JSON"`
}

func (ProductionOrderDevice) TableName() string {
	return "production_order_devices"
}

func (p *ProductionOrder) AfterDelete(tx *gorm.DB) error {
	if err := tx.Where("production_order_id = ?", p.ID).Delete(&ProductionBatch{}).Error; err != nil {
		return err
	}
	return tx.Where("production_order_id = ?", p.ID).Delete(&ProductionOrderDevice{}).Error
}
