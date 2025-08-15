BEGIN;

-- Table: Ensure that the `fivenet_user` table `updated_at` column (if it exists) is renamed to `last_seen`
set @x := (select 1 from information_schema.columns where table_name = 'fivenet_user' and column_name = 'updated_at' and table_schema = database());
set @sql := if( @x is not null and @x > 0, 'ALTER TABLE fivenet_user CHANGE updated_at last_seen timestamp on update CURRENT_TIMESTAMP NULL', 'select ''fivenet_user updated_at column exists.''');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: fivenet_accounts_oauth2 change `expires_in` column type from int32 to int64 (go oauth2 uses int64 for expires_in field)
ALTER TABLE `fivenet_accounts_oauth2` MODIFY COLUMN `expires_in` bigint NULL;

COMMIT;
