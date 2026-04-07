package device

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	deviceModel "github.com/flipped-aurora/gin-vue-admin/server/model/device"
	deviceReq "github.com/flipped-aurora/gin-vue-admin/server/model/device/request"
	"gorm.io/gorm"
)

// CreateFirmwareTag 创建固件标签
func (s *FirmwareTagService) CreateFirmwareTag(tag *deviceModel.FirmwareTag) error {
	return global.GVA_DB.Create(tag).Error
}

// DeleteFirmwareTag 删除固件标签
func (s *FirmwareTagService) DeleteFirmwareTag(id string) error {
	return global.GVA_DB.Delete(&deviceModel.FirmwareTag{}, "id = ?", id).Error
}

// DeleteFirmwareTagByIds 批量删除固件标签
func (s *FirmwareTagService) DeleteFirmwareTagByIds(ids commonReq.IdsReq) error {
	return global.GVA_DB.Delete(&[]deviceModel.FirmwareTag{}, "id in ?", ids.Ids).Error
}

// UpdateFirmwareTag 更新固件标签
func (s *FirmwareTagService) UpdateFirmwareTag(tag deviceModel.FirmwareTag) error {
	return global.GVA_DB.Model(&deviceModel.FirmwareTag{}).Where("id = ?", tag.ID).Updates(&tag).Error
}

// GetFirmwareTag 获取固件标签详情
func (s *FirmwareTagService) GetFirmwareTag(id string) (tag deviceModel.FirmwareTag, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&tag).Error
	return
}

// GetFirmwareTagInfoList 获取固件标签分页列表
func (s *FirmwareTagService) GetFirmwareTagInfoList(info deviceReq.FirmwareTagSearch) (list []deviceModel.FirmwareTag, total int64, err error) {
	db := global.GVA_DB.Model(&deviceModel.FirmwareTag{})
	if info.TagCode != "" {
		db = db.Where("tag_code LIKE ?", "%"+info.TagCode+"%")
	}
	if info.TagName != "" {
		db = db.Where("tag_name LIKE ?", "%"+info.TagName+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}
	err = db.Order("sort asc, id desc").Find(&list).Error
	return
}

// SetFirmwareTags 设置固件版本标签
func (s *FirmwareVersionTagRelService) SetFirmwareTags(req deviceReq.SetFirmwareTagsRequest) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var firmwareCount int64
		if err := tx.Model(&deviceModel.FirmwareVersion{}).Where("id = ?", req.FirmwareID).Count(&firmwareCount).Error; err != nil {
			return err
		}
		if firmwareCount == 0 {
			return errors.New("固件版本不存在")
		}

		if len(req.TagIDs) > 0 {
			var tagCount int64
			if err := tx.Model(&deviceModel.FirmwareTag{}).Where("id in ?", req.TagIDs).Count(&tagCount).Error; err != nil {
				return err
			}
			if int(tagCount) != len(req.TagIDs) {
				return errors.New("标签数据不存在或已失效")
			}
		}

		if err := tx.Unscoped().Where("firmware_id = ?", req.FirmwareID).Delete(&deviceModel.FirmwareVersionTagRel{}).Error; err != nil {
			return err
		}

		if len(req.TagIDs) == 0 {
			return nil
		}

		rels := make([]deviceModel.FirmwareVersionTagRel, 0, len(req.TagIDs))
		for _, tagID := range req.TagIDs {
			rels = append(rels, deviceModel.FirmwareVersionTagRel{
				FirmwareID: req.FirmwareID,
				TagID:      tagID,
			})
		}
		return tx.Create(&rels).Error
	})
}
