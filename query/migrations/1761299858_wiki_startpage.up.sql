BEGIN;

ALTER TABLE `fivenet_wiki_pages` ADD `startpage` tinyint(1) DEFAULT 0 NOT NULL;
ALTER TABLE `fivenet_wiki_pages` CHANGE `startpage` `startpage` tinyint(1) DEFAULT 0 NOT NULL AFTER `public`;
ALTER TABLE `fivenet_wiki_pages` ADD INDEX `idx_startpage` (`startpage`);

ALTER TABLE `fivenet_documents_approval_policies` DROP COLUMN `due_at`;

COMMIT;
