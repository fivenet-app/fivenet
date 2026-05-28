BEGIN;

ALTER TABLE `fivenet_centrum_dispatchers` ADD COLUMN `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) AFTER `user_id`;

ALTER TABLE `fivenet_centrum_units_users` ADD COLUMN `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) AFTER `user_id`;

COMMIT;
