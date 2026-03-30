CREATE TABLE `meetup` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `organizer_id` bigint(20) unsigned NOT NULL COMMENT '发起用户ID',
   `game_id` bigint(20) unsigned NOT NULL COMMENT '关联桌游ID',
   `title` varchar(128) NOT NULL COMMENT '标题',
   `max_players` int(11) NOT NULL DEFAULT '4' COMMENT '最大人数',
   `current_players` int(11) DEFAULT '1' COMMENT '当前人数(含发起人)',
   `start_time` datetime NOT NULL COMMENT '开始时间',
   `end_time` datetime DEFAULT NULL COMMENT '结束时间',
   `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:报名中 1:满员 2:已结束 3:已取消 4:已过期',
   `enable_audit` tinyint(1) DEFAULT '0' COMMENT '是否开启审核 0:否 1:是',
   `geo_hash` varchar(12) NOT NULL DEFAULT '' COMMENT 'Geohash(冗余用于快速筛选)',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_organizer` (`organizer_id`),
   KEY `idx_start_time` (`start_time`),
   KEY `idx_geo_hash` (`geo_hash`(6))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组局基本信息表';

CREATE TABLE `meetup_location` (
   `meetup_id` bigint(20) unsigned NOT NULL,
   `poi_name` varchar(255) NOT NULL COMMENT 'POI名称(如某桌游店)',
   `address` varchar(512) NOT NULL COMMENT '详细地址',
   `latitude` decimal(10, 6) NOT NULL COMMENT '纬度',
   `longitude` decimal(10, 6) NOT NULL COMMENT '经度',
   `city_code` varchar(64) DEFAULT '' COMMENT '城市编码',
   PRIMARY KEY (`meetup_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组局关联位置表';

CREATE TABLE `meetup_member` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `meetup_id` bigint(20) unsigned NOT NULL,
   `user_id` bigint(20) unsigned NOT NULL,
   `role` tinyint(4) NOT NULL DEFAULT '0' COMMENT '角色 0:成员 1:发起人',
   `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:待审核 1:已加入 2:已拒绝 3:已退出',
   `joined_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_meetup_user` (`meetup_id`, `user_id`),
   KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组局成员状态表';
