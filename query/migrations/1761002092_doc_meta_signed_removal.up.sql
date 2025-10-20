BEGIN;

-- Table: `fivenet_documents_meta` - Remove signature index
set @x := (select 1 from information_schema.statistics where table_name = 'fivenet_documents_meta' and index_name = 'idx_fivenet_documents_meta_signed' and table_schema = database());
set @sql := if( @x is null or @x = 0, 'select ''fivenet_documents_meta idx_fivenet_documents_meta_signed index exists.''', 'ALTER TABLE fivenet_documents_meta DROP INDEX `idx_fivenet_documents_meta_signed`;');
PREPARE stmt FROM @sql;
EXECUTE stmt;

set @x := (select 1 from information_schema.statistics where table_name = 'fivenet_documents_meta' and index_name = 'idx_documents_meta_signed' and table_schema = database());
set @sql := if( @x is null or @x = 0, 'select ''fivenet_documents_meta idx_documents_meta_signed index exists.''', 'ALTER TABLE fivenet_documents_meta DROP INDEX `idx_documents_meta_signed`;');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Remove signature-related columns
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `signed`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_required_remaining`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_declined_count`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_pending_count`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_any_declined`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_required_total`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_collected_valid`;
ALTER TABLE `fivenet_documents_meta` DROP COLUMN `sig_policies_active`;

COMMIT;
