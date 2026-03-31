CREATE TABLE `post` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `user_id` bigint(20) unsigned NOT NULL COMMENT '发布者ID',
   `content` text COMMENT '帖子内容(支持Markdown)',
   `images` json DEFAULT NULL COMMENT '图片URL列表(JSON数组)',
   `visibility` tinyint(4) NOT NULL DEFAULT '0' COMMENT '可见性 0:公开 1:仅好友 2:私密',
   `related_game_id` bigint(20) unsigned DEFAULT '0' COMMENT '关联桌游ID(可选)',
   `like_count` int(11) DEFAULT '0' COMMENT '点赞数(冗余字段)',
   `comment_count` int(11) DEFAULT '0' COMMENT '评论数(冗余字段)',
   `collect_count` int(11) DEFAULT '0' COMMENT '收藏数(冗余字段)',
   `is_deleted` tinyint(1) DEFAULT '0' COMMENT '是否被删除',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_user_time` (`user_id`, `created_at`),
   KEY `idx_game_time` (`related_game_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区动态表';

CREATE TABLE `comment` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `post_id` bigint(20) unsigned NOT NULL COMMENT '关联帖子ID',
   `user_id` bigint(20) unsigned NOT NULL COMMENT '评论者ID',
   `content` varchar(1024) NOT NULL COMMENT '评论内容',
   `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '父评论ID(一级评论为0)',
   `root_id` bigint(20) unsigned DEFAULT '0' COMMENT '根评论ID(二级及以下评论所属的一级评论ID)',
   `like_count` int(11) DEFAULT '0' COMMENT '点赞数',
   `status` tinyint(4) DEFAULT '0' COMMENT '状态 0:正常 1:隐藏',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_post_root` (`post_id`, `root_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

CREATE TABLE `social_interaction` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
   `user_id` bigint(20) unsigned NOT NULL COMMENT '操作用户ID',
   `target_id` bigint(20) unsigned NOT NULL COMMENT '目标ID(Post/Comment)',
   `target_type` tinyint(4) NOT NULL COMMENT '目标类型 1:Post 2:Comment',
   `action` tinyint(4) NOT NULL COMMENT '动作类型 1:Like 2:Collect',
   `status` tinyint(4) DEFAULT '0' COMMENT '状态 0:有效 1:取消',
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_user_target_action` (`user_id`, `target_id`, `target_type`, `action`),
   KEY `idx_target_action` (`target_id`, `target_type`, `action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社交互动记录表(点赞/收藏)';
