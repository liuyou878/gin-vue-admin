package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/router"
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	router.Router.InspectionItem.Init(private)
	router.Router.Template.Init(private)
	router.Router.ProductionOrder.Init(private)
	router.Router.WorkOrder.Init(private)
}
