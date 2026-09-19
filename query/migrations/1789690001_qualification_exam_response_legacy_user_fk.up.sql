BEGIN;

-- Remove the legacy FK to the ESX `users` table if it is still present.
-- The application uses `fivenet_user`; some databases may already have the
-- corrected schema, so this must be conditional.
SET @constraint_name := (
    SELECT kcu.CONSTRAINT_NAME
    FROM information_schema.KEY_COLUMN_USAGE AS kcu
    WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
      AND kcu.TABLE_NAME = 'fivenet_qualifications_exam_responses'
      AND kcu.COLUMN_NAME = 'user_id'
      AND kcu.REFERENCED_TABLE_NAME = 'users'
      AND kcu.CONSTRAINT_NAME = 'fivenet_qualifications_exam_responses_user_id'
    LIMIT 1
);

SET @ddl := IF(
    @constraint_name IS NULL,
    'SELECT 1',
    CONCAT(
        'ALTER TABLE `fivenet_qualifications_exam_responses` DROP FOREIGN KEY `',
        @constraint_name,
        '`'
    )
);

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
