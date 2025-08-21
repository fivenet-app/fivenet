BEGIN;

-- Table: `fivenet_internet_*` remove internet feature tables for now
DROP TABLE IF EXISTS `fivenet_internet_ads`;
DROP TABLE IF EXISTS `fivenet_internet_pages`;
DROP TABLE IF EXISTS `fivenet_internet_domains_access`;
DROP TABLE IF EXISTS `fivenet_internet_domains`;
DROP TABLE IF EXISTS `fivenet_internet_tlds`;

COMMIT;
