INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(3), NOW(3), '/firmwareVersion/publishFirmwareVersion', '发布固件版本', '设备固件-固件版本', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM sys_apis WHERE path = '/firmwareVersion/publishFirmwareVersion' AND method = 'POST' AND deleted_at IS NULL
);

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(3), NOW(3), '/firmwareVersion/setFirmwareStable', '设置稳定版本', '设备固件-固件版本', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM sys_apis WHERE path = '/firmwareVersion/setFirmwareStable' AND method = 'POST' AND deleted_at IS NULL
);

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(3), NOW(3), '/firmwareVersion/voidFirmwareVersion', '作废固件版本', '设备固件-固件版本', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM sys_apis WHERE path = '/firmwareVersion/voidFirmwareVersion' AND method = 'POST' AND deleted_at IS NULL
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '888', '/firmwareVersion/publishFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/firmwareVersion/publishFirmwareVersion' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '888', '/firmwareVersion/setFirmwareStable', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/firmwareVersion/setFirmwareStable' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '888', '/firmwareVersion/voidFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/firmwareVersion/voidFirmwareVersion' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '8881', '/firmwareVersion/publishFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '8881' AND v1 = '/firmwareVersion/publishFirmwareVersion' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '8881', '/firmwareVersion/setFirmwareStable', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '8881' AND v1 = '/firmwareVersion/setFirmwareStable' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '8881', '/firmwareVersion/voidFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '8881' AND v1 = '/firmwareVersion/voidFirmwareVersion' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '9528', '/firmwareVersion/publishFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '9528' AND v1 = '/firmwareVersion/publishFirmwareVersion' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '9528', '/firmwareVersion/setFirmwareStable', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '9528' AND v1 = '/firmwareVersion/setFirmwareStable' AND v2 = 'POST'
);

INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', '9528', '/firmwareVersion/voidFirmwareVersion', 'POST'
WHERE NOT EXISTS (
  SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '9528' AND v1 = '/firmwareVersion/voidFirmwareVersion' AND v2 = 'POST'
);
