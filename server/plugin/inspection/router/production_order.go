package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var ProductionOrder = new(productionOrderRouter)

type productionOrderRouter struct{}

func (r *productionOrderRouter) Init(private *gin.RouterGroup) {
	{
		group := private.Group("productionOrder").Use(middleware.OperationRecord())
		group.POST("submitDeviceData", apiProductionOrder.SubmitDeviceData)
		group.POST("createProductionOrder", apiProductionOrder.CreateProductionOrder)
		group.POST("confirmReworkReceived", apiProductionOrder.ConfirmReworkReceived)
		group.POST("confirmReworkDone", apiProductionOrder.ConfirmReworkDone)
		group.DELETE("deleteProductionOrder", apiProductionOrder.DeleteProductionOrder)
		group.DELETE("forceDeleteProductionOrder", apiProductionOrder.ForceDeleteProductionOrder)
		group.DELETE("deleteSubmittedDevice", apiProductionOrder.DeleteSubmittedDevice)
		group.PUT("updateProductionOrder", apiProductionOrder.UpdateProductionOrder)
		group.POST("assignBatch", apiProductionOrder.AssignBatch)
		group.POST("scanAssignBatch", apiProductionOrder.ScanAssignBatch)
		group.POST("createBatch", apiProductionOrder.CreateBatch)
	}
	{
		group := private.Group("productionOrder")
		group.GET("findProductionOrder", apiProductionOrder.FindProductionOrder)
		group.GET("getProductionOrderList", apiProductionOrder.GetProductionOrderList)
		group.GET("getSubmittedDeviceList", apiProductionOrder.GetSubmittedDeviceList)
		group.GET("findSubmittedDevice", apiProductionOrder.FindSubmittedDevice)
		group.GET("getDeviceStatusLogs", apiProductionOrder.GetDeviceStatusLogs)
	}
}
