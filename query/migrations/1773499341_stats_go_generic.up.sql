BEGIN;

RENAME TABLE `fivenet_documents_stats_daily_rollup` TO `fivenet_stats_daily_rollup`;

COMMIT;
