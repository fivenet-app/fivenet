BEGIN;

DELETE FROM `fivenet_documents_stats_metric` WHERE `metric_key` = 'case_count';

COMMIT;
