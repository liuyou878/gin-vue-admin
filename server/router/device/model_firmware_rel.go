package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ModelFirmwareRelRouter struct{}

// InitModelFirmwareRelRouter 初始化型号固件关系路由
func (r *ModelFirmwareRelRouter) InitModelFirmwareRelRouter(Router *gin.RouterGroup) {
	modelFirmwareRelRouter := Router.Group("modelFirmware").Use(middleware.OperationRecord())
	modelFirmwareRelRouterWithoutRecord := Router.Group("modelFirmware")
	{
		modelFirmwareRelRouter.POST("createModelFirmwareRel", modelFirmwareRelApi.CreateModelFirmwareRel)             // 创建型号固件关系
		modelFirmwareRelRouter.DELETE("deleteModelFirmwareRel", modelFirmwareRelApi.DeleteModelFirmwareRel)           // 删除型号固件关系
		modelFirmwareRelRouter.DELETE("deleteModelFirmwareRelByIds", modelFirmwareRelApi.DeleteModelFirmwareRelByIds) // 批量删除型号固件关系
		modelFirmwareRelRouter.PUT("updateModelFirmwareRel", modelFirmwareRelApi.UpdateModelFirmwareRel)              // 更新型号固件关系
		modelFirmwareRelRouter.POST("setModelFirmwareRecommended", modelFirmwareRelApi.SetModelFirmwareRecommended)   // 设置推荐版本
		modelFirmwareRelRouter.POST("setModelFirmwareTestResult", modelFirmwareRelApi.SetModelFirmwareTestResult)     // 设置测试结果
	}
	{
		modelFirmwareRelRouterWithoutRecord.GET("findModelFirmwareRel", modelFirmwareRelApi.FindModelFirmwareRel)       // 获取型号固件关系详情
		modelFirmwareRelRouterWithoutRecord.GET("getModelFirmwareRelList", modelFirmwareRelApi.GetModelFirmwareRelList) // 获取型号固件关系列表
	}
}
