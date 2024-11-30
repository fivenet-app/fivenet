BEGIN;

ALTER TABLE fivenet_mailer_threads ADD COLUMN `creator_email` varchar(80) NOT NULL;
ALTER TABLE fivenet_mailer_messages ADD COLUMN `creator_email` varchar(80) NOT NULL;

COMMIT;
