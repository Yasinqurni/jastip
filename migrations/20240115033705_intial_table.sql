-- migrate:up
CREATE TABLE `users` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `phone_number` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

CREATE TABLE `subscribes` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `user_id` varchar(255) NOT NULL,
  `class_id` varchar(255) NOT NULL,
  `status` ENUM ('expired', 'active', 'waiting'),
  `expire_date` timestamp NOT NULL,
  `active_at` timestamp,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

CREATE TABLE `subscribe_class` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `time` varchar(255) NOT NULL,
  `price` varchar(255) NOT NULL,
  `description` text,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

CREATE TABLE `admins` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `user_id` varchar(255),
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

CREATE TABLE `rutes` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `user_id` varchar(255),
  `from` varchar(255) NOT NULL,
  `to` varchar(255) NOT NULL,
  `price` varchar(255) NOT NULL,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

CREATE TABLE `packages` (
  `id` varchar(255) UNIQUE PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `owner` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  `phone_number` varchar(255) NOT NULL,
  `admin_id` varchar(255) NOT NULL,
  `weight` double NOT NULL,
  `amount` double NOT NULL,
  `rute_id` varchar(255) NOT NULL,
  `status` ENUM ('received', 'sending', 'arrived', 'picked'),
  `notes` text,
  `created_at` timestamp,
  `updated_at` timestamp,
  `deleted_at` timestamp
);

ALTER TABLE `subscribes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `subscribes` ADD FOREIGN KEY (`class_id`) REFERENCES `subscribe_class` (`id`);

ALTER TABLE `packages` ADD FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`);

ALTER TABLE `packages` ADD FOREIGN KEY (`rute_id`) REFERENCES `rutes` (`id`);


-- migrate:down
DROP TABLE `users`;
DROP TABLE `subscribes`;
DROP TABLE `subscribe_class`;
DROP TABLE `admins`;
DROP TABLE `rutes`;
DROP TABLE `packages`;
