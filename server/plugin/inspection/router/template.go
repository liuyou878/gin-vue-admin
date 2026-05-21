package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var Template = new(templateRouter)

type templateRouter struct{}

func (r *templateRouter) Init(private *gin.RouterGroup) {
	{
		group := private.Group("inspectionTemplate").Use(middleware.OperationRecord())
		group.POST("createTemplate", apiTemplate.CreateTemplate)
		group.POST("copyTemplate", apiTemplate.CopyTemplate)
		group.DELETE("deleteTemplate", apiTemplate.DeleteTemplate)
		group.PUT("updateTemplate", apiTemplate.UpdateTemplate)
	}
	{
		group := private.Group("inspectionTemplate")
		group.GET("findTemplate", apiTemplate.FindTemplate)
		group.GET("getTemplateList", apiTemplate.GetTemplateList)
	}
}
