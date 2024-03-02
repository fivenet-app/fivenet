BEGIN;

-- Table: fivenet_config
CREATE TABLE
    IF NOT EXISTS `fivenet_config` (
        `key` bigint(20) NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `app_config` longtext,
        `plugin_config` longtext,
        PRIMARY KEY (`key`)
    ) ENGINE = InnoDB;

COMMIT;
