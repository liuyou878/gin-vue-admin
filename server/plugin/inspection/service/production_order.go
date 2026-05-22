package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model/request"
	"gorm.io/gorm"
)

var ProductionOrder = new(productionOrderSvc)

type productionOrderSvc struct{}

type submitDeviceInfoPayload struct {
	FirmwareVersion          string `json:"firmwareVersion"`
	MainboardFirmwareVersion string `json:"mainboardFirmwareVersion"`
	Device                   struct {
		FirmwareVersion          string `json:"firmwareVersion"`
		MainboardFirmwareVersion string `json:"mainboardFirmwareVersion"`
	} `json:"device"`
}

type normalizedSubmitDevice struct {
	SN                       string
	Model                    string
	PNCode                   string
	FirmwareVersion          string
	MainboardFirmwareVersion string
	TimeLicense              string
	RegionLicense            string
	NtripCode                string
	DeviceInfo               string
}

func (s *productionOrderSvc) CreateProductionOrder(req *request.CreateProductionOrder) error {
	now := time.Now()
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		po := model.ProductionOrder{
			MONumber:                 req.MONumber,
			TemplateID:               req.TemplateID,
			ProductName:              req.ProductName,
			Model:                    req.Model,
			FirmwareVersion:          req.FirmwareVersion,
			MainboardFirmwareVersion: req.MainboardFirmwareVersion,
			PNCode:                   req.PNCode,
			InstrumentCategory:       req.InstrumentCategory,
			Status:                   0,
			SubmitDate:               &now,
			Remark:                   req.Remark,
		}
		if err := tx.Create(&po).Error; err != nil {
			return err
		}
		for i, sn := range req.SNs {
			if sn == "" {
				continue
			}
			device := model.ProductionOrderDevice{
				ProductionOrderID: po.ID,
				SN:                sn,
				LineNumber:        i + 1,
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *productionOrderSvc) DeleteProductionOrder(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var po model.ProductionOrder
		if err := tx.Where("id = ?", id).First(&po).Error; err != nil {
			return err
		}
		if po.Status >= 1 {
			return errors.New("已确认的订单不允许删除")
		}
		return tx.Unscoped().Delete(&po).Error
	})
}

func (s *productionOrderSvc) ForceDeleteProductionOrder(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var po model.ProductionOrder
		if err := tx.Where("id = ?", id).First(&po).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&po).Error
	})
}

func (s *productionOrderSvc) UpdateProductionOrder(req *request.UpdateProductionOrder) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var po model.ProductionOrder
		if err := tx.Where("id = ?", req.ID).First(&po).Error; err != nil {
			return err
		}

		if req.Status != nil && *req.Status != po.Status {
			return errors.New("生产订单状态不支持在此页面直接修改，请到批次流程中处理")
		}

		updates := map[string]interface{}{
			"mo_number":                  req.MONumber,
			"product_name":               req.ProductName,
			"model":                      req.Model,
			"firmware_version":           req.FirmwareVersion,
			"mainboard_firmware_version": req.MainboardFirmwareVersion,
			"pn_code":                    req.PNCode,
			"instrument_category":        req.InstrumentCategory,
			"remark":                     req.Remark,
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		if err := tx.Model(&model.ProductionOrder{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *productionOrderSvc) FindProductionOrder(id string) (model.ProductionOrder, error) {
	var po model.ProductionOrder
	err := global.GVA_DB.Preload("Template").Preload("Batches.Template").Preload("Batches.Devices").Where("id = ?", id).First(&po).Error
	if err != nil {
		return po, err
	}
	var devices []model.ProductionOrderDevice
	if err := global.GVA_DB.Where("production_order_id = ?", id).Order("line_number asc").Find(&devices).Error; err != nil {
		return po, err
	}
	po.Devices = devices
	fillOrderHeaderFromDevices(&po, devices)
	po.DeviceCount = len(devices)
	for i := range po.Batches {
		po.Batches[i].DeviceCount = len(po.Batches[i].Devices)
	}
	fillBatchSummary(&po)
	var pass, fail, returned, rework, recheck int64
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", id, "pass").Count(&pass)
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", id, "fail").Count(&fail)
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", id, "returned").Count(&returned)
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", id, "rework").Count(&rework)
	global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status IN ?", id, []string{"pending_recheck", "rechecking"}).Count(&recheck)
	po.PassCount = int(pass)
	po.FailCount = int(fail)
	po.ReworkCount = int(returned + rework)
	po.RecheckCount = int(recheck)
	po.AbnormalCount = int(fail + returned + rework + recheck)
	return po, err
}

func (s *productionOrderSvc) GetProductionOrderList(search request.ProductionOrderSearch) (list []model.ProductionOrder, total int64, err error) {
	db := global.GVA_DB.Model(&model.ProductionOrder{})
	if search.MONumber != "" {
		db = db.Where("mo_number LIKE ?", "%"+search.MONumber+"%")
	}
	if search.Model != "" {
		db = db.Where("model LIKE ?", "%"+search.Model+"%")
	}
	if search.BatchNumber != "" {
		db = db.Where(
			"EXISTS (SELECT 1 FROM production_batches pb WHERE pb.production_order_id = production_orders.id AND pb.batch_number LIKE ?)",
			"%"+search.BatchNumber+"%",
		)
	}
	if search.SN != "" {
		db = db.Where(
			"EXISTS (SELECT 1 FROM production_order_devices pod WHERE pod.production_order_id = production_orders.id AND pod.sn LIKE ?)",
			"%"+search.SN+"%",
		)
	}
	if search.InstrumentCategory != "" {
		db = db.Where("instrument_category = ?", search.InstrumentCategory)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Order("id desc").Find(&list).Error
	for i := range list {
		var count, pass, fail, returned, rework, recheck int64
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ?", list[i].ID).Count(&count)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", list[i].ID, "pass").Count(&pass)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", list[i].ID, "fail").Count(&fail)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", list[i].ID, "returned").Count(&returned)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status = ?", list[i].ID, "rework").Count(&rework)
		global.GVA_DB.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ? AND status IN ?", list[i].ID, []string{"pending_recheck", "rechecking"}).Count(&recheck)
		list[i].DeviceCount = int(count)
		list[i].PassCount = int(pass)
		list[i].FailCount = int(fail)
		list[i].ReworkCount = int(returned + rework)
		list[i].RecheckCount = int(recheck)
		list[i].AbnormalCount = int(fail + returned + rework + recheck)
		if list[i].TemplateID != nil {
			var tmpl model.InspectionTemplate
			if global.GVA_DB.Where("id = ?", *list[i].TemplateID).First(&tmpl).Error == nil {
				list[i].Template = &tmpl
			}
		}
		var batches []model.ProductionBatch
		if err := global.GVA_DB.Where("production_order_id = ?", list[i].ID).Order("id asc").Find(&batches).Error; err == nil {
			list[i].Batches = batches
			fillBatchSummary(&list[i])
		}
		if needsOrderHeaderFallback(&list[i]) {
			var firstDevice model.ProductionOrderDevice
			if err := global.GVA_DB.Where("production_order_id = ?", list[i].ID).Order("line_number asc, id asc").First(&firstDevice).Error; err == nil {
				fillOrderHeaderFromDevices(&list[i], []model.ProductionOrderDevice{firstDevice})
			}
		}
	}
	return list, total, err
}

func (s *productionOrderSvc) GetSubmittedDeviceList(search request.SubmittedDeviceSearch) (list []model.SubmittedDeviceListItem, total int64, err error) {
	db := global.GVA_DB.Table("production_order_devices AS pod").
		Select("pod.id, pod.production_order_id, pod.batch_id, po.mo_number, pb.batch_number, pod.sn, pod.model, po.instrument_category, pod.pn_code, pod.firmware_version, pod.mainboard_firmware_version, pod.time_license, pod.region_license, pod.ntrip_code, pod.status, pod.return_reason, pod.return_at, pod.return_by_name, po.submitter_name, po.submit_date, pod.created_at").
		Joins("JOIN production_orders po ON po.id = pod.production_order_id").
		Joins("LEFT JOIN production_batches pb ON pb.id = pod.batch_id")

	if search.ProductionOrderID > 0 {
		db = db.Where("pod.production_order_id = ?", search.ProductionOrderID)
	}
	if search.BatchID > 0 {
		db = db.Where("pod.batch_id = ?", search.BatchID)
	}
	if search.MONumber != "" {
		db = db.Where("po.mo_number LIKE ?", "%"+search.MONumber+"%")
	}
	if search.BatchNumber != "" {
		db = db.Where("pb.batch_number LIKE ?", "%"+search.BatchNumber+"%")
	}
	if search.SN != "" {
		db = db.Where("pod.sn LIKE ?", "%"+search.SN+"%")
	}
	if search.Model != "" {
		db = db.Where("pod.model LIKE ?", "%"+search.Model+"%")
	}
	if search.Status != "" {
		db = db.Where("pod.status = ?", search.Status)
	}
	if search.Statuses != "" {
		db = db.Where("pod.status IN ?", splitDeviceStatuses(search.Statuses))
	}

	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Order("pod.id desc").Limit(search.PageSize).Offset(search.PageSize * (search.Page - 1)).Scan(&list).Error
	return list, total, err
}

func (s *productionOrderSvc) ConfirmReworkDone(req *request.ConfirmReworkDone, operatorID uint, operatorName string) error {
	ids := make([]uint, 0, len(req.DeviceIDs)+1)
	if req.DeviceID > 0 {
		ids = append(ids, req.DeviceID)
	}
	ids = append(ids, req.DeviceIDs...)
	if len(ids) == 0 {
		return errors.New("请选择返工设备")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var devices []model.ProductionOrderDevice
		if err := tx.Where("id IN ?", ids).Find(&devices).Error; err != nil {
			return err
		}
		if len(devices) != len(ids) {
			return errors.New("部分设备不存在")
		}
		for _, device := range devices {
			if device.Status != "rework" {
				return fmt.Errorf("设备 %s 不是返工中，不能确认返工完成", device.SN)
			}
		}

		for _, device := range devices {
			if req.Remark != "" {
				if err := tx.Model(&device).Update("return_reason", req.Remark).Error; err != nil {
					return err
				}
				device.ReturnReason = req.Remark
			}
			if err := updateDeviceStatusWithLog(tx, device, "pending_recheck", firstNonEmpty(req.Remark, "生产确认返工完成"), operatorID, operatorName); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *productionOrderSvc) ConfirmReworkReceived(req *request.ConfirmReworkReceived, operatorID uint, operatorName string) error {
	ids := make([]uint, 0, len(req.DeviceIDs)+1)
	if req.DeviceID > 0 {
		ids = append(ids, req.DeviceID)
	}
	ids = append(ids, req.DeviceIDs...)
	if len(ids) == 0 {
		return errors.New("请选择待接收返工设备")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var devices []model.ProductionOrderDevice
		if err := tx.Where("id IN ?", ids).Find(&devices).Error; err != nil {
			return err
		}
		if len(devices) != len(ids) {
			return errors.New("部分设备不存在")
		}
		for _, device := range devices {
			if device.Status != "returned" {
				return fmt.Errorf("设备 %s 不是待生产接收状态，不能确认接收返工", device.SN)
			}
		}

		for _, device := range devices {
			reason := firstNonEmpty(req.Remark, "生产确认接收返工")
			if err := updateDeviceStatusWithLog(tx, device, "rework", reason, operatorID, operatorName); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *productionOrderSvc) GetDeviceStatusLogs(deviceID string) ([]model.ProductionOrderDeviceStatusLog, error) {
	var logs []model.ProductionOrderDeviceStatusLog
	err := global.GVA_DB.Where("production_order_device_id = ?", deviceID).Order("id asc").Find(&logs).Error
	return logs, err
}

func (s *productionOrderSvc) FindSubmittedDevice(id string) (model.ProductionOrderDevice, error) {
	var device model.ProductionOrderDevice
	err := global.GVA_DB.Preload("Batch").Preload("ProductionOrder").Where("id = ?", id).First(&device).Error
	return device, err
}

func (s *productionOrderSvc) DeleteSubmittedDevice(id string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var device model.ProductionOrderDevice
		if err := tx.Where("id = ?", id).First(&device).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("production_order_device_id = ?", device.ID).Delete(&model.InspectionDeviceResult{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&device).Error; err != nil {
			return err
		}

		if device.BatchID != nil {
			var batchDeviceCount int64
			if err := tx.Model(&model.ProductionOrderDevice{}).Where("batch_id = ?", *device.BatchID).Count(&batchDeviceCount).Error; err != nil {
				return err
			}
			if batchDeviceCount == 0 {
				if err := tx.Unscoped().Delete(&model.ProductionBatch{}, "id = ?", *device.BatchID).Error; err != nil {
					return err
				}
			}
		}

		var orderDeviceCount int64
		if err := tx.Model(&model.ProductionOrderDevice{}).Where("production_order_id = ?", device.ProductionOrderID).Count(&orderDeviceCount).Error; err != nil {
			return err
		}
		if orderDeviceCount == 0 {
			if err := tx.Unscoped().Delete(&model.ProductionOrder{}, "id = ?", device.ProductionOrderID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// SubmitDeviceData 生产工具提交全量数据
func (s *productionOrderSvc) SubmitDeviceData(req *request.SubmitDeviceData, submitterID uint, submitterName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		devices := normalizeSubmitDevices(req)
		if len(devices) == 0 {
			return errors.New("至少需要提交一台设备")
		}

		var po model.ProductionOrder
		err := tx.Unscoped().Where("mo_number = ?", req.MONumber).First(&po).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil && po.DeletedAt.Valid {
			if err := tx.Unscoped().Delete(&po).Error; err != nil {
				return err
			}
			err = gorm.ErrRecordNotFound
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			header := buildOrderHeaderFromSubmit(req, devices[0])
			po = model.ProductionOrder{
				MONumber:                 req.MONumber,
				ProductName:              header.ProductName,
				Model:                    header.Model,
				FirmwareVersion:          header.FirmwareVersion,
				MainboardFirmwareVersion: header.MainboardFirmwareVersion,
				PNCode:                   header.PNCode,
				InstrumentCategory:       header.InstrumentCategory,
				Status:                   0,
				SubmitterID:              &submitterID,
				SubmitterName:            submitterName,
				SubmitDate:               &now,
			}
			if err := tx.Create(&po).Error; err != nil {
				return err
			}
		} else {
			if po.Status != 0 && !submitOnlyReworkDevices(tx, po.ID, devices) {
				return fmt.Errorf("生产订单 %s 已确认或已进入检测流程，不允许继续提交设备数据", req.MONumber)
			}
			updates := map[string]interface{}{
				"submitter_id":   submitterID,
				"submitter_name": submitterName,
				"submit_date":    time.Now(),
			}
			for key, value := range fillMissingOrderHeaderFromSubmit(&po, req, devices[0]) {
				updates[key] = value
			}
			if err := tx.Model(&po).Updates(updates).Error; err != nil {
				return err
			}
		}

		var batchID *uint
		if req.BatchNumber != "" {
			var batch model.ProductionBatch
			err := tx.Where("production_order_id = ? AND batch_number = ?", po.ID, req.BatchNumber).First(&batch).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				batch = model.ProductionBatch{
					ProductionOrderID: po.ID,
					BatchNumber:       req.BatchNumber,
					Status:            0,
				}
				if err := tx.Create(&batch).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			batchID = &batch.ID
		}

		for i, item := range devices {
			if item.SN == "" {
				continue
			}
			var existing model.ProductionOrderDevice
			if err := tx.Where("sn = ?", item.SN).First(&existing).Error; err == nil {
				if existing.ProductionOrderID != po.ID {
					return fmt.Errorf("SN %s 已存在于其他生产订单中", item.SN)
				}
				nextStatus := "pending"
				if existing.Status == "rework" {
					nextStatus = "pending_recheck"
				}
				updates := map[string]interface{}{
					"model":                      item.Model,
					"pn_code":                    item.PNCode,
					"time_license":               item.TimeLicense,
					"region_license":             item.RegionLicense,
					"ntrip_code":                 item.NtripCode,
					"line_number":                i + 1,
					"device_info":                item.DeviceInfo,
					"firmware_version":           firstNonEmpty(item.FirmwareVersion, extractFirmwareVersion(item.DeviceInfo)),
					"mainboard_firmware_version": firstNonEmpty(item.MainboardFirmwareVersion, extractMainboardFirmwareVersion(item.DeviceInfo)),
					"status":                     nextStatus,
					"return_reason":              "",
					"return_at":                  nil,
					"return_by_id":               nil,
					"return_by_name":             "",
				}
				if batchID != nil {
					updates["batch_id"] = *batchID
				}
				if err := tx.Model(&existing).Updates(updates).Error; err != nil {
					return err
				}
				if existing.Status != nextStatus {
					reason := "生产工具更新设备数据"
					if existing.Status == "rework" && nextStatus == "pending_recheck" {
						reason = "生产工具重新提交返工设备"
					}
					updated := existing
					updated.Model = item.Model
					updated.PNCode = item.PNCode
					updated.FirmwareVersion = firstNonEmpty(item.FirmwareVersion, extractFirmwareVersion(item.DeviceInfo))
					updated.MainboardFirmwareVersion = firstNonEmpty(item.MainboardFirmwareVersion, extractMainboardFirmwareVersion(item.DeviceInfo))
					updated.TimeLicense = item.TimeLicense
					updated.RegionLicense = item.RegionLicense
					updated.NtripCode = item.NtripCode
					updated.LineNumber = i + 1
					updated.DeviceInfo = item.DeviceInfo
					if batchID != nil {
						updated.BatchID = batchID
					}
					if err := updateDeviceStatusLogOnly(tx, updated, nextStatus, reason, submitterID, submitterName); err != nil {
						return err
					}
				}
				continue
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			device := model.ProductionOrderDevice{
				ProductionOrderID:        po.ID,
				BatchID:                  batchID,
				SN:                       item.SN,
				Model:                    item.Model,
				PNCode:                   item.PNCode,
				FirmwareVersion:          firstNonEmpty(item.FirmwareVersion, extractFirmwareVersion(item.DeviceInfo)),
				MainboardFirmwareVersion: firstNonEmpty(item.MainboardFirmwareVersion, extractMainboardFirmwareVersion(item.DeviceInfo)),
				TimeLicense:              item.TimeLicense,
				RegionLicense:            item.RegionLicense,
				NtripCode:                item.NtripCode,
				LineNumber:               i + 1,
				DeviceInfo:               item.DeviceInfo,
				Status:                   "pending",
			}
			if err := tx.Create(&device).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// AssignBatch 分配设备到批次
func (s *productionOrderSvc) AssignBatch(req *request.AssignBatch, operatorID uint, operatorName string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var batch model.ProductionBatch
		if err := tx.Where("id = ?", req.BatchID).First(&batch).Error; err != nil {
			return errors.New("批次不存在")
		}
		for _, sn := range req.SNs {
			var device model.ProductionOrderDevice
			if err := tx.Where("sn = ?", sn).First(&device).Error; err != nil {
				return errors.New("SN不存在: " + sn)
			}
			if device.ProductionOrderID != batch.ProductionOrderID {
				return errors.New("SN不属于该生产号: " + sn)
			}
			tx.Model(&device).Update("batch_id", req.BatchID)
		}
		return createBatchFlowLog(tx, batch, "生产分批", operatorID, operatorName)
	})
}

func (s *productionOrderSvc) ScanAssignBatch(req *request.ScanAssignBatch, operatorID uint, operatorName string) error {
	sns := normalizeSNList(req.SNs)
	if len(sns) == 0 {
		return errors.New("请先扫码加入设备")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order model.ProductionOrder
		if err := tx.Where("id = ?", req.ProductionOrderID).First(&order).Error; err != nil {
			return errors.New("生产订单不存在")
		}

		batchNumber := strings.TrimSpace(req.BatchNumber)
		if batchNumber == "" {
			var err error
			batchNumber, err = nextBatchNumber(tx, order)
			if err != nil {
				return err
			}
		}

		var batch model.ProductionBatch
		err := tx.Where("production_order_id = ? AND batch_number = ?", req.ProductionOrderID, batchNumber).First(&batch).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			batch = model.ProductionBatch{
				ProductionOrderID: req.ProductionOrderID,
				BatchNumber:       batchNumber,
				Status:            0,
			}
			if err := tx.Create(&batch).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if batch.Status != 0 {
			return errors.New("该批次已派检或已进入检测流程，不能继续扫码加入设备")
		}

		for _, sn := range sns {
			var device model.ProductionOrderDevice
			if err := tx.Where("production_order_id = ? AND sn = ?", req.ProductionOrderID, sn).First(&device).Error; err != nil {
				return fmt.Errorf("SN %s 不存在于当前生产订单", sn)
			}
			if device.BatchID != nil && *device.BatchID != batch.ID {
				var oldBatch model.ProductionBatch
				if err := tx.Where("id = ?", *device.BatchID).First(&oldBatch).Error; err == nil {
					return fmt.Errorf("SN %s 已在批次 %s 中", sn, oldBatch.BatchNumber)
				}
				return fmt.Errorf("SN %s 已在其他批次中", sn)
			}
			if device.Status != "pending" {
				return fmt.Errorf("SN %s 当前状态为 %s，不能加入生产批次", sn, device.Status)
			}
			if err := tx.Model(&device).Update("batch_id", batch.ID).Error; err != nil {
				return err
			}
		}
		return createBatchFlowLog(tx, batch, "生产分批", operatorID, operatorName)
	})
}

func createBatchFlowLog(tx *gorm.DB, batch model.ProductionBatch, action string, operatorID uint, operatorName string) error {
	var existing model.ProductionBatchStatusLog
	err := tx.Where("production_batch_id = ? AND from_status = ? AND to_status = ? AND action = ?", batch.ID, batch.Status, batch.Status, action).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(&model.ProductionBatchStatusLog{
			ProductionBatchID: batch.ID,
			FromStatus:        batch.Status,
			ToStatus:          batch.Status,
			Action:            action,
			OperatorID:        &operatorID,
			OperatorName:      operatorName,
		}).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func nextBatchNumber(tx *gorm.DB, order model.ProductionOrder) (string, error) {
	dateText := time.Now().Format("20060102")
	prefix := fmt.Sprintf("%s-%s-", order.MONumber, dateText)
	var count int64
	if err := tx.Model(&model.ProductionBatch{}).
		Where("production_order_id = ? AND batch_number LIKE ?", order.ID, prefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}

	for seq := int(count) + 1; seq < 1000; seq++ {
		batchNumber := fmt.Sprintf("%s%02d", prefix, seq)
		var exists int64
		if err := tx.Model(&model.ProductionBatch{}).
			Where("production_order_id = ? AND batch_number = ?", order.ID, batchNumber).
			Count(&exists).Error; err != nil {
			return "", err
		}
		if exists == 0 {
			return batchNumber, nil
		}
	}
	return "", errors.New("今日批次流水号已用完")
}

func normalizeSNList(sns []string) []string {
	seen := map[string]struct{}{}
	list := make([]string, 0, len(sns))
	for _, sn := range sns {
		value := strings.TrimSpace(sn)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		list = append(list, value)
	}
	return list
}

func extractFirmwareVersion(deviceInfo string) string {
	if deviceInfo == "" {
		return ""
	}

	var payload submitDeviceInfoPayload
	if err := json.Unmarshal([]byte(deviceInfo), &payload); err != nil {
		return ""
	}
	if payload.Device.FirmwareVersion != "" {
		return payload.Device.FirmwareVersion
	}
	return payload.FirmwareVersion
}

func normalizeSubmitDevices(req *request.SubmitDeviceData) []normalizedSubmitDevice {
	devices := make([]normalizedSubmitDevice, 0)
	for _, item := range req.Devices {
		sn := strings.TrimSpace(item.SN)
		if sn == "" {
			continue
		}
		sn = strings.TrimSpace(sn)
		devices = append(devices, normalizedSubmitDevice{
			SN:                       sn,
			Model:                    item.Model,
			PNCode:                   item.PNCode,
			FirmwareVersion:          item.FirmwareVersion,
			MainboardFirmwareVersion: item.MainboardFirmwareVersion,
			TimeLicense:              item.TimeLicense,
			RegionLicense:            item.RegionLicense,
			NtripCode:                item.NtripCode,
			DeviceInfo:               item.DeviceInfo,
		})
	}
	return devices
}

func submitOnlyReworkDevices(tx *gorm.DB, productionOrderID uint, devices []normalizedSubmitDevice) bool {
	if len(devices) == 0 {
		return false
	}
	for _, item := range devices {
		var existing model.ProductionOrderDevice
		if err := tx.Where("production_order_id = ? AND sn = ?", productionOrderID, item.SN).First(&existing).Error; err != nil {
			return false
		}
		if existing.Status != "rework" {
			return false
		}
	}
	return true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func buildOrderHeaderFromSubmit(req *request.SubmitDeviceData, device normalizedSubmitDevice) model.ProductionOrder {
	modelName := firstNonEmpty(device.Model, extractModelFromDeviceInfo(device.DeviceInfo))
	return model.ProductionOrder{
		ProductName:              modelName,
		Model:                    modelName,
		FirmwareVersion:          firstNonEmpty(device.FirmwareVersion, extractFirmwareVersion(device.DeviceInfo)),
		MainboardFirmwareVersion: firstNonEmpty(device.MainboardFirmwareVersion, extractMainboardFirmwareVersion(device.DeviceInfo)),
		PNCode:                   firstNonEmpty(device.PNCode, extractPNCodeFromDeviceInfo(device.DeviceInfo)),
		InstrumentCategory:       strings.TrimSpace(req.InstrumentCategory),
	}
}

func fillMissingOrderHeaderFromSubmit(order *model.ProductionOrder, req *request.SubmitDeviceData, device normalizedSubmitDevice) map[string]interface{} {
	header := buildOrderHeaderFromSubmit(req, device)
	updates := map[string]interface{}{}
	if strings.TrimSpace(order.ProductName) == "" && header.ProductName != "" {
		updates["product_name"] = header.ProductName
	}
	if strings.TrimSpace(order.Model) == "" && header.Model != "" {
		updates["model"] = header.Model
	}
	if strings.TrimSpace(order.FirmwareVersion) == "" && header.FirmwareVersion != "" {
		updates["firmware_version"] = header.FirmwareVersion
	}
	if strings.TrimSpace(order.MainboardFirmwareVersion) == "" && header.MainboardFirmwareVersion != "" {
		updates["mainboard_firmware_version"] = header.MainboardFirmwareVersion
	}
	if strings.TrimSpace(order.PNCode) == "" && header.PNCode != "" {
		updates["pn_code"] = header.PNCode
	}
	if strings.TrimSpace(order.InstrumentCategory) == "" && header.InstrumentCategory != "" {
		updates["instrument_category"] = header.InstrumentCategory
	}
	return updates
}

func needsOrderHeaderFallback(order *model.ProductionOrder) bool {
	return strings.TrimSpace(order.ProductName) == "" ||
		strings.TrimSpace(order.Model) == "" ||
		strings.TrimSpace(order.FirmwareVersion) == "" ||
		strings.TrimSpace(order.MainboardFirmwareVersion) == "" ||
		strings.TrimSpace(order.PNCode) == ""
}

func fillOrderHeaderFromDevices(order *model.ProductionOrder, devices []model.ProductionOrderDevice) {
	if len(devices) == 0 {
		return
	}
	for _, device := range devices {
		if strings.TrimSpace(order.ProductName) == "" {
			order.ProductName = firstNonEmpty(device.Model, extractModelFromDeviceInfo(device.DeviceInfo))
		}
		if strings.TrimSpace(order.Model) == "" {
			order.Model = firstNonEmpty(device.Model, extractModelFromDeviceInfo(device.DeviceInfo))
		}
		if strings.TrimSpace(order.FirmwareVersion) == "" {
			order.FirmwareVersion = firstNonEmpty(device.FirmwareVersion, extractFirmwareVersion(device.DeviceInfo))
		}
		if strings.TrimSpace(order.MainboardFirmwareVersion) == "" {
			order.MainboardFirmwareVersion = firstNonEmpty(device.MainboardFirmwareVersion, extractMainboardFirmwareVersion(device.DeviceInfo))
		}
		if strings.TrimSpace(order.PNCode) == "" {
			order.PNCode = firstNonEmpty(device.PNCode, extractPNCodeFromDeviceInfo(device.DeviceInfo))
		}
		if !needsOrderHeaderFallback(order) {
			return
		}
	}
}

func extractModelFromDeviceInfo(deviceInfo string) string {
	if deviceInfo == "" {
		return ""
	}

	var payload struct {
		Model  string `json:"model"`
		Device struct {
			FullType string `json:"fullType"`
			Model    string `json:"model"`
		} `json:"device"`
	}
	if err := json.Unmarshal([]byte(deviceInfo), &payload); err != nil {
		return ""
	}
	return firstNonEmpty(payload.Device.FullType, payload.Device.Model, payload.Model)
}

func extractPNCodeFromDeviceInfo(deviceInfo string) string {
	if deviceInfo == "" {
		return ""
	}

	var payload struct {
		PNCode string `json:"pnCode"`
		Device struct {
			PN string `json:"pn"`
		} `json:"device"`
	}
	if err := json.Unmarshal([]byte(deviceInfo), &payload); err != nil {
		return ""
	}
	return firstNonEmpty(payload.PNCode, payload.Device.PN)
}

func extractMainboardFirmwareVersion(deviceInfo string) string {
	if deviceInfo == "" {
		return ""
	}

	var payload struct {
		MainboardFirmwareVersion string `json:"mainboardFirmwareVersion"`
		Device                   struct {
			MainboardFirmwareVersion string `json:"mainboardFirmwareVersion"`
		} `json:"device"`
	}
	if err := json.Unmarshal([]byte(deviceInfo), &payload); err != nil {
		return ""
	}
	return firstNonEmpty(payload.Device.MainboardFirmwareVersion, payload.MainboardFirmwareVersion)
}

func fillBatchSummary(order *model.ProductionOrder) {
	order.BatchCount = len(order.Batches)
	if len(order.Batches) == 0 {
		order.BatchSummary = "-"
		return
	}
	order.Status = aggregateOrderStatus(order.Status, order.Batches)
	if len(order.Batches) == 1 {
		order.BatchSummary = order.Batches[0].BatchNumber
		return
	}

	names := make([]string, 0, len(order.Batches))
	for _, batch := range order.Batches {
		if batch.BatchNumber != "" {
			names = append(names, batch.BatchNumber)
		}
	}
	if len(names) == 0 {
		order.BatchSummary = fmt.Sprintf("多批次(%d)", len(order.Batches))
		return
	}
	if len(names) <= 2 {
		order.BatchSummary = strings.Join(names, ", ")
		return
	}
	order.BatchSummary = fmt.Sprintf("%s 等%d个批次", names[0], len(order.Batches))
}

func aggregateOrderStatus(current int, batches []model.ProductionBatch) int {
	status := current
	for _, batch := range batches {
		if batch.Status > status {
			status = batch.Status
		}
	}
	return status
}
