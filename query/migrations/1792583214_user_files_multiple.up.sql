BEGIN;

-- A file is the relationship identity. This permits several investigation
-- attachments for one citizen while retaining a single, replaceable mugshot.
ALTER TABLE `fivenet_user_files`
  ADD COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) AFTER `kind`,
  DROP PRIMARY KEY,
  DROP INDEX `idx_file_id`,
  ADD PRIMARY KEY (`file_id`),
  ADD KEY `idx_user_kind` (`user_id`, `kind`);

COMMIT;
