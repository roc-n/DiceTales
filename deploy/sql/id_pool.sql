CREATE TABLE `id_pool` (
  `id` bigint unsigned NOT NULL COMMENT '号码',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0:未分配 1:已分配/使用中 2:保留/冻结',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ID号码池';

CREATE TABLE `group_id_pool` (
  `id` bigint unsigned NOT NULL COMMENT '号码',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0:未分配 1:已分配',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群主ID号码池';