BEGIN;

ALTER TABLE `fivenet_centrum_markers`
    ADD COLUMN `public` tinyint(1) NOT NULL DEFAULT 0 AFTER `job`;

COMMIT;
