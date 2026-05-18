package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var InspectionItem = new(inspectionItem)

type inspectionItem struct{}

func (r *inspectionItem) Init(private *gin.RouterGroup) {
	{
		group := private.Group("inspectionItem").Use(middleware.OperationRecord())
		group.POST("createItem", apiInspectionItem.CreateItem)
		group.DELETE("deleteItem", apiInspectionItem.DeleteItem)
		group.DELETE("deleteItemByIds", apiInspectionItem.DeleteItemByIds)
		group.PUT("updateItem", apiInspectionItem.UpdateItem)
	}
	{
		group := private.Group("inspectionItem")
		group.GET("findItem", apiInspectionItem.FindItem)
		group.GET("getItemList", apiInspectionItem.GetItemList)
	}
}
