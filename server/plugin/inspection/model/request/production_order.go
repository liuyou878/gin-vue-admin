package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type ProductionOrderSearch struct {
	MONumber           string `json:"moNumber" form:"moNumber"`
	Model              string `json:"model" form:"model"`
	InstrumentCategory string `json:"instrumentCategory" form:"instrumentCategory"`
	Status             *int   `json:"status" form:"status"`
	request.PageInfo
}

type CreateProductionOrder struct {
	MONumber           string   `json:"moNumber" binding:"required"`
	TemplateID         *uint    `json:"templateID"`
	ProductName        string   `json:"productName"`
	Model              string   `json:"model"`
	FirmwareVersion    string   `json:"firmwareVersion"`
	InstrumentCategory string   `json:"instrumentCategory"`
	BatchNumber        string   `json:"batchNumber"`
	Remark             string   `json:"remark"`
	SNs                []string `json:"sns"`
}

type UpdateProductionOrder struct {
	ID                 uint   `json:"ID" binding:"required"`
	MONumber           string `json:"moNumber" binding:"required"`
	TemplateID         *uint  `json:"templateID"`
	ProductName        string `json:"productName"`
	Model              string `json:"model"`
	FirmwareVersion    string `json:"firmwareVersion"`
	InstrumentCategory string `json:"instrumentCategory"`
	Status             *int   `json:"status"`
	Remark             string `json:"remark"`
}

// SubmitDeviceData 生产工具提交全量数据
type SubmitDeviceData struct {
	MONumber           string           `json:"moNumber" binding:"required"`
	BatchNumber        string           `json:"batchNumber"`
	DeviceType         string           `json:"deviceType" binding:"required"`
	InstrumentCategory string           `json:"instrumentCategory"`
	SNs                []string         `json:"sns" binding:"required,min=1"`
	PNCode             string           `json:"pnCode"`
	TimeLicense        string           `json:"timeLicense"`
	RegionLicense      string           `json:"regionLicense"`
	NtripCode          string           `json:"ntripCode"`
	DeviceInfo         string           `json:"deviceInfo"`
}

// AssignBatch 分配序列号到批次
type AssignBatch struct {
	BatchID uint     `json:"batchID" binding:"required"`
	SNs     []string `json:"sns" binding:"required,min=1"`
}

// CreateBatch 创建批次
type CreateBatch struct {
	ProductionOrderID uint   `json:"productionOrderID" binding:"required"`
	BatchNumber       string `json:"batchNumber" binding:"required"`
}
