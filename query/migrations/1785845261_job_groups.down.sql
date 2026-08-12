BEGIN;

DROP TABLE IF EXISTS `fivenet_job_groups_visibility_subject`;
DROP TABLE IF EXISTS `fivenet_job_groups_access`;
DROP TABLE IF EXISTS `fivenet_job_group_activity`;
DROP TABLE IF EXISTS `fivenet_job_group_rule_qualification_items`;
DROP TABLE IF EXISTS `fivenet_job_group_rule_qualifications`;
DROP TABLE IF EXISTS `fivenet_job_group_rule_grades`;
DROP TABLE IF EXISTS `fivenet_job_group_rules`;
DROP TABLE IF EXISTS `fivenet_job_group_member_exclusions`;
DROP TABLE IF EXISTS `fivenet_job_group_manual_members`;
DROP TABLE IF EXISTS `fivenet_job_group_leaders`;
DROP TABLE IF EXISTS `fivenet_job_groups`;

COMMIT;
