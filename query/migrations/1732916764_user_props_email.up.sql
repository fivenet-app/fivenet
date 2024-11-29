BEGIN;

ALTER TABLE fivenet_user_props ADD COLUMN `email` varchar(80) DEFAULT NULL;

COMMIT;
