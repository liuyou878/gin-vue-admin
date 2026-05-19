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
	SNs                []string `json:"sns" binding:"required,min=1"`
}

type UpdateProductionOrder struct {
	ID                 uint     `json:"ID" binding:"required"`
	MONumber           string   `json:"moNumber" binding:"required"`
	TemplateID         *uint    `json:"templateID"`
	ProductName        string   `json:"productName"`
	Model              string   `json:"model"`
	FirmwareVersion    string   `json:"firmwareVersion"`
	InstrumentCategory string   `json:"instrumentCategory"`
	BatchNumber        string   `json:"batchNumber"`
	Status             *int     `json:"status"`
	Remark             string   `json:"remark"`
	SNs                []string `json:"sns" binding:"required,min=1"`
}
