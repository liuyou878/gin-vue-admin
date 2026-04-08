package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type DeviceCategorySearch struct {
	Name   string `json:"name" form:"name"`     // 类别名称
	Code   string `json:"code" form:"code"`     // 类别编码
	Status *int   `json:"status" form:"status"` // 状态
	commonReq.PageInfo
}

type DeviceModelSearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"` // 设备类别ID
	ModelCode  string `json:"modelCode" form:"modelCode"`   // 型号编码
	ModelName  string `json:"modelName" form:"modelName"`   // 型号名称
	Status     *int   `json:"status" form:"status"`         // 状态
	commonReq.PageInfo
}

type FirmwareVersionSearch struct {
	VersionCode     string      `json:"versionCode" form:"versionCode"`           // 版本号
	VersionName     string      `json:"versionName" form:"versionName"`           // 版本名称
	Status          string      `json:"status" form:"status"`                     // 状态
	PublishStatus   string      `json:"publishStatus" form:"publishStatus"`       // 发布状态
	IsLatest        *bool       `json:"isLatest" form:"isLatest"`                 // 是否最新版本
	IsStable        *bool       `json:"isStable" form:"isStable"`                 // 是否稳定版本
	UploadedAtRange []time.Time `json:"uploadedAtRange" form:"uploadedAtRange[]"` // 上传时间范围
	commonReq.PageInfo
}

type ModelFirmwareRelSearch struct {
	ModelID       uint   `json:"modelId" form:"modelId"`             // 型号ID
	FirmwareID    uint   `json:"firmwareId" form:"firmwareId"`       // 固件版本ID
	TestResult    string `json:"testResult" form:"testResult"`       // 测试结果
	IsRecommended *bool  `json:"isRecommended" form:"isRecommended"` // 是否推荐版本
	commonReq.PageInfo
}

type FirmwareTagSearch struct {
	TagCode string `json:"tagCode" form:"tagCode"` // 标签编码
	TagName string `json:"tagName" form:"tagName"` // 标签名称
	Status  *int   `json:"status" form:"status"`   // 状态
	commonReq.PageInfo
}

type FirmwareVersionLogSearch struct {
	FirmwareID     uint        `json:"firmwareId" form:"firmwareId"`           // 固件版本ID
	ModelID        uint        `json:"modelId" form:"modelId"`                 // 型号ID
	Action         string      `json:"action" form:"action"`                   // 动作
	OperateAtRange []time.Time `json:"operateAtRange" form:"operateAtRange[]"` // 操作时间范围
	commonReq.PageInfo
}

type SetFirmwareTagsRequest struct {
	FirmwareID uint   `json:"firmwareId" binding:"required"` // 固件版本ID
	TagIDs     []uint `json:"tagIds" binding:"required"`     // 标签ID列表
}

type ChangeFirmwareVersionStatusRequest struct {
	ID       uint   `json:"id" binding:"required"`     // 固件版本ID
	Status   string `json:"status" binding:"required"` // 目标状态
	Operator string `json:"operator"`                  // 操作人
	Content  string `json:"content"`                   // 日志内容
}

type SetModelFirmwareRecommendedRequest struct {
	ID       uint   `json:"id" binding:"required"` // 型号固件关系ID
	Operator string `json:"operator"`              // 操作人
	Content  string `json:"content"`               // 日志内容
}

type PublishFirmwareVersionRequest struct {
	ID       uint   `json:"id" binding:"required"` // 固件版本ID
	Operator string `json:"operator"`              // 操作人
	Content  string `json:"content"`               // 日志内容
}

type SetFirmwareStableRequest struct {
	ID       uint   `json:"id" binding:"required"` // 固件版本ID
	Stable   bool   `json:"stable"`                // 是否稳定版本
	Operator string `json:"operator"`              // 操作人
	Content  string `json:"content"`               // 日志内容
}

type VoidFirmwareVersionRequest struct {
	ID         uint   `json:"id" binding:"required"` // 固件版本ID
	Operator   string `json:"operator"`              // 操作人
	VoidReason string `json:"voidReason"`            // 作废原因
	Content    string `json:"content"`               // 日志内容
}

type SetModelFirmwareTestResultRequest struct {
	ID         uint       `json:"id" binding:"required"`         // 型号固件关系ID
	TestResult string     `json:"testResult" binding:"required"` // 测试结果
	Tester     string     `json:"tester"`                        // 测试人
	TestedAt   *time.Time `json:"testedAt"`                      // 测试时间
	Operator   string     `json:"operator"`                      // 操作人
	Content    string     `json:"content"`                       // 日志内容
}
