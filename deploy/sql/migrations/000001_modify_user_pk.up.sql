-- 1. 取消原来的主键限制
ALTER TABLE `user` DROP PRIMARY KEY;

-- 2. 将原来的 id 字段重命名为 account并保留原类型
ALTER TABLE `user` CHANGE `id` `account` VARCHAR(24) NOT NULL COMMENT '用户账号';

-- 3. 加上唯一索引
ALTER TABLE `user` ADD UNIQUE INDEX `uk_account` (`account`);

-- 4. 新增一个 id 字段，类型为 BIGINT AUTO_INCREMENT，并设为主键 (将其置为第一列)
ALTER TABLE `user` ADD COLUMN `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键' FIRST;
