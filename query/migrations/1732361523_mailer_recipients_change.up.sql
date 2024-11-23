BEGIN;

ALTER TABLE `fivenet_mailer_threads_recipients` ADD COLUMN `email` varchar(80) NOT NULL;

COMMIT;
