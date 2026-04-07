CREATE TABLE IF NOT EXISTS `alpha_device_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '类别名称',
  `code` varchar(100) NOT NULL COMMENT '类别编码',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态:1启用 2禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alpha_device_categories_code` (`code`),
  KEY `idx_alpha_device_categories_deleted_at` (`deleted_at`),
  KEY `idx_alpha_device_categories_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='设备类别表';

CREATE TABLE IF NOT EXISTS `alpha_device_models` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `category_id` bigint unsigned NOT NULL COMMENT '设备类别ID',
  `model_code` varchar(100) NOT NULL COMMENT '型号编码',
  `model_name` varchar(150) NOT NULL COMMENT '型号名称',
  `series_name` varchar(100) DEFAULT NULL COMMENT '系列名称',
  `generation` varchar(50) DEFAULT NULL COMMENT '代际',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态:1启用 2禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alpha_device_models_model_code` (`model_code`),
  KEY `idx_alpha_device_models_category_id` (`category_id`),
  KEY `idx_alpha_device_models_deleted_at` (`deleted_at`),
  KEY `idx_alpha_device_models_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='设备型号表';

CREATE TABLE IF NOT EXISTS `alpha_firmware_versions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `version_code` varchar(100) NOT NULL COMMENT '版本号',
  `version_name` varchar(150) NOT NULL COMMENT '版本名称',
  `package_url` varchar(500) DEFAULT NULL COMMENT '安装包地址',
  `package_name` varchar(255) DEFAULT NULL COMMENT '安装包名称',
  `checksum` varchar(128) DEFAULT NULL COMMENT '校验值',
  `status` varchar(32) NOT NULL DEFAULT 'draft' COMMENT '状态:draft/testing/tested_pass/stable/deprecated',
  `release_note` text COMMENT '版本说明',
  `test_summary` text COMMENT '测试总结',
  `is_latest` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否最新版本',
  `is_stable` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否稳定版本',
  `uploaded_by` varchar(100) DEFAULT NULL COMMENT '上传人',
  `uploaded_at` datetime(3) DEFAULT NULL COMMENT '上传时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_alpha_firmware_versions_version_code` (`version_code`),
  KEY `idx_alpha_firmware_versions_status` (`status`),
  KEY `idx_alpha_firmware_versions_is_latest` (`is_latest`),
  KEY `idx_alpha_firmware_versions_is_stable` (`is_stable`),
  KEY `idx_alpha_firmware_versions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='固件版本表';

CREATE TABLE IF NOT EXISTS `alpha_model_firmware_rels` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `model_id` bigint unsigned NOT NULL COMMENT '型号ID',
  `firmware_id` bigint unsigned NOT NULL COMMENT '固件版本ID',
  `is_supported` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否支持',
  `is_recommended` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否推荐版本',
  `test_result` varchar(32) NOT NULL DEFAULT 'pending' COMMENT '测试结果:pending/testing/passed/failed',
  `tested_at` datetime(3) DEFAULT NULL COMMENT '测试时间',
  `tester` varchar(100) DEFAULT NULL COMMENT '测试人',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alpha_model_firmware_rels_model_firmware` (`model_id`,`firmware_id`),
  KEY `idx_alpha_model_firmware_rels_firmware_id` (`firmware_id`),
  KEY `idx_alpha_model_firmware_rels_test_result` (`test_result`),
  KEY `idx_alpha_model_firmware_rels_is_recommended` (`is_recommended`),
  KEY `idx_alpha_model_firmware_rels_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='型号固件关系表';

CREATE TABLE IF NOT EXISTS `alpha_firmware_change_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `firmware_id` bigint unsigned NOT NULL COMMENT '固件版本ID',
  `change_type` varchar(32) NOT NULL COMMENT '变更类型:feature/fix/optimize/breaking',
  `title` varchar(150) NOT NULL COMMENT '变更标题',
  `content` text COMMENT '变更内容',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_alpha_firmware_change_items_firmware_id` (`firmware_id`),
  KEY `idx_alpha_firmware_change_items_change_type` (`change_type`),
  KEY `idx_alpha_firmware_change_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='固件变更项表';

CREATE TABLE IF NOT EXISTS `alpha_firmware_version_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `firmware_id` bigint unsigned NOT NULL COMMENT '固件版本ID',
  `model_id` bigint unsigned DEFAULT NULL COMMENT '型号ID',
  `action` varchar(32) NOT NULL COMMENT '动作:upload/bind_model/start_testing/test_pass/test_fail/set_recommended/set_stable/unset_stable/deprecate',
  `from_status` varchar(32) DEFAULT NULL COMMENT '原状态',
  `to_status` varchar(32) DEFAULT NULL COMMENT '目标状态',
  `operator` varchar(100) DEFAULT NULL COMMENT '操作人',
  `operate_at` datetime(3) DEFAULT NULL COMMENT '操作时间',
  `content` text COMMENT '日志内容',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_alpha_firmware_version_logs_firmware_id` (`firmware_id`),
  KEY `idx_alpha_firmware_version_logs_model_id` (`model_id`),
  KEY `idx_alpha_firmware_version_logs_action` (`action`),
  KEY `idx_alpha_firmware_version_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='固件版本日志表';

CREATE TABLE IF NOT EXISTS `alpha_firmware_tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `tag_code` varchar(50) NOT NULL COMMENT '标签编码',
  `tag_name` varchar(100) NOT NULL COMMENT '标签名称',
  `tag_color` varchar(20) DEFAULT NULL COMMENT '标签颜色',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态:1启用 2禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alpha_firmware_tags_tag_code` (`tag_code`),
  KEY `idx_alpha_firmware_tags_deleted_at` (`deleted_at`),
  KEY `idx_alpha_firmware_tags_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='固件标签表';

CREATE TABLE IF NOT EXISTS `alpha_firmware_version_tag_rels` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `firmware_id` bigint unsigned NOT NULL COMMENT '固件版本ID',
  `tag_id` bigint unsigned NOT NULL COMMENT '标签ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alpha_firmware_version_tag_rels_firmware_tag` (`firmware_id`,`tag_id`),
  KEY `idx_alpha_firmware_version_tag_rels_tag_id` (`tag_id`),
  KEY `idx_alpha_firmware_version_tag_rels_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='固件版本标签关系表';

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'major_change', CONVERT(0xE5A4A7E694B9 USING utf8mb4), '#E6A23C', 1, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'major_change'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'fix', CONVERT(0xE4BFAEE5A48D427567 USING utf8mb4), '#F56C6C', 2, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'fix'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'optimize', CONVERT(0xE4BC98E58C96 USING utf8mb4), '#409EFF', 3, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'optimize'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'feature', CONVERT(0xE696B0E5A29EE58A9FE883BD USING utf8mb4), '#67C23A', 4, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'feature'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'compatibility', CONVERT(0xE585BCE5AEB9E680A7E8B083E695B4 USING utf8mb4), '#909399', 5, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'compatibility'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'security', CONVERT(0xE5AE89E585A8E4BFAEE5A48D USING utf8mb4), '#8E44AD', 6, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'security'
);

INSERT INTO `alpha_firmware_tags` (`tag_code`, `tag_name`, `tag_color`, `sort`, `status`, `created_at`, `updated_at`)
SELECT 'performance', CONVERT(0xE680A7E883BDE68F90E58D87 USING utf8mb4), '#00B894', 7, 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `alpha_firmware_tags` WHERE `tag_code` = 'performance'
);
