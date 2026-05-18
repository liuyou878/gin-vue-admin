package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/api"

var (
	Router            = new(router)
	apiInspectionItem = api.Api.InspectionItem
	apiTemplate       = api.Api.Template
)

type router struct {
	InspectionItem inspectionItem
	Template       templateRouter
}
