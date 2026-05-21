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
		group.POST("assignOrderTemplate", apiWorkOrder.AssignOrderTemplate)
		group.POST("startInspection", apiWorkOrder.StartInspection)
		group.POST("startRecheck", apiWorkOrder.StartRecheck)
		group.POST("saveResults", apiWorkOrder.SaveResults)
		group.POST("completeInspection", apiWorkOrder.CompleteInspection)
		group.POST("completeRecheck", apiWorkOrder.CompleteRecheck)
		group.POST("returnDevices", apiWorkOrder.ReturnDevices)
	}
	{
		group := private.Group("workOrder")
		group.GET("getInspectionBatchList", apiWorkOrder.GetInspectionBatchList)
		group.GET("getInspectionDetail", apiWorkOrder.GetInspectionDetail)
	}
}
