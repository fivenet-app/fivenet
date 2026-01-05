BEGIN;

ALTER TABLE `fivenet_documents` MODIFY `summary` VARCHAR(256) NOT NULL;

ALTER TABLE `fivenet_documents` DROP INDEX `idx_fivenet_documents_content`;

ALTER TABLE `fivenet_documents` CHANGE `content` `content_json` LONGTEXT JSON NOT NULL;

ALTER TABLE `fivenet_documents` ADD COLUMN `content_text` LONGTEXT NOT NULL AFTER `content_json`;

ALTER TABLE `fivenet_documents`
  ADD COLUMN `word_count` INT UNSIGNED NOT NULL DEFAULT 0 AFTER `summary`,
  ADD COLUMN `first_heading` VARCHAR(256) NOT NULL DEFAULT '' AFTER `word_count`;

ALTER TABLE `fivenet_documents` ADD FULLTEXT KEY `idx_fivenet_documents_content_text` (`content_text`);

ALTER TABLE `fivenet_documents_comments` CHANGE `comment` `content` LONGTEXT NULL;

COMMIT;
