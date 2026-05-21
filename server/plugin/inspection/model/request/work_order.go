package request

type StartInspection struct {
	ID uint `json:"ID" binding:"required"`
}

type AssignBatchTemplate struct {
	ID         uint `json:"ID" binding:"required"`
	TemplateID uint `json:"templateID" binding:"required"`
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

type CompleteInspection struct {
	ID uint `json:"ID" binding:"required"`
}

type InspectionBatchSearch struct {
	MONumber string `json:"moNumber" form:"moNumber"`
	Model    string `json:"model" form:"model"`
	Status   *int   `json:"status" form:"status"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
