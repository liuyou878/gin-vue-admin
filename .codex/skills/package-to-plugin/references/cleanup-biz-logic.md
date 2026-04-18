---
name: cleanup-biz-logic
description: 规范在将模块迁移为插件后，如何从 gorm_biz.go 和 router_biz.go 中安全地清理原有代码
---

# 核心原则
在将业务模块（如 `example`）迁移为插件后，需要在主项目的初始化文件中移除原有引用。**此操作必须极其谨慎，严禁影响其他业务模块。**

# 1. 清理 `server/initialize/gorm_biz.go`

## 目标
仅移除 `db.AutoMigrate()` 参数列表中属于目标模块的 Model。

## 操作规范
1.  **定位**: 找到 `initBizModel` 函数中的 `db.AutoMigrate(...)` 调用。
2.  **识别**: 识别出属于目标模块的 Model（通常以模块名开头或在模块包下）。
3.  **删除**: 仅删除该 Model 的参数传递。
4.  **保留**: **绝对不要** 删除 `db.AutoMigrate` 函数调用本身，也**绝对不要** 修改其他模块的 Model。

## 示例 (移除 `example` 模块)

**修改前**:
```go
func initBizModel(db *gorm.DB) error {
	err := db.AutoMigrate(
		&system.SysApi{},
		&example.ExaCustomer{}, // 待删除
		&other.OtherModel{},
	)
	return err
}
```

**修改后**:
```go
func initBizModel(db *gorm.DB) error {
	err := db.AutoMigrate(
		&system.SysApi{},
        // &example.ExaCustomer{}, 已删除
		&other.OtherModel{},
	)
	return err
}
```

# 2. 清理 `server/initialize/router_biz.go`

## 目标
仅移除 `initBizRouter` 函数中属于目标模块的路由初始化调用。

## 操作规范
1.  **定位**: 找到 `initBizRouter` 函数。
2.  **识别**: 识别出调用目标模块路由初始化方法（如 `exampleRouter.InitApiRouter`）的代码行。
3.  **删除**: 仅删除这些特定的调用行。
4.  **保留**: **绝对不要** 修改 `routers` 变量的定义或其他模块的路由注册。

## 示例 (移除 `example` 模块)

**修改前**:
```go
func initBizRouter(routers *server.Routers, publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	exampleRouter := routers.Example
	otherRouter := routers.Other

	exampleRouter.InitCustomerRouter(privateGroup) // 待删除
	exampleRouter.InitFileUploadAndDownloadRouter(privateGroup) // 待删除

	otherRouter.InitOrderRouter(privateGroup)
}
```

**修改后**:
```go
func initBizRouter(routers *server.Routers, publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	exampleRouter := routers.Example // 如果不再使用，也可以删除此变量声明，但需确保没有其他地方引用
	otherRouter := routers.Other

    // exampleRouter 相关调用已删除

	otherRouter.InitOrderRouter(privateGroup)
}
```
