package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/device"

type DeviceCategoryListResponse struct {
	List  []device.DeviceCategory `json:"list"`
	Total int64                   `json:"total"`
	Page  int                     `json:"page"`
	Size  int                     `json:"pageSize"`
}

type DeviceModelListResponse struct {
	List  []device.DeviceModel `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"pageSize"`
}

type FirmwareVersionListResponse struct {
	List  []device.FirmwareVersion `json:"list"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"pageSize"`
}

type ModelFirmwareRelListResponse struct {
	List  []device.ModelFirmwareRel `json:"list"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"pageSize"`
}

type FirmwareTagListResponse struct {
	List  []device.FirmwareTag `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"pageSize"`
}

type FirmwareVersionLogListResponse struct {
	List  []device.FirmwareVersionLog `json:"list"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Size  int                         `json:"pageSize"`
}
