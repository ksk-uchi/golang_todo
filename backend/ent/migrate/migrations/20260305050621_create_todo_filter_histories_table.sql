-- Create "todo_filter_histories" table
CREATE TABLE `todo_filter_histories` (
  `id` char(36) NOT NULL,
  `query` varchar(400) NOT NULL,
  `function_name` varchar(100) NULL,
  `args` json NULL,
  `result_todo_ids` json NULL,
  `created_at` timestamp NOT NULL,
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `todo_filter_histories_users_todo_filter_histories` (`user_id`),
  CONSTRAINT `todo_filter_histories_users_todo_filter_histories` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
