BEGIN;

-- Table: `fivenet_jobs_conduct` - Add deleted_at column
ALTER TABLE `fivenet_jobs_conduct`
  ADD COLUMN `deleted_at` datetime(3) AFTER `updated_at`;

ALTER TABLE `fivenet_jobs_conduct`
  ADD KEY `idx_fivenet_jobs_conduct_deleted_at` (`deleted_at`);

-- Table: `fivenet_jobs_labels` - Add deleted_at column
ALTER TABLE `fivenet_jobs_labels`
  ADD COLUMN `deleted_at` datetime(3) AFTER `job`;

ALTER TABLE `fivenet_jobs_labels`
  ADD KEY `idx_fivenet_jobs_labels_deleted_at` (`deleted_at`);

-- Table: `fivenet_jobs_user_props` - Add deleted_at column
ALTER TABLE `fivenet_jobs_user_props`
  ADD COLUMN `deleted_at` datetime(3) AFTER `job`;

ALTER TABLE `fivenet_jobs_user_props`
  ADD KEY `idx_fivenet_jobs_user_props_deleted_at` (`deleted_at`);

COMMIT;
