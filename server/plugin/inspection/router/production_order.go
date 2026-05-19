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
		group.POST("createProductionOrder", apiProductionOrder.CreateProductionOrder)
		group.DELETE("deleteProductionOrder", apiProductionOrder.DeleteProductionOrder)
		group.PUT("updateProductionOrder", apiProductionOrder.UpdateProductionOrder)
	}
	{
		group := private.Group("productionOrder")
		group.GET("findProductionOrder", apiProductionOrder.FindProductionOrder)
		group.GET("getProductionOrderList", apiProductionOrder.GetProductionOrderList)
	}
}
