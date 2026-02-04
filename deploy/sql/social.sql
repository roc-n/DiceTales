CREATE TABLE `friend` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL,
   `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL,
   `remark` varchar(255) DEFAULT NULL,
   `add_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `friend_request` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
   `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
   `req_msg` varchar(255) DEFAULT NULL,
   `req_time` timestamp  NOT NULL,
   `req_remark` varchar(255) NOT NULL,
   `handle_result`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
   `handle_msg` varchar(255) DEFAULT NULL,
   `handled_at`timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `groups` (
  `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT '',
  `status`  tinyint COLLATE utf8mb4_unicode_ci NOT NULL,
  `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
  `group_type` int(11) NOT NULL ,
  `is_verify` boolean NOT NULL ,
  `notification` varchar(255) DEFAULT NULL,
  `notification_uid` varchar(64) DEFAULT NULL,
  `description` varchar(512) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group_member` (
 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
 `group_id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL,
 `uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
 `role_level`  tinyint COLLATE utf8mb4_unicode_ci NOT NULL ,
 `join_time` timestamp NULL DEFAULT NULL,
 `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
 `inviter_uid` varchar(64) DEFAULT NULL,
 `operator_uid` varchar(64) DEFAULT NULL,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group_request` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
  `group_id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL,
  `req_msg` varchar(255) DEFAULT NULL,
  `req_time` timestamp NULL DEFAULT NULL,
  `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `inviter_uid` varchar(64) DEFAULT NULL,
  `handle_uid` varchar(64) DEFAULT NULL,
  `handle_time` timestamp NULL DEFAULT NULL,
  `handle_result`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

