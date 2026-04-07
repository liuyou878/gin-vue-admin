package device

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FirmwareVersionLogApi struct{}

// FindFirmwareVersionLog 获取固件日志详情
func (a *FirmwareVersionLogApi) FindFirmwareVersionLog(c *gin.Context) {
	id := c.Query("ID")
	log, err := firmwareVersionLogService.GetFirmwareVersionLog(id)
	if err != nil {
		global.GVA_LOG.Error("查询固件日志失败!", zap.Error(err))
		response.FailWithMessage("查询固件日志失败:"+err.Error(), c)
		return
	}
	response.OkWithData(log, c)
}

// GetFirmwareVersionLogList 获取固件日志列表
func (a *FirmwareVersionLogApi) GetFirmwareVersionLogList(c *gin.Context) {
	var pageInfo deviceReq.FirmwareVersionLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := firmwareVersionLogService.GetFirmwareVersionLogInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取固件日志列表失败!", zap.Error(err))
		response.FailWithMessage("获取固件日志列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
