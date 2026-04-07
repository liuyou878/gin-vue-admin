# 设备固件版本管理模块设计说明

## 1. 背景说明

当前项目里已经存在 `sysVersion` 模块，但这个模块的用途是 Gin-Vue-Admin 自带的系统快照导入导出，主要管理的是：

- 菜单
- API
- 字典

它不适合直接拿来做公司的设备固件版本管理。

本次要做的业务模块目标是：

- 按设备型号管理固件版本，不细化到单台设备
- 能看到某个型号有哪些固件版本
- 能区分某个版本是否测试中、测试通过、稳定、废弃
- 能记录每个版本改了什么、优化了什么、修复了什么
- 能给版本打标签，例如大改、修复 Bug、优化、新增功能
- 能看到某个型号当前推荐使用哪个版本
- 能区分“最新版本”和“稳定版本”

结论：

- 设备固件管理应当作为一个新的业务模块开发
- 不建议继续在现有 `sysVersion` 上硬改

## 2. 一期范围

### 2.1 一期要做的内容

- 设备类别管理
- 设备型号管理
- 固件版本管理
- 型号和固件版本关系管理
- 固件变更内容管理
- 固件标签管理
- 固件状态流转日志
- 按类别、型号、版本、状态进行查询

### 2.2 一期先不做的内容

- 单台设备固件版本跟踪
- OTA 升级任务下发
- 固件包推送
- 真正的设备升级回滚执行
- 审批流

## 3. 核心业务概念

### 3.1 设备类别

用于归类和筛选，例如：

- RTK
- 全站仪
- 手簿
- 接收机

设备类别主要用于管理和查询，不直接作为固件兼容绑定单位。

### 3.2 设备型号

设备型号是固件兼容和推荐策略的最小单位。

例如：

- G3X 一代
- G3X 二代
- TS-600

固件应该绑定到“型号”层，而不是只绑定到“类别”层。

### 3.3 固件版本

表示一次上传的固件版本记录，例如：

- `v1.0.0`
- `v1.0.1`
- `v2.3.5-beta`

每上传一个新版本，都应新增一条记录，不应直接覆盖旧版本。

### 3.4 型号-固件关系

用于描述：

- 某个型号支持哪些固件版本
- 哪个版本当前推荐
- 哪个版本已经测试通过

### 3.5 固件变更项

用于记录这个版本到底改了什么，例如：

- 修复蓝牙连接问题
- 优化 RTK 启动速度
- 新增参数配置功能

### 3.6 固件日志

用于记录固件生命周期的重要动作，例如：

- 上传版本
- 开始测试
- 测试通过
- 设置为稳定版
- 设置为推荐版
- 废弃版本

### 3.7 固件标签

用于给版本打上业务标签，便于快速筛选和归类。

例如：

- 大改
- 修复 Bug
- 优化
- 新增功能
- 兼容性调整
- 安全修复
- 性能提升

说明：

- 日志记录的是“发生了什么动作”
- 标签表达的是“这个版本属于什么类型变化”
- 两者应当并存，不能互相替代

## 4. 推荐数据模型

一期建议使用 6 张表。

## 4.1 `alpha_device_categories`

用途：设备类别表。

建议字段：

- `id`
- `name`
- `code`
- `sort`
- `status`
- `remark`
- `created_at`
- `updated_at`

说明：

- 用于维护设备大类
- 如 `RTK`、`全站仪`

## 4.2 `alpha_device_models`

用途：设备型号表。

建议字段：

- `id`
- `category_id`
- `model_code`
- `model_name`
- `series_name`
- `generation`
- `status`
- `remark`
- `created_at`
- `updated_at`

说明：

- `series_name` 可存 `G3X`
- `generation` 可存 `1代`、`2代`
- 实际固件绑定仍然在 `alpha_device_models` 这一层完成

## 4.3 `alpha_firmware_versions`

用途：固件版本主表。

建议字段：

- `id`
- `version_code`
- `version_name`
- `package_url`
- `package_name`
- `checksum`
- `status`
- `release_note`
- `test_summary`
- `is_latest`
- `is_stable`
- `uploaded_by`
- `uploaded_at`
- `remark`
- `created_at`
- `updated_at`

建议状态值：

- `draft`
- `testing`
- `tested_pass`
- `stable`
- `deprecated`

说明：

- `is_latest` 表示最新版本
- `is_stable` 表示稳定版本
- “最新”不等于“稳定”

## 4.4 `alpha_model_firmware_rels`

用途：型号和固件版本关系表。

建议字段：

- `id`
- `model_id`
- `firmware_id`
- `is_supported`
- `is_recommended`
- `is_default`
- `test_result`
- `tested_at`
- `tester`
- `remark`
- `created_at`
- `updated_at`

建议测试结果值：

- `pending`
- `testing`
- `passed`
- `failed`

说明：

- `is_recommended` 用于表示该型号当前推荐使用哪个版本
- 一个型号可以对应多个版本
- 一个版本也可以对应多个型号
- 如果 `is_default` 和 `is_recommended` 业务含义重复，一期可以先去掉 `is_default`

## 4.5 `alpha_firmware_change_items`

用途：固件变更项表。

建议字段：

- `id`
- `firmware_id`
- `change_type`
- `title`
- `content`
- `sort`
- `created_at`
- `updated_at`

建议变更类型：

- `feature`
- `fix`
- `optimize`
- `breaking`

说明：

- 用于结构化记录每个版本的变更内容
- 不建议只把所有变更写成一大段纯文本

## 4.6 `alpha_firmware_version_logs`

用途：固件状态和关键动作日志表。

建议字段：

- `id`
- `firmware_id`
- `model_id`
- `action`
- `from_status`
- `to_status`
- `operator`
- `operate_at`
- `content`
- `created_at`
- `updated_at`

建议动作值：

- `upload`
- `bind_model`
- `start_testing`
- `test_pass`
- `test_fail`
- `set_recommended`
- `set_stable`
- `unset_stable`
- `deprecate`

说明：

- `model_id` 可为空
- 如果动作只针对版本本身，可以不填 `model_id`
- 如果动作是“某型号切换推荐版本”，则应记录 `model_id`

## 4.7 `alpha_firmware_tags`

用途：固件标签定义表。

建议字段：

- `id`
- `tag_code`
- `tag_name`
- `tag_color`
- `sort`
- `status`
- `remark`
- `created_at`
- `updated_at`

建议标签值：

- `major_change`
- `fix`
- `optimize`
- `feature`
- `compatibility`
- `security`
- `performance`

说明：

- 用于统一维护版本标签
- 支持后续继续扩展，不必把标签写死在代码里

## 4.8 `alpha_firmware_version_tag_rels`

用途：固件版本和标签关系表。

建议字段：

- `id`
- `firmware_id`
- `tag_id`
- `created_at`
- `updated_at`

说明：

- 一个版本可以挂多个标签
- 一个标签也可以被多个版本复用
- 例如某个版本既可以是“修复 Bug”，也可以是“优化”

## 5. 为什么不是 4 张表，而是 6 张表

如果只做下面 4 张表：

- `alpha_device_categories`
- `alpha_device_models`
- `alpha_firmware_versions`
- `alpha_model_firmware_rels`

那么系统只能知道“当前结果”，但不能清楚回答下面这些问题：

- 这个版本到底改了什么
- 谁把这个版本从测试中改成稳定版
- 某个型号什么时候把推荐版本从 A 切换到 B

所以还需要：

- `alpha_firmware_change_items`
  用于记录版本改动内容
- `alpha_firmware_version_logs`
  用于记录过程留痕

这 6 张表才是既能看当前状态、又能看历史过程的最小可用方案。

如果还要支持“标签轴”，建议再增加：

- `alpha_firmware_tags`
- `alpha_firmware_version_tag_rels`

这样系统就可以同时支持：

- 状态轴
  例如测试中、稳定版、废弃
- 标签轴
  例如大改、修复 Bug、优化、新增功能

## 6. 关键业务规则

## 6.1 最新、稳定、推荐要分开

这三个概念不能混用。

- `latest`：当前最新版本
- `stable`：已经确认可稳定使用的版本
- `recommended`：某个型号当前推荐使用的版本

示例：

- `v1.0.5` 可能是最新版本
- `v1.0.4` 可能仍然是稳定版本
- 对于 `G3X 一代`，推荐版本可能还是 `v1.0.4`

## 6.2 新版本上传必须新增记录

上传新版本时，不要把旧版本改成新版本。

正确方式：

- 新增一条 `alpha_firmware_versions`
- 按需调整 `alpha_model_firmware_rels`
- 写入 `alpha_firmware_version_logs`

这样旧版本历史才能一直保留下来。

## 6.3 推荐版本必须按型号控制

推荐版本不能只在固件表里定义，而必须在 `alpha_model_firmware_rels` 中按型号维护。

原因：

- 同一个固件版本，可能对某个型号已经稳定
- 但对另一个型号还没有完成测试

## 6.4 变更内容建议结构化

建议做法：

- `release_note` 存版本总结
- `alpha_firmware_change_items` 存详细变更项

这样既方便展示，也方便后续统计和查询。

## 6.5 标签和日志要分开

推荐同时保留两种维度：

- `alpha_firmware_version_logs`
  记录动作和状态变化
- `alpha_firmware_tags` 与 `alpha_firmware_version_tag_rels`
  记录版本标签

示例：

- 某版本日志中记录：
  - 上传版本
  - 开始测试
  - 测试通过
  - 设置稳定版
- 同时给该版本打标签：
  - 修复 Bug
  - 优化

这样后续既能按流程查，也能按类型查。

## 7. 主要页面设计

一期建议至少包含以下页面。

## 7.1 设备类别管理

功能：

- 列表
- 新增
- 编辑
- 启用禁用

## 7.2 设备型号管理

功能：

- 列表
- 新增
- 编辑
- 按设备类别筛选
- 维护系列和代际信息

## 7.3 固件版本管理

功能：

- 查看所有固件版本
- 新增版本记录
- 维护固件包地址或上传信息
- 编辑版本说明
- 选择版本标签
- 调整版本状态
- 查看变更项
- 查看操作日志

建议列表字段：

- 版本号
- 版本名称
- 状态
- 是否最新
- 是否稳定
- 标签
- 上传人
- 上传时间

## 7.4 型号固件关系管理

功能：

- 选择某个型号
- 查看该型号支持的所有固件版本
- 设置推荐版本
- 记录测试结果
- 记录测试人和测试时间

建议列表字段：

- 型号
- 固件版本
- 测试结果
- 是否推荐
- 测试时间
- 测试人

## 7.5 固件详情页

功能：

- 查看版本基础信息
- 查看版本标签
- 查看适配型号
- 查看变更项
- 查看日志

## 8. 后端模块建议

建议作为全新的业务模块开发，不复用当前 `sysVersion`。

建议目录：

- `server/model/device`
- `server/model/device/request`
- `server/model/device/response`
- `server/service/device`
- `server/api/v1/device`
- `server/router/device`
- `server/source/device`

建议实体名称：

- `DeviceCategory`
- `DeviceModel`
- `FirmwareVersion`
- `ModelFirmwareRel`
- `FirmwareChangeItem`
- `FirmwareVersionLog`
- `FirmwareTag`
- `FirmwareVersionTagRel`

## 9. 前端模块建议

建议目录：

- `web/src/api/deviceCategory.js`
- `web/src/api/deviceModel.js`
- `web/src/api/firmwareVersion.js`
- `web/src/api/modelFirmware.js`
- `web/src/view/deviceFirmware/category/index.vue`
- `web/src/view/deviceFirmware/model/index.vue`
- `web/src/view/deviceFirmware/version/index.vue`
- `web/src/view/deviceFirmware/modelVersion/index.vue`
- `web/src/view/deviceFirmware/version/detail.vue`
- `web/src/view/deviceFirmware/tag/index.vue`

建议菜单结构：

- 设备固件管理
- 设备类别
- 设备型号
- 固件版本
- 型号固件关系
- 固件标签

## 10. 接口建议

一期建议接口如下。

## 10.1 设备类别

- `POST /deviceCategory/create`
- `PUT /deviceCategory/update`
- `DELETE /deviceCategory/delete`
- `GET /deviceCategory/find`
- `GET /deviceCategory/list`

## 10.2 设备型号

- `POST /deviceModel/create`
- `PUT /deviceModel/update`
- `DELETE /deviceModel/delete`
- `GET /deviceModel/find`
- `GET /deviceModel/list`

## 10.3 固件版本

- `POST /firmwareVersion/create`
- `PUT /firmwareVersion/update`
- `DELETE /firmwareVersion/delete`
- `GET /firmwareVersion/find`
- `GET /firmwareVersion/list`
- `POST /firmwareVersion/changeStatus`
- `GET /firmwareVersion/logs`
- `GET /firmwareVersion/changes`

## 10.4 型号固件关系

- `POST /modelFirmware/create`
- `PUT /modelFirmware/update`
- `DELETE /modelFirmware/delete`
- `GET /modelFirmware/list`
- `POST /modelFirmware/setRecommended`
- `POST /modelFirmware/setTestResult`

## 10.5 固件变更项

- `POST /firmwareChange/create`
- `PUT /firmwareChange/update`
- `DELETE /firmwareChange/delete`
- `GET /firmwareChange/list`

## 10.6 固件标签

- `POST /firmwareTag/create`
- `PUT /firmwareTag/update`
- `DELETE /firmwareTag/delete`
- `GET /firmwareTag/find`
- `GET /firmwareTag/list`
- `POST /firmwareVersion/setTags`

## 11. 推荐开发顺序

为降低风险，一期建议按下面顺序开发。

## 第一阶段：基础数据

- 建 `alpha_device_categories`
- 建 `alpha_device_models`
- 建 `alpha_firmware_versions`

目标：

- 系统可以管理设备类别、设备型号、固件版本基础数据

## 第二阶段：关系管理

- 建 `alpha_model_firmware_rels`
- 支持把固件版本绑定到型号
- 支持设置推荐版本
- 支持记录测试结果

目标：

- 系统可以回答：
- 某型号有哪些固件版本
- 某型号当前推荐哪个版本
- 哪些版本已经测试通过

## 第三阶段：过程留痕

- 建 `alpha_firmware_change_items`
- 建 `alpha_firmware_version_logs`

目标：

- 系统可以回答：
- 这个版本改了什么
- 什么时候从测试中变成稳定版
- 谁执行了这些关键操作

## 第三阶段补充：标签能力

- 建 `alpha_firmware_tags`
- 建 `alpha_firmware_version_tag_rels`
- 支持在版本编辑页选择多个标签

目标：

- 系统可以回答：
- 哪些版本属于“大改”
- 哪些版本是“修复 Bug”
- 哪些版本同时包含“优化”和“新增功能”

## 第四阶段：前端集成

- 加菜单
- 加 CRUD 页面
- 加详情页
- 加状态操作入口

## 12. 推荐状态流转

建议版本状态流转如下：

`draft -> testing -> tested_pass -> stable -> deprecated`

如果测试失败，可以按业务选择：

- `testing -> draft`
- 或 `testing -> deprecated`

建议规则：

- 每个型号同一时间只能有一个推荐版本
- 稳定版本可以按业务决定是否允许多个
- 最新版本通常应只有一个

## 13. 初版界面策略

一期建议先采用标准的 Gin-Vue-Admin CRUD 风格，不必一开始做得太重。

建议首版界面包括：

- 搜索区域
- 列表表格
- 新增编辑抽屉
- 详情抽屉
- 标签选择器
- 变更项列表
- 日志列表

这样最容易快速落地，也方便后面继续扩展。

## 14. 风险与注意事项

## 14.1 不建议复用 `sysVersion`

原因：

- 现有 `sysVersion` 是系统快照工具
- 设备固件管理是业务模块
- 两者语义完全不同

## 14.5 不建议只靠日志表达版本类型

原因：

- 日志只能看出“发生了什么动作”
- 不能直接筛选“哪些版本是优化类更新”
- 也不方便统计“本月新增功能类版本有多少”

## 14.2 不建议把“最新”和“稳定”混为一体

原因：

- 最新不一定稳定
- 稳定也不一定是最新

## 14.3 不建议只按设备类别绑定固件

原因：

- 固件兼容性一般是按型号决定的
- 类别层级太粗

## 14.4 不建议不做日志

原因：

- 业务上通常都需要查：
- 改了什么
- 什么时候改的
- 谁改的

## 15. 最终建议

一期建议按“设备型号级固件管理”来做，核心原则如下：

- 设备类别只做分类和筛选
- 固件兼容和推荐关系绑定到设备型号
- 分开管理最新、稳定、推荐三个概念
- 必须保留变更内容和过程日志
- 不复用当前 `sysVersion`

一期最合适的表结构组合为：

- `alpha_device_categories`
- `alpha_device_models`
- `alpha_firmware_versions`
- `alpha_model_firmware_rels`
- `alpha_firmware_change_items`
- `alpha_firmware_version_logs`

如果同时要支持“版本标签轴”，建议再增加：

- `alpha_firmware_tags`
- `alpha_firmware_version_tag_rels`

如果这份说明确认没有问题，下一步建议直接进入：

1. 先细化这 6 张表的字段定义
2. 再做后端 model、service、api、router
3. 最后补前端页面和菜单
