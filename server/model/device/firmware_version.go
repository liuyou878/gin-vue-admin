package device

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FirmwareVersion 固件版本表
type FirmwareVersion struct {
	global.GVA_MODEL
	VersionCode   string                  `json:"versionCode" form:"versionCode" gorm:"column:version_code;comment:版本号;size:100;not null"`                                                          // 版本号
	VersionName   string                  `json:"versionName" form:"versionName" gorm:"column:version_name;comment:版本名称;size:150;not null"`                                                         // 版本名称
	PackageURL    string                  `json:"packageUrl" form:"packageUrl" gorm:"column:package_url;comment:安装包地址;size:500"`                                                                    // 安装包地址
	PackageName   string                  `json:"packageName" form:"packageName" gorm:"column:package_name;comment:安装包名称;size:255"`                                                                 // 安装包名称
	Checksum      string                  `json:"checksum" form:"checksum" gorm:"column:checksum;comment:校验值;size:128"`                                                                             // 校验值
	Status        string                  `json:"status" form:"status" gorm:"column:status;comment:开发状态:pending_test/testing/tested_pass/test_failed/pending_release;size:32;default:pending_test"` // 开发状态
	PublishStatus string                  `json:"publishStatus" form:"publishStatus" gorm:"column:publish_status;comment:发布状态:unpublished/published/voided;size:32;default:unpublished"`            // 发布状态
	ReleaseNote   string                  `json:"releaseNote" form:"releaseNote" gorm:"column:release_note;comment:版本说明;type:text"`                                                                 // 版本说明
	TestSummary   string                  `json:"testSummary" form:"testSummary" gorm:"column:test_summary;comment:测试总结;type:text"`                                                                 // 测试总结
	IsLatest      bool                    `json:"isLatest" form:"isLatest" gorm:"column:is_latest;comment:是否最新版本;default:false"`                                                                    // 是否最新版本
	IsStable      bool                    `json:"isStable" form:"isStable" gorm:"column:is_stable;comment:是否稳定版本;default:false"`                                                                    // 是否稳定版本
	UploadedBy    string                  `json:"uploadedBy" form:"uploadedBy" gorm:"column:uploaded_by;comment:上传人;size:100"`                                                                      // 上传人
	UploadedAt    *time.Time              `json:"uploadedAt" form:"uploadedAt" gorm:"column:uploaded_at;comment:上传时间"`                                                                              // 上传时间
	PublishedBy   string                  `json:"publishedBy" form:"publishedBy" gorm:"column:published_by;comment:发布人;size:100"`                                                                   // 发布人
	PublishedAt   *time.Time              `json:"publishedAt" form:"publishedAt" gorm:"column:published_at;comment:发布时间"`                                                                           // 发布时间
	VoidedBy      string                  `json:"voidedBy" form:"voidedBy" gorm:"column:voided_by;comment:作废人;size:100"`                                                                            // 作废人
	VoidedAt      *time.Time              `json:"voidedAt" form:"voidedAt" gorm:"column:voided_at;comment:作废时间"`                                                                                    // 作废时间
	VoidReason    string                  `json:"voidReason" form:"voidReason" gorm:"column:void_reason;comment:作废原因;type:text"`                                                                    // 作废原因
	Remark        string                  `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:255"`                                                                                    // 备注
	ChangeItems   []FirmwareChangeItem    `json:"changeItems" gorm:"foreignKey:FirmwareID"`                                                                                                         // 变更项
	Tags          []FirmwareVersionTagRel `json:"tags" gorm:"foreignKey:FirmwareID"`                                                                                                                // 标签关系
}

// TableName 固件版本表
func (FirmwareVersion) TableName() string {
	return "alpha_firmware_versions"
}
