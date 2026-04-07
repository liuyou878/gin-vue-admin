package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/device"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	Device  device.RouterGroup
	System  system.RouterGroup
	Example example.RouterGroup
}
