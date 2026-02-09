-- Modify "todos" table
ALTER TABLE `todos` DROP FOREIGN KEY `todos_users_todos`;
-- Modify "todos" table
ALTER TABLE `todos` MODIFY COLUMN `user_id` bigint NOT NULL, ADD CONSTRAINT `todos_users_todos` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE;
