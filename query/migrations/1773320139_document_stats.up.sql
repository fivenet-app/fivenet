BEGIN;

CREATE TABLE `fivenet_documents_stats_metric` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,

  `document_id` bigint(20) unsigned NOT NULL,

  `job` varchar(40) NOT NULL,
  `source_key` varchar(64) NOT NULL, -- e.g. penalty_calculator
  `metric_key` varchar(64) NOT NULL, -- e.g. case_count, fine_total, law_count

  `dimension1` varchar(128) NULL, -- e.g. law_id
  `dimension2` varchar(128) NULL,
  `dimension3` varchar(128) NULL,

  `value` bigint(20) unsigned NOT NULL,

  `occurred_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,

  UNIQUE KEY `uq_metric` (
    `document_id`,
    `source_key`,
    `metric_key`,
    `dimension1`,
    `dimension2`,
    `dimension3`
  ),

  PRIMARY KEY (`id`),
  KEY `ix_doc` (`document_id`),
  KEY `ix_job_source_time` (`job`, `source_key`, `occurred_at`),
  KEY `ix_metric_time` (`source_key`, `metric_key`, `occurred_at`),
  KEY `ix_dim1` (`dimension1`),
  CONSTRAINT `fk_fivenet_documents_stats_daily_rollup_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `fivenet_documents_stats_daily_rollup` (
  `day` date NOT NULL,
  `job` varchar(40) NOT NULL,

  `source_kind` varchar(32) NOT NULL,
  `source_key` varchar(64) NOT NULL,
  `metric_key` varchar(64) NOT NULL,

  `dimension1` varchar(128) NOT NULL DEFAULT '',
  `dimension2` varchar(128) NOT NULL DEFAULT '',
  `dimension3` varchar(128) NOT NULL DEFAULT '',

  `value` bigint(20) unsigned NOT NULL,

  PRIMARY KEY (`day`, `job`, `source_kind`, `source_key`, `metric_key`, `dimension1`, `dimension2`, `dimension3`),
  KEY `ix_rollup_source_day` (`source_kind`, `source_key`, `metric_key`, `day`),
  KEY `ix_rollup_job_day` (`job`, `day`)
) ENGINE=InnoDB;

COMMIT;
