CREATE TABLE `game` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '中文名',
   `name_en` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '英文名',
   `cover_img` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '封面图URL',
   `score` decimal(3,1) DEFAULT '0.0' COMMENT '评分',
   `score_count` int(11) DEFAULT '0' COMMENT '打分人数',
   `min_players` int(11) DEFAULT '1' COMMENT '最少人数',
   `max_players` int(11) DEFAULT '1' COMMENT '最多人数',
   `min_recommended_players` int(11) DEFAULT NULL COMMENT '最少推荐人数',
   `max_recommended_players` int(11) DEFAULT NULL COMMENT '最多推荐人数',
   `need_host` tinyint(1) DEFAULT '0' COMMENT '是否需要主持人 0-否 1-是',
   `rank` int(11) DEFAULT NULL COMMENT '排行榜名次',
   `year` int(4) DEFAULT NULL COMMENT '发行年份',
   `description` text COLLATE utf8mb4_unicode_ci COMMENT '桌游简介',
   `difficulty` decimal(3,1) DEFAULT NULL COMMENT '上手难度(满分10级)',
   `duration_per_player` int(11) DEFAULT NULL COMMENT '人均时长(分钟)',
   `setup_time` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设置时长',
   `language_dependency` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '语言需求文本描述，例如：较高',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='桌游信息表';

CREATE TABLE `game_category_info` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `game_id` int(11) unsigned NOT NULL COMMENT '桌游ID，关联game表的id',
   `category` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏类别',
   `mode` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏模式',
   `theme` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏主题',
   `mechanic` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏机制',
   `portability` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '便携程度',
   `table_requirement` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '桌面要求',
   `suitable_age` int(11) DEFAULT NULL COMMENT '适合年龄',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_game_id` (`game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='桌游分类信息表';