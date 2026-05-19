package request

type StartInspection struct {
	ID uint `json:"ID" binding:"required"`
}

type SaveInspectionResult struct {
	ProductionOrderID uint                      `json:"productionOrderID" binding:"required"`
	DeviceStatuses    []DeviceStatusItem        `json:"deviceStatuses"`
	DeviceResults     []DeviceResultItem        `json:"deviceResults"`
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
