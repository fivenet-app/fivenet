BEGIN;

-- Table: `fivenet_documents_categories` - Fix document category color (primary doesn't exist; blue is default)
ALTER TABLE `fivenet_documents_categories`
  MODIFY COLUMN `color` char(7) DEFAULT 'blue' NULL;
UPDATE `fivenet_documents_categories` SET color = 'blue' WHERE color = 'primary';

-- Table: `fivenet_documents_categories` - Add created_at and deleted_at columns
ALTER TABLE `fivenet_documents_categories`
  ADD COLUMN `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) AFTER `id`,
  ADD COLUMN `deleted_at` datetime(3) AFTER `created_at`;

ALTER TABLE `fivenet_documents_categories`
  ADD KEY `idx_fivenet_documents_categories_deleted_at` (`deleted_at`);

-- Table: `fivenet_calendar` - Fix color
ALTER TABLE `fivenet_calendar`
  MODIFY COLUMN `color` varchar(24) DEFAULT 'blue' NULL;

UPDATE `fivenet_calendar` SET color = 'blue' WHERE color = 'primary';

COMMIT;
