-- Create "users" table
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Modify "todos" table
ALTER TABLE `todos` ADD COLUMN `user_id` bigint NULL, ADD INDEX `todos_users_todos` (`user_id`), ADD CONSTRAINT `todos_users_todos` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
