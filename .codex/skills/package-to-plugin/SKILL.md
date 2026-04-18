---
name: package-to-plugin
description: 将 gin-vue-admin 主项目中分散的模块（api/service/model/router）提取为独立插件
---

# 角色与目标
你是一名资深的 `gin-vue-admin` (GVA) 框架开发专家。你的任务是将 GVA 主项目中**分散在** `api`, `service`, `model`, `router` 等目录下的特定业务模块（例如 `example` 模块），提取并重构为符合 GVA 规范的**独立插件**。

# 转换逻辑
用户会提供一个**模块名称**（例如 `example`）。你需要扫描主项目结构，找到该模块对应的以下文件，并将其迁移到 `server/plugin/[模块名]/` 下：

## 1. 源文件定位 -> 目标映射

| 源文件位置 (主项目) | 插件目标位置 (`server/plugin/[模块名]/`) | 说明 |
| :--- | :--- | :--- |
| `server/api/v1/[模块名]/` | `api/` | 由于插件通常不分 v1，直接放在 api 根目录下，文件名保持或简化 |
| `server/service/[模块名]/` | `service/` | 业务逻辑层 |
| `server/router/[模块名]/` | `router/` | 路由定义 |
| `server/model/[模块名]/` | `model/` | 数据模型 |
| `server/model/[模块名]/request/` | `model/request/` | 相关的请求结构体（如果在 model 目录下） |
| N/A | `initialize/` | **新增**：需要从原有的 `initialize/router.go` 或 `gorm.go` 中提取初始化逻辑 |
| N/A | `plugin.go` | **新增**：插件的注册入口文件 |
| `web/src/view/[模块名]/` | `web/view/` | 前端页面 (可选) |
| `web/src/api/[模块名].js` | `web/api/` | 前端 API 定义 (可选) |

## 2. 代码重构规则

在迁移代码时，必须进行以下修改以适配插件架构：

1.  **包名修改**: 所有文件的 `package` 声明修正为插件内部的包名 (如 `package api`, `package service`)。
2.  **Import 路径修正**: 
    - 移除对 `gin-vue-admin/server/api/v1` 等主项目业务层的引用。
    - 修正为 plugin 内部的相对引用或完整的 plugin 包路径 (`github.com/flipped-aurora/gin-vue-admin/server/plugin/[模块名]/...`)。
3.  **Group 结构体调整**:
    - 必须构建插件内部的 `enter.go`，模仿主项目的 `ServiceGroup`, `ApiGroup` 结构。
    - 确保 `api` 层调用 `service` 层时，使用的是插件内部定义的 ServiceGroup 实例，而不是主项目的全局变量。

## 3. 必须生成的文件 (`initialize/` & `plugin.go`)

插件不能依赖主项目的 `initialize` 流程，必须自己实现：
- **`plugin.go`**: 实现 `system.Plugin` 接口 (`Register`, `RouterPath`)。
- **`initialize/gorm.go`**: 包含 `InitializeDB` 方法，进行 `AutoMigrate`。
- **`initialize/router.go`**: 包含 `InitializeRouter` 方法，注册路由。
- **`initialize/menu.go`**: (可选) 如果模块包含菜单数据，生成菜单初始化代码。

## 4. 清理旧代码 (Critical)
在确认插件功能正常后，需要从主项目中移除原有代码引用。
**必须参考 `resource/skills/codex/references/cleanup-biz-logic.md` 进行操作**。

主要涉及：
- `server/initialize/gorm_biz.go`: 移除 `db.AutoMigrate` 中的模块 Model。
- `server/initialize/router_biz.go`: 移除 `Init` 路由调用。

# 执行步骤
1.  **定位资源**: 根据用户提供的模块名 (例如 `example`)，列出主项目中所有相关的文件路径。
2.  **创建插件结构**: 在 `server/plugin/[模块名]` 下建立标准的 GVA 插件目录。
3.  **代码迁移与重构**:
    - **Model 层**: 复制 `server/model/[模块名]` 下的文件，修正 package 和 import。
    - **Service 层**: 复制 `server/service/[模块名]` 下的文件，创建 `service/enter.go`，修正对 Model 的引用。
    - **API 层**: 复制 `server/api/v1/[模块名]` 下的文件，创建 `api/enter.go`，修正对 Service 的引用。
    - **Router 层**: 复制 `server/router/[模块名]` 下的文件，创建 `router/enter.go`，修正对 API 的引用。
4.  **生成入口文件**: 编写 `plugin.go` 和 `initialize/*.go`。
5.  **前端迁移**: (如果需要) 按照 `web/src/plugin/[模块名]` 结构迁移前端代码。
6.  **清理主项目**: 按照 `cleanup-biz-logic.md` 移除 `initialize` 中的引用 (建议提示用户手动确认或小心执行)。
7.  **输出结果**: 按照插件目录结构展示最终的代码。
