package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceModelApi struct{}

// CreateDeviceModel 创建设备型号
func (a *DeviceModelApi) CreateDeviceModel(c *gin.Context) {
	var model deviceModel.DeviceModel
	if err := c.ShouldBindJSON(&model); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceModelService.CreateDeviceModel(&model); err != nil {
		global.GVA_LOG.Error("创建设备型号失败!", zap.Error(err))
		response.FailWithMessage("创建设备型号失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDeviceModel 删除设备型号
func (a *DeviceModelApi) DeleteDeviceModel(c *gin.Context) {
	id := c.Query("ID")
	if err := deviceModelService.DeleteDeviceModel(id); err != nil {
		global.GVA_LOG.Error("删除设备型号失败!", zap.Error(err))
		response.FailWithMessage("删除设备型号失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDeviceModelByIds 批量删除设备型号
func (a *DeviceModelApi) DeleteDeviceModelByIds(c *gin.Context) {
	var ids commonReq.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceModelService.DeleteDeviceModelByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除设备型号失败!", zap.Error(err))
		response.FailWithMessage("批量删除设备型号失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateDeviceModel 更新设备型号
func (a *DeviceModelApi) UpdateDeviceModel(c *gin.Context) {
	var model deviceModel.DeviceModel
	if err := c.ShouldBindJSON(&model); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceModelService.UpdateDeviceModel(model); err != nil {
		global.GVA_LOG.Error("更新设备型号失败!", zap.Error(err))
		response.FailWithMessage("更新设备型号失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDeviceModel 获取设备型号详情
func (a *DeviceModelApi) FindDeviceModel(c *gin.Context) {
	id := c.Query("ID")
	model, err := deviceModelService.GetDeviceModel(id)
	if err != nil {
		global.GVA_LOG.Error("查询设备型号失败!", zap.Error(err))
		response.FailWithMessage("查询设备型号失败:"+err.Error(), c)
		return
	}
	response.OkWithData(model, c)
}

// GetDeviceModelList 获取设备型号列表
func (a *DeviceModelApi) GetDeviceModelList(c *gin.Context) {
	var pageInfo deviceReq.DeviceModelSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := deviceModelService.GetDeviceModelInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取设备型号列表失败!", zap.Error(err))
		response.FailWithMessage("获取设备型号列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
