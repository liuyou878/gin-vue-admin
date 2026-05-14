-- 为「版本流程管理」菜单添加功能按钮权限
-- 执行前请确认菜单名称与实际数据库一致

INSERT INTO sys_base_menu_btns (created_at, updated_at, deleted_at, name, `desc`, sys_base_menu_id)
SELECT NOW(), NOW(), NULL, t.name, t.desc, m.id
FROM (
  SELECT 'add' AS name, '新增固件' AS `desc`
  UNION ALL SELECT 'edit', '编辑信息'
  UNION ALL SELECT 'updatePackage', '更新包'
  UNION ALL SELECT 'startTest', '开始测试'
  UNION ALL SELECT 'submitTestResult', '测试结果'
  UNION ALL SELECT 'reject', '驳回'
  UNION ALL SELECT 'publish', '发布'
  UNION ALL SELECT 'setRecommended', '设为当前推荐'
  UNION ALL SELECT 'voidRelease', '下架'
  UNION ALL SELECT 'onShelf', '上架'
  UNION ALL SELECT 'remove', '移除'
  UNION ALL SELECT 'deleteRelation', '移除关联'
  UNION ALL SELECT 'viewLog', '查看日志'
  UNION ALL SELECT 'download', '下载'
  UNION ALL SELECT 'viewPublic', '公开下载页'
) t
JOIN sys_base_menus m ON m.name = 'deviceFirmwareWorkflow'
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menu_btns b
  WHERE b.sys_base_menu_id = m.id AND b.name = t.name
);
