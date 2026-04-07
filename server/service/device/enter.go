package device

type ServiceGroup struct {
	DeviceCategoryService
	DeviceModelService
	FirmwareVersionService
	ModelFirmwareRelService
	FirmwareChangeItemService
	FirmwareVersionLogService
	FirmwareTagService
	FirmwareVersionTagRelService
}
