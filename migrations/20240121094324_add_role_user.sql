-- migrate:up
ALTER TABLE `users` ADD COLUMN `role` ENUM('admin', 'user') DEFAULT 'user';

-- migrate:down
ALTER TABLE `users` DROP COLUMN `role`;
