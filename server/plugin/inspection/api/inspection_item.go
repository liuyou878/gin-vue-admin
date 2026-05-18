package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"github.com/gin-gonic/gin"
)

var InspectionItem = new(inspectionItem)

type inspectionItem struct{}

// CreateItem 新增检测项
// @Tags     InspectionItem
// @Summary  新增检测项
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateInspectionItem true "检测项信息"
// @Success  200 {object} response.Response{msg=string} "创建成功"
// @Router   /inspectionItem/createItem [post]
func (a *inspectionItem) CreateItem(c *gin.Context) {
	var req request.CreateInspectionItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	item := &model.InspectionItem{
		Name:       req.Name,
		ResultType: req.ResultType,
		Unit:       req.Unit,
		MinValue:   req.MinValue,
		MaxValue:   req.MaxValue,
		Remark:     req.Remark,
	}
	if err := serviceInspectionItem.CreateItem(item); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteItem 删除检测项
// @Tags     InspectionItem
// @Summary  删除检测项
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "检测项ID"
// @Success  200 {object} response.Response{msg=string} "删除成功"
// @Router   /inspectionItem/deleteItem [delete]
func (a *inspectionItem) DeleteItem(c *gin.Context) {
	id := c.Query("id")
	if err := serviceInspectionItem.DeleteItem(id); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteItemByIds 批量删除检测项
// @Tags     InspectionItem
// @Summary  批量删除检测项
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    ids query string true "检测项ID列表(逗号分隔)"
// @Success  200 {object} response.Response{msg=string} "批量删除成功"
// @Router   /inspectionItem/deleteItemByIds [delete]
func (a *inspectionItem) DeleteItemByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.FailWithMessage("ids不能为空", c)
		return
	}
	ids := make([]string, 0)
	for _, s := range splitIDs(idsStr) {
		if s != "" {
			ids = append(ids, s)
		}
	}
	if err := serviceInspectionItem.DeleteItemByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateItem 更新检测项
// @Tags     InspectionItem
// @Summary  更新检测项
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateInspectionItem true "检测项信息"
// @Success  200 {object} response.Response{msg=string} "更新成功"
// @Router   /inspectionItem/updateItem [put]
func (a *inspectionItem) UpdateItem(c *gin.Context) {
	var req request.UpdateInspectionItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	item := &model.InspectionItem{
		Name:       req.Name,
		ResultType: req.ResultType,
		Unit:       req.Unit,
		MinValue:   req.MinValue,
		MaxValue:   req.MaxValue,
		Remark:     req.Remark,
	}
	item.ID = req.ID
	if err := serviceInspectionItem.UpdateItem(item); err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindItem 根据ID获取检测项
// @Tags     InspectionItem
// @Summary  根据ID获取检测项
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query string true "检测项ID"
// @Success  200 {object} response.Response{data=model.InspectionItem,msg=string} "查询成功"
// @Router   /inspectionItem/findItem [get]
func (a *inspectionItem) FindItem(c *gin.Context) {
	id := c.Query("id")
	item, err := serviceInspectionItem.FindItem(id)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithData(item, c)
}

// GetItemList 获取检测项列表
// @Tags     InspectionItem
// @Summary  获取检测项列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Param    name query string false "检测项名称"
// @Param    resultType query string false "结果类型"
// @Success  200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router   /inspectionItem/getItemList [get]
func (a *inspectionItem) GetItemList(c *gin.Context) {
	var search request.InspectionItemSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInspectionItem.GetItemList(search)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "查询成功", c)
}

func splitIDs(s string) []string {
	var ids []string
	cur := ""
	for _, ch := range s {
		if ch == ',' {
			if cur != "" {
				ids = append(ids, cur)
				cur = ""
			}
		} else if ch != ' ' {
			cur += string(ch)
		}
	}
	if cur != "" {
		ids = append(ids, cur)
	}
	return ids
}
