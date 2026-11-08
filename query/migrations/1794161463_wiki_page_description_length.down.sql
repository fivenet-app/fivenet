BEGIN;

ALTER TABLE `fivenet_wiki_pages`
  MODIFY COLUMN `description` varchar(128) NOT NULL;

COMMIT;
