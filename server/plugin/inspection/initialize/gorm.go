package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.InspectionItem),
		new(model.InspectionTemplate),
		new(model.InspectionTemplateItem),
		new(model.ProductionOrder),
		new(model.ProductionOrderDevice),
		new(model.InspectionDeviceResult),
	)
	if err != nil {
		err = errors.Wrap(err, "检测插件表注册失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}

	seedInspectionItems(ctx)
	seedDefaultTemplate(ctx)
}

func seedInspectionItems(ctx context.Context) {
	var count int64
	if err := global.GVA_DB.WithContext(ctx).Model(&model.InspectionItem{}).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		return
	}

	items := []model.InspectionItem{
		{Name: "老化测试", ResultType: "pass_fail"},
		{Name: "防水/气密性", ResultType: "pass_fail"},
		{Name: "蜂鸣器检查", ResultType: "pass_fail"},
		{Name: "接收机设置", ResultType: "pass_fail"},
		{Name: "电台RX移动站", ResultType: "pass_fail"},
		{Name: "电台TX基准站", ResultType: "both", Unit: "V", MinValue: f64ptr(12.0), MaxValue: f64ptr(15.0)},
		{Name: "蓝牙/WIFI", ResultType: "pass_fail"},
		{Name: "文件传输", ResultType: "pass_fail"},
		{Name: "USB充电/电量检查", ResultType: "number", Unit: "%", MinValue: f64ptr(85), MaxValue: f64ptr(100)},
		{Name: "底摄标定", ResultType: "pass_fail"},
		{Name: "相机精度", ResultType: "pass_fail"},
		{Name: "静态测量", ResultType: "pass_fail"},
		{Name: "RTK测量", ResultType: "pass_fail"},
		{Name: "外观质量", ResultType: "pass_fail"},
		{Name: "时间注册码", ResultType: "pass_fail"},
		{Name: "电子围栏", ResultType: "pass_fail"},
	}
	for _, item := range items {
		global.GVA_DB.WithContext(ctx).Create(&item)
	}
}

func seedDefaultTemplate(ctx context.Context) {
	var count int64
	if err := global.GVA_DB.WithContext(ctx).Model(&model.InspectionTemplate{}).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		return
	}

	var items []model.InspectionItem
	global.GVA_DB.WithContext(ctx).Find(&items)

	tmpl := model.InspectionTemplate{
		Name:            "G3X标准检测",
		ProductName:     "GNSS接收机（RTK）",
		Model:           "G3X",
		FirmwareVersion: "UM980-11833",
		Status:          1,
	}
	if err := global.GVA_DB.WithContext(ctx).Create(&tmpl).Error; err != nil {
		return
	}

	for i, item := range items {
		global.GVA_DB.WithContext(ctx).Create(&model.InspectionTemplateItem{
			TemplateID: tmpl.ID,
			ItemID:     item.ID,
			Sort:       i + 1,
		})
	}
}

func f64ptr(v float64) *float64 {
	return &v
}
