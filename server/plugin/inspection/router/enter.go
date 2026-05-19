package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/api"

var (
	Router               = new(router)
	apiInspectionItem    = api.Api.InspectionItem
	apiTemplate          = api.Api.Template
	apiProductionOrder   = api.Api.ProductionOrder
)

type router struct {
	InspectionItem  inspectionItem
	Template        templateRouter
	ProductionOrder productionOrderRouter
}
