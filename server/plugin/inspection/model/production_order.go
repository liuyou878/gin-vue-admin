package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ProductionOrder 生产订单（生产号）
type ProductionOrder struct {
	global.GVA_MODEL
	MONumber                 string                  `json:"moNumber" gorm:"column:mo_number;uniqueIndex;size:100;not null;comment:生产订单号"`
	TemplateID               *uint                   `json:"templateID" gorm:"column:template_id;comment:关联检测模板ID"`
	Template                 *InspectionTemplate     `json:"template" gorm:"foreignKey:TemplateID"`
	ProductName              string                  `json:"productName" gorm:"column:product_name;size:100;comment:产品名称"`
	Model                    string                  `json:"model" gorm:"column:model;size:100;comment:型号"`
	FirmwareVersion          string                  `json:"firmwareVersion" gorm:"column:firmware_version;size:100;comment:固件版本"`
	MainboardFirmwareVersion string                  `json:"mainboardFirmwareVersion" gorm:"column:mainboard_firmware_version;size:100;comment:主板固件版本"`
	PNCode                   string                  `json:"pnCode" gorm:"column:pn_code;size:50;comment:订单PN码"`
	InstrumentCategory       string                  `json:"instrumentCategory" gorm:"column:instrument_category;size:50;comment:仪器类别"`
	Status                   int                     `json:"status" gorm:"column:status;default:0;comment:状态(0=未派检,1=待检测接收,2=检测中,3=待确认,4=已完成)"`
	SubmitterID              *uint                   `json:"submitterID" gorm:"column:submitter_id;comment:提交人ID"`
	SubmitterName            string                  `json:"submitterName" gorm:"column:submitter_name;size:100;comment:提交人姓名"`
	InspectorID              *uint                   `json:"inspectorID" gorm:"column:inspector_id;comment:检测人ID"`
	InspectorName            string                  `json:"inspectorName" gorm:"column:inspector_name;size:100;comment:检测人姓名"`
	InspectionDate           *time.Time              `json:"inspectionDate" gorm:"column:inspection_date;comment:检测日期"`
	SubmitDate               *time.Time              `json:"submitDate" gorm:"column:submit_date;comment:提交日期"`
	Remark                   string                  `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	Batches                  []ProductionBatch       `json:"batches" gorm:"foreignKey:ProductionOrderID"`
	Devices                  []ProductionOrderDevice `json:"devices" gorm:"-"`
	DeviceCount              int                     `json:"deviceCount" gorm:"-"`
	PassCount                int                     `json:"passCount" gorm:"-"`
	FailCount                int                     `json:"failCount" gorm:"-"`
	ReworkCount              int                     `json:"reworkCount" gorm:"-"`
	RecheckCount             int                     `json:"recheckCount" gorm:"-"`
	AbnormalCount            int                     `json:"abnormalCount" gorm:"-"`
	BatchCount               int                     `json:"batchCount" gorm:"-"`
	BatchSummary             string                  `json:"batchSummary" gorm:"-"`
}

func (ProductionOrder) TableName() string {
	return "production_orders"
}

// ProductionBatch 生产批次
type ProductionBatch struct {
	global.GVA_MODEL
	ProductionOrderID uint                    `json:"productionOrderID" gorm:"column:production_order_id;index;not null;comment:生产订单ID"`
	BatchNumber       string                  `json:"batchNumber" gorm:"column:batch_number;size:100;not null;comment:批次号"`
	TemplateID        *uint                   `json:"templateID" gorm:"column:template_id;comment:检测模板ID"`
	Template          *InspectionTemplate     `json:"template" gorm:"foreignKey:TemplateID"`
	Status            int                     `json:"status" gorm:"column:status;default:0;comment:状态(0=未派检,1=待检测接收,2=检测中,3=待确认,4=已完成)"`
	InspectorID       *uint                   `json:"inspectorID" gorm:"column:inspector_id;comment:检测人ID"`
	InspectorName     string                  `json:"inspectorName" gorm:"column:inspector_name;size:100;comment:检测人姓名"`
	InspectionDate    *time.Time              `json:"inspectionDate" gorm:"column:inspection_date;comment:检测日期"`
	Devices           []ProductionOrderDevice `json:"devices" gorm:"foreignKey:BatchID"`
	DeviceCount       int                     `json:"deviceCount" gorm:"-"`
	PassCount         int                     `json:"passCount" gorm:"-"`
	FailCount         int                     `json:"failCount" gorm:"-"`
}

func (ProductionBatch) TableName() string {
	return "production_batches"
}

type ProductionBatchStatusLog struct {
	global.GVA_MODEL
	ProductionBatchID uint             `json:"productionBatchID" gorm:"column:production_batch_id;index;not null;comment:生产批次ID"`
	Batch             *ProductionBatch `json:"batch" gorm:"foreignKey:ProductionBatchID"`
	FromStatus        int              `json:"fromStatus" gorm:"column:from_status;comment:变更前状态"`
	ToStatus          int              `json:"toStatus" gorm:"column:to_status;comment:变更后状态"`
	Action            string           `json:"action" gorm:"column:action;size:100;comment:动作"`
	Reason            string           `json:"reason" gorm:"column:reason;size:500;comment:原因/备注"`
	OperatorID        *uint            `json:"operatorID" gorm:"column:operator_id;comment:操作人ID"`
	OperatorName      string           `json:"operatorName" gorm:"column:operator_name;size:100;comment:操作人姓名"`
}

func (ProductionBatchStatusLog) TableName() string {
	return "production_batch_status_logs"
}

type InspectionBatchListItem struct {
	ID                       uint                `json:"ID"`
	ProductionOrderID        uint                `json:"productionOrderID"`
	MONumber                 string              `json:"moNumber"`
	BatchNumber              string              `json:"batchNumber"`
	ProductName              string              `json:"productName"`
	Model                    string              `json:"model"`
	FirmwareVersion          string              `json:"firmwareVersion"`
	MainboardFirmwareVersion string              `json:"mainboardFirmwareVersion"`
	PNCode                   string              `json:"pnCode"`
	InstrumentCategory       string              `json:"instrumentCategory"`
	Status                   int                 `json:"status"`
	TemplateID               *uint               `json:"templateID"`
	Template                 *InspectionTemplate `json:"template"`
	InspectorID              *uint               `json:"inspectorID"`
	InspectorName            string              `json:"inspectorName"`
	InspectionDate           *time.Time          `json:"inspectionDate"`
	DeviceCount              int                 `json:"deviceCount"`
	PassCount                int                 `json:"passCount"`
	FailCount                int                 `json:"failCount"`
	ReworkCount              int                 `json:"reworkCount"`
	RecheckCount             int                 `json:"recheckCount"`
	RecheckingCount          int                 `json:"recheckingCount"`
	CreatedAt                time.Time           `json:"CreatedAt"`
}

type SubmittedDeviceListItem struct {
	ID                       uint       `json:"ID"`
	ProductionOrderID        uint       `json:"productionOrderID"`
	BatchID                  *uint      `json:"batchID"`
	MONumber                 string     `json:"moNumber"`
	BatchNumber              string     `json:"batchNumber"`
	SN                       string     `json:"sn"`
	Model                    string     `json:"model"`
	InstrumentCategory       string     `json:"instrumentCategory"`
	PNCode                   string     `json:"pnCode"`
	FirmwareVersion          string     `json:"firmwareVersion"`
	MainboardFirmwareVersion string     `json:"mainboardFirmwareVersion"`
	TimeLicense              string     `json:"timeLicense"`
	RegionLicense            string     `json:"regionLicense"`
	NtripCode                string     `json:"ntripCode"`
	Status                   string     `json:"status"`
	ReturnReason             string     `json:"returnReason"`
	ReturnAt                 *time.Time `json:"returnAt"`
	ReturnByName             string     `json:"returnByName"`
	SubmitterName            string     `json:"submitterName"`
	SubmitDate               *time.Time `json:"submitDate"`
	CreatedAt                time.Time  `json:"CreatedAt"`
}

// ProductionOrderDevice 生产订单设备（序列号）
type ProductionOrderDevice struct {
	global.GVA_MODEL
	ProductionOrderID        uint             `json:"productionOrderID" gorm:"column:production_order_id;index;not null;comment:生产订单ID"`
	BatchID                  *uint            `json:"batchID" gorm:"column:batch_id;index;comment:批次ID(可空)"`
	Batch                    *ProductionBatch `json:"batch" gorm:"foreignKey:BatchID"`
	ProductionOrder          *ProductionOrder `json:"productionOrder" gorm:"foreignKey:ProductionOrderID"`
	SN                       string           `json:"sn" gorm:"column:sn;size:100;uniqueIndex;not null;comment:机身码"`
	Model                    string           `json:"model" gorm:"column:model;size:100;comment:设备型号"`
	PNCode                   string           `json:"pnCode" gorm:"column:pn_code;size:20;comment:PN码"`
	FirmwareVersion          string           `json:"firmwareVersion" gorm:"column:firmware_version;size:100;comment:固件版本"`
	MainboardFirmwareVersion string           `json:"mainboardFirmwareVersion" gorm:"column:mainboard_firmware_version;size:100;comment:主板固件版本"`
	TimeLicense              string           `json:"timeLicense" gorm:"column:time_license;size:200;comment:时间注册码"`
	RegionLicense            string           `json:"regionLicense" gorm:"column:region_license;size:200;comment:围栏注册码"`
	NtripCode                string           `json:"ntripCode" gorm:"column:ntrip_code;size:200;comment:Ntrip码"`
	LineNumber               int              `json:"lineNumber" gorm:"column:line_number;default:0;comment:行号"`
	Status                   string           `json:"status" gorm:"column:status;size:20;default:pending;comment:设备状态(pending/pass/fail/returned/rework/pending_recheck/rechecking)"`
	ReturnReason             string           `json:"returnReason" gorm:"column:return_reason;size:500;comment:退回原因"`
	ReturnAt                 *time.Time       `json:"returnAt" gorm:"column:return_at;comment:退回时间"`
	ReturnByID               *uint            `json:"returnByID" gorm:"column:return_by_id;comment:退回人ID"`
	ReturnByName             string           `json:"returnByName" gorm:"column:return_by_name;size:100;comment:退回人姓名"`
	DeviceInfo               string           `json:"deviceInfo" gorm:"column:device_info;type:text;comment:GETALL完整信息JSON"`
}

func (ProductionOrderDevice) TableName() string {
	return "production_order_devices"
}

type ProductionOrderDeviceStatusLog struct {
	global.GVA_MODEL
	ProductionOrderDeviceID uint                   `json:"productionOrderDeviceID" gorm:"column:production_order_device_id;index;not null;comment:设备ID"`
	Device                  *ProductionOrderDevice `json:"device" gorm:"foreignKey:ProductionOrderDeviceID"`
	FromStatus              string                 `json:"fromStatus" gorm:"column:from_status;size:30;comment:变更前状态"`
	ToStatus                string                 `json:"toStatus" gorm:"column:to_status;size:30;not null;comment:变更后状态"`
	Reason                  string                 `json:"reason" gorm:"column:reason;size:500;comment:原因/备注"`
	OperatorID              *uint                  `json:"operatorID" gorm:"column:operator_id;comment:操作人ID"`
	OperatorName            string                 `json:"operatorName" gorm:"column:operator_name;size:100;comment:操作人姓名"`
}

func (ProductionOrderDeviceStatusLog) TableName() string {
	return "production_order_device_status_logs"
}

func (p *ProductionOrder) AfterDelete(tx *gorm.DB) error {
	if err := tx.Unscoped().Where(
		"production_order_device_id IN (?)",
		tx.Unscoped().Model(&ProductionOrderDevice{}).Select("id").Where("production_order_id = ?", p.ID),
	).Delete(&InspectionDeviceResult{}).Error; err != nil {
		return err
	}
	if err := tx.Unscoped().Where("production_order_id = ?", p.ID).Delete(&ProductionBatch{}).Error; err != nil {
		return err
	}
	return tx.Unscoped().Where("production_order_id = ?", p.ID).Delete(&ProductionOrderDevice{}).Error
}
