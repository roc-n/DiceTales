CREATE TABLE `game` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `name` varchar(128) NOT NULL COMMENT '中文名',
   `name_en` varchar(255) DEFAULT '' COMMENT '英文名',
   `cover_img` varchar(512) DEFAULT '' COMMENT '封面图URL',
   `publish_year` int(11) DEFAULT NULL COMMENT '发行年份',
   `score` decimal(3,1) DEFAULT '0.0' COMMENT '评分',
   `score_count` int(11) DEFAULT '0' COMMENT '打分人数',
   `min_players` int(11) DEFAULT '1' COMMENT '最少人数',
   `max_players` int(11) DEFAULT '1' COMMENT '最多人数',
   `duration_min` int(11) DEFAULT '0' COMMENT '最短时长(分钟)',
   `duration_max` int(11) DEFAULT '0' COMMENT '最长时长(分钟)',
   `complexity` decimal(3,2) DEFAULT '0.00' COMMENT '重度(1.0-5.0)',
   `description` text COMMENT '简介',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_name` (`name`),
   KEY `idx_score` (`score`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌游核心信息表';

CREATE TABLE `game_resource` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `game_id` bigint(20) unsigned NOT NULL COMMENT '桌游ID',
   `type` tinyint(4) NOT NULL COMMENT '资源类型 1:规则书 2:教学视频 3:扩展包',
   `title` varchar(128) NOT NULL COMMENT '资源标题',
   `url` varchar(512) NOT NULL COMMENT '资源链接',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_game_id` (`game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌游资源表';

CREATE TABLE `tag` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `name` varchar(64) NOT NULL COMMENT '标签名',
   `category` tinyint(4) NOT NULL COMMENT '类型 1:机制 2:主题 3:分类',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_name_cat` (`name`, `category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌游标签元数据表';

CREATE TABLE `game_tag_relation` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `game_id` bigint(20) unsigned NOT NULL,
   `tag_id` bigint(20) unsigned NOT NULL,
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_game_tag` (`game_id`, `tag_id`),
   KEY `idx_tag_game` (`tag_id`, `game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌游标签关联表';
