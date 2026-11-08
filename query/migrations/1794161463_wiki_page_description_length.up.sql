BEGIN;

ALTER TABLE `fivenet_wiki_pages`
  MODIFY COLUMN `description` varchar(512) NOT NULL;

COMMIT;
