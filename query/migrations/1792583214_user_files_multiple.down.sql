BEGIN;

-- The former schema permits one row per user/kind. Retain the oldest
-- relationship for each pair when rolling back from the multi-file model.
DELETE newer
FROM `fivenet_user_files` newer
INNER JOIN `fivenet_user_files` older
  ON older.`user_id` = newer.`user_id`
  AND older.`kind` = newer.`kind`
  AND (
    older.`created_at` < newer.`created_at`
    OR (
      older.`created_at` = newer.`created_at`
      AND older.`file_id` < newer.`file_id`
    )
  );

ALTER TABLE `fivenet_user_files`
  DROP PRIMARY KEY,
  DROP INDEX `idx_user_kind`,
  DROP COLUMN `created_at`,
  ADD PRIMARY KEY (`user_id`, `kind`),
  ADD KEY `idx_file_id` (`file_id`);

COMMIT;
