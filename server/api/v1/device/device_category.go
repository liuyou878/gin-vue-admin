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

type DeviceCategoryApi struct{}

// CreateDeviceCategory 创建设备类别
func (a *DeviceCategoryApi) CreateDeviceCategory(c *gin.Context) {
	var category deviceModel.DeviceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceCategoryService.CreateDeviceCategory(&category); err != nil {
		global.GVA_LOG.Error("创建设备类别失败!", zap.Error(err))
		response.FailWithMessage("创建设备类别失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDeviceCategory 删除设备类别
func (a *DeviceCategoryApi) DeleteDeviceCategory(c *gin.Context) {
	id := c.Query("ID")
	if err := deviceCategoryService.DeleteDeviceCategory(id); err != nil {
		global.GVA_LOG.Error("删除设备类别失败!", zap.Error(err))
		response.FailWithMessage("删除设备类别失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDeviceCategoryByIds 批量删除设备类别
func (a *DeviceCategoryApi) DeleteDeviceCategoryByIds(c *gin.Context) {
	var ids commonReq.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceCategoryService.DeleteDeviceCategoryByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除设备类别失败!", zap.Error(err))
		response.FailWithMessage("批量删除设备类别失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateDeviceCategory 更新设备类别
func (a *DeviceCategoryApi) UpdateDeviceCategory(c *gin.Context) {
	var category deviceModel.DeviceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := deviceCategoryService.UpdateDeviceCategory(category); err != nil {
		global.GVA_LOG.Error("更新设备类别失败!", zap.Error(err))
		response.FailWithMessage("更新设备类别失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDeviceCategory 获取设备类别详情
func (a *DeviceCategoryApi) FindDeviceCategory(c *gin.Context) {
	id := c.Query("ID")
	category, err := deviceCategoryService.GetDeviceCategory(id)
	if err != nil {
		global.GVA_LOG.Error("查询设备类别失败!", zap.Error(err))
		response.FailWithMessage("查询设备类别失败:"+err.Error(), c)
		return
	}
	response.OkWithData(category, c)
}

// GetDeviceCategoryList 获取设备类别列表
func (a *DeviceCategoryApi) GetDeviceCategoryList(c *gin.Context) {
	var pageInfo deviceReq.DeviceCategorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := deviceCategoryService.GetDeviceCategoryInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取设备类别列表失败!", zap.Error(err))
		response.FailWithMessage("获取设备类别列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
