BEGIN;

-- Table: `fivenet_centrum_settings` - Add deleted_at column
ALTER TABLE `fivenet_centrum_settings`
  ADD COLUMN `deleted_at` datetime(3) AFTER `job`;

ALTER TABLE `fivenet_centrum_settings`
  ADD KEY `idx_fivenet_centrum_settings_deleted_at` (`deleted_at`);

COMMIT;
