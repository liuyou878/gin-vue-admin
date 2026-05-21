package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var WorkOrder = new(workOrderRouter)

type workOrderRouter struct{}

func (r *workOrderRouter) Init(private *gin.RouterGroup) {
	{
		group := private.Group("workOrder").Use(middleware.OperationRecord())
		group.POST("assignBatchTemplate", apiWorkOrder.AssignBatchTemplate)
		group.POST("startInspection", apiWorkOrder.StartInspection)
		group.POST("saveResults", apiWorkOrder.SaveResults)
		group.POST("completeInspection", apiWorkOrder.CompleteInspection)
	}
	{
		group := private.Group("workOrder")
		group.GET("getInspectionBatchList", apiWorkOrder.GetInspectionBatchList)
		group.GET("getInspectionDetail", apiWorkOrder.GetInspectionDetail)
	}
}
