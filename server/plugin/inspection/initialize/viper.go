package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/inspection/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Config config.Config

func Viper() {
	err := global.GVA_VP.UnmarshalKey("inspection", &Config)
	if err != nil {
		err = errors.Wrap(err, "初始化检测配置文件失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
