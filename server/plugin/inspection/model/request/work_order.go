package request

type StartInspection struct {
	ID uint `json:"ID" binding:"required"`
}

type AssignBatchTemplate struct {
	ID         uint `json:"ID" binding:"required"`
	TemplateID uint `json:"templateID" binding:"required"`
}

type AssignOrderTemplate struct {
	ProductionOrderID  uint   `json:"productionOrderID" binding:"required"`
	TemplateID         uint   `json:"templateID" binding:"required"`
	InstrumentCategory string `json:"instrumentCategory"`
}

type SaveInspectionResult struct {
	BatchID        uint               `json:"batchID" binding:"required"`
	DeviceStatuses []DeviceStatusItem `json:"deviceStatuses"`
	DeviceResults  []DeviceResultItem `json:"deviceResults"`
}

type DeviceStatusItem struct {
	DeviceID uint   `json:"deviceID"`
	Status   string `json:"status"`
}

type DeviceResultItem struct {
	DeviceID     uint     `json:"deviceID"`
	ItemID       uint     `json:"itemID"`
	PassResult   *bool    `json:"passResult"`
	NumberResult *float64 `json:"numberResult"`
	Remark       string   `json:"remark"`
}

type SaveSingleInspectionResult struct {
	BatchID      uint     `json:"batchID" binding:"required"`
	DeviceID     uint     `json:"deviceID" binding:"required"`
	ItemID       uint     `json:"itemID" binding:"required"`
	PassResult   *bool    `json:"passResult"`
	NumberResult *float64 `json:"numberResult"`
	Remark       string   `json:"remark"`
	Status       string   `json:"status"`
}

type CompleteInspection struct {
	ID uint `json:"ID" binding:"required"`
}

type ConfirmInspectionComplete struct {
	ID uint `json:"ID" binding:"required"`
}

type StartRecheck struct {
	ID       uint `json:"ID" binding:"required"`
	DeviceID uint `json:"deviceID"`
}

type CompleteRecheck struct {
	ID uint `json:"ID" binding:"required"`
}

type ReturnDevices struct {
	BatchID   uint   `json:"batchID" binding:"required"`
	DeviceIDs []uint `json:"deviceIDs" binding:"required,min=1"`
	Reason    string `json:"reason"`
}

type InspectionBatchSearch struct {
	MONumber     string `json:"moNumber" form:"moNumber"`
	Model        string `json:"model" form:"model"`
	BatchNumber  string `json:"batchNumber" form:"batchNumber"`
	SN           string `json:"sn" form:"sn"`
	Status       *int   `json:"status" form:"status"`
	DeviceStatus string `json:"deviceStatus" form:"deviceStatus"`
	Page         int    `json:"page" form:"page"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
}
