BEGIN;

ALTER TABLE `fivenet_notifications`
  DROP FOREIGN KEY `fk_fivenet_notifications_actor_user_id`,
  DROP INDEX `idx_fivenet_notifications_user_read_created`,
  DROP INDEX `idx_fivenet_notifications_user_archived_created`,
  DROP INDEX `idx_fivenet_notifications_user_starred_created`,
  DROP INDEX `idx_fivenet_notifications_actor_user_id`,
  DROP INDEX `idx_fivenet_notifications_entity`,
  DROP COLUMN `kind`,
  DROP COLUMN `actor_user_id`,
  DROP COLUMN `entity_type`,
  DROP COLUMN `entity_id`,
  DROP COLUMN `archived_at`;

COMMIT;
