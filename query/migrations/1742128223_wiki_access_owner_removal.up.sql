BEGIN;

UPDATE `fivenet_wiki_pages_access` SET `access` = 4 WHERE `access` = 5;

COMMIT;
