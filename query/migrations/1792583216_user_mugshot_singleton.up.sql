BEGIN;

-- Keep the newest mugshot relation before enforcing singleton semantics.
DELETE newer
FROM `fivenet_user_files` newer
INNER JOIN `fivenet_user_files` older
  ON older.`user_id` = newer.`user_id`
  AND older.`kind` = 'mugshot'
  AND newer.`kind` = 'mugshot'
  AND (
    older.`created_at` > newer.`created_at`
    OR (
      older.`created_at` = newer.`created_at`
      AND older.`file_id` > newer.`file_id`
    )
  );

-- Attachments remain multi-valued, while mugshots become one relation per user.
-- Keep this column virtual: MySQL rejects an indexed stored generated column
-- based on a foreign-key column with ON DELETE CASCADE.
ALTER TABLE `fivenet_user_files`
  ADD COLUMN `mugshot_user_id` int
    GENERATED ALWAYS AS (IF(`kind` = 'mugshot', `user_id`, NULL)) VIRTUAL,
  ADD UNIQUE KEY `uq_fivenet_user_files_mugshot_user` (`mugshot_user_id`);

COMMIT;
