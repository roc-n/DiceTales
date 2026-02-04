CREATE TABLE `user` (
         `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL,
         `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
         `nickname` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
         `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
         `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
         `status` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
         `bio` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
         `sex` tinyint COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 0,
         `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
         `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
         `deleted_at` timestamp NULL DEFAULT NULL,
         `city` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
         PRIMARY KEY (`id`),
         UNIQUE KEY `uniq_phone` (`phone`),
         INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;