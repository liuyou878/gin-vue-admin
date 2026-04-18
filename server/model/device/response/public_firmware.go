package response

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/device"
)

type PublicFirmwareDownloadItem struct {
	RelationID    uint                   `json:"relationId"`
	Category      device.DeviceCategory  `json:"category"`
	Model         device.DeviceModel     `json:"model"`
	Firmware      device.FirmwareVersion `json:"firmware"`
	IsRecommended bool                   `json:"isRecommended"`
	PackageSize   int64                  `json:"packageSize"`
}

type PublicFirmwareDownloadPackageItem struct {
	LogID         uint                   `json:"logId"`
	RelationID    uint                   `json:"relationId"`
	Category      device.DeviceCategory  `json:"category"`
	Model         device.DeviceModel     `json:"model"`
	Firmware      device.FirmwareVersion `json:"firmware"`
	Action        string                 `json:"action"`
	OperateAt     *time.Time             `json:"operateAt"`
	IsRecommended bool                   `json:"isRecommended"`
	PackageSize   int64                  `json:"packageSize"`
}

type PublicFirmwareDownloadPageResponse struct {
	Categories         []device.DeviceCategory             `json:"categories"`
	Models             []device.DeviceModel                `json:"models"`
	SelectedCategoryID uint                                `json:"selectedCategoryId"`
	SelectedModelID    uint                                `json:"selectedModelId"`
	PrimaryType        string                              `json:"primaryType"`
	Current            *PublicFirmwareDownloadItem         `json:"current"`
	Stable             *PublicFirmwareDownloadItem         `json:"stable"`
	Latest             *PublicFirmwareDownloadItem         `json:"latest"`
	History            []PublicFirmwareDownloadItem        `json:"history"`
	Packages           []PublicFirmwareDownloadPackageItem `json:"packages"`
}
