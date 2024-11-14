BEGIN;

DROP TABLE `fivenet_msgs_threads_job_access`;
ALTER TABLE `fivenet_msgs_threads` DROP COLUMN `archived`;
ALTER TABLE `fivenet_msgs_threads_user_state` ADD `archived` tinyint(1) DEFAULT '0';

DELETE FROM `fivenet_permissions` WHERE `category` = 'MessengerService';

COMMIT;
