-- 1. 移除新增的 id 字段 (由于是主键，直接删除该列即可移除主键)
ALTER TABLE `user` DROP COLUMN `id`;

-- 2. 移除 account 列上的唯一索引
ALTER TABLE `user` DROP INDEX `uk_account`;

-- 3. 将 account 字段重命名回 id
ALTER TABLE `user` CHANGE `account` `id` VARCHAR(24) NOT NULL COMMENT '用户账号';

-- 4. 重新将 id 设为主键
ALTER TABLE `user` ADD PRIMARY KEY (`id`);
