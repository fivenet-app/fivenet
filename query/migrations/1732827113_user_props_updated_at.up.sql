BEGIN;

ALTER TABLE fivenet_user_props ADD COLUMN `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) AFTER `user_id`;

COMMIT;
