ALTER TABLE alpha_firmware_versions
  MODIFY COLUMN status VARCHAR(32) NOT NULL DEFAULT 'pending_test' COMMENT '开发状态:pending_test/testing/tested_pass/test_failed/pending_release';

ALTER TABLE alpha_firmware_versions
  ADD COLUMN publish_status VARCHAR(32) NOT NULL DEFAULT 'unpublished' COMMENT '发布状态:unpublished/published/voided' AFTER status;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN published_by VARCHAR(100) NULL COMMENT '发布人' AFTER uploaded_at;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN published_at DATETIME(3) NULL COMMENT '发布时间' AFTER published_by;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN voided_by VARCHAR(100) NULL COMMENT '作废人' AFTER published_at;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN voided_at DATETIME(3) NULL COMMENT '作废时间' AFTER voided_by;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN void_reason TEXT NULL COMMENT '作废原因' AFTER voided_at;

ALTER TABLE alpha_firmware_versions
  ADD COLUMN package_file_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '安装包文件ID' AFTER package_name;
