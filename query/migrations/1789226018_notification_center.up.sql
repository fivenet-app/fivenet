BEGIN;

ALTER TABLE `fivenet_notifications`
  ADD COLUMN `kind` smallint(5) NOT NULL DEFAULT 0 AFTER `category`,
  ADD COLUMN `actor_user_id` int(11) DEFAULT NULL AFTER `user_id`,
  ADD COLUMN `entity_type` varchar(64) DEFAULT NULL AFTER `actor_user_id`,
  ADD COLUMN `entity_id` bigint(20) DEFAULT NULL AFTER `entity_type`,
  ADD COLUMN `archived_at` datetime(3) DEFAULT NULL AFTER `read_at`,
  ADD KEY `idx_fivenet_notifications_user_read_created` (`user_id`, `read_at`, `archived_at`, `created_at`, `id`),
  ADD KEY `idx_fivenet_notifications_user_archived_created` (`user_id`, `archived_at`, `created_at`, `id`),
  ADD KEY `idx_fivenet_notifications_user_starred_created` (`user_id`, `starred`, `archived_at`, `created_at`, `id`),
  ADD KEY `idx_fivenet_notifications_actor_user_id` (`actor_user_id`),
  ADD KEY `idx_fivenet_notifications_entity` (`entity_type`, `entity_id`),
  ADD CONSTRAINT `fk_fivenet_notifications_actor_user_id` FOREIGN KEY (`actor_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

COMMIT;
