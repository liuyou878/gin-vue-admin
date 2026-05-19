package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/service"

var (
	Api                    = new(api)
	serviceInspectionItem  = service.Service.InspectionItem
	serviceTemplate        = service.Service.Template
	serviceProductionOrder = service.Service.ProductionOrder
	serviceWorkOrder       = service.Service.WorkOrder
)

type api struct {
	InspectionItem  inspectionItem
	Template        templateApi
	ProductionOrder productionOrderApi
	WorkOrder       workOrderApi
}
