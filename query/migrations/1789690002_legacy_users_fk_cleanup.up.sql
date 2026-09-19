BEGIN;

-- Remove legacy FKs to the ESX `users` table. The application uses
-- `fivenet_user` as its canonical user table.
SET @constraint_name := (
    SELECT kcu.CONSTRAINT_NAME
    FROM information_schema.KEY_COLUMN_USAGE AS kcu
    WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
      AND kcu.TABLE_NAME = 'fivenet_centrum_markers'
      AND kcu.COLUMN_NAME = 'creator_id'
      AND kcu.REFERENCED_TABLE_NAME = 'users'
      AND kcu.CONSTRAINT_NAME = 'fk_fivenet_centrum_markers_user_id'
    LIMIT 1
);

SET @ddl := IF(
    @constraint_name IS NULL,
    'SELECT 1',
    CONCAT(
        'ALTER TABLE `fivenet_centrum_markers` DROP FOREIGN KEY `',
        @constraint_name,
        '`'
    )
);

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @constraint_name := (
    SELECT kcu.CONSTRAINT_NAME
    FROM information_schema.KEY_COLUMN_USAGE AS kcu
    WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
      AND kcu.TABLE_NAME = 'fivenet_user_labels'
      AND kcu.COLUMN_NAME = 'user_id'
      AND kcu.REFERENCED_TABLE_NAME = 'users'
      AND kcu.CONSTRAINT_NAME = 'fk_fivenet_user_citizen_attributes_user_id'
    LIMIT 1
);

SET @ddl := IF(
    @constraint_name IS NULL,
    'SELECT 1',
    CONCAT(
        'ALTER TABLE `fivenet_user_labels` DROP FOREIGN KEY `',
        @constraint_name,
        '`'
    )
);

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
