BEGIN;

ALTER TABLE `fivenet_centrum_dispatches_asgmts`
    DROP INDEX `idx_fivenet_centrum_dispatches_asgmts_dispatch_id`;

ALTER TABLE `fivenet_centrum_job_access`
    DROP INDEX `idx_fivenet_centrum_job_access_unit_id`;

ALTER TABLE `fivenet_centrum_markers`
    DROP INDEX `fk_fivenet_centrum_markers_user_id`;

ALTER TABLE `fivenet_centrum_units_users`
    DROP INDEX `idx_fivenet_centrum_units_users_unit_id`;

ALTER TABLE `fivenet_centrum_user_locations`
    DROP INDEX `idx_job`;

ALTER TABLE `fivenet_documents_pins`
    DROP INDEX `idx_fivenet_documents_pins_document_id`;

ALTER TABLE `fivenet_documents_references`
    DROP INDEX `idx_fivenet_documents_references_source_document_id`;

ALTER TABLE `fivenet_documents_relations`
    DROP INDEX `idx_fivenet_documents_relations_document_id`;

ALTER TABLE `fivenet_documents_requests`
    DROP INDEX `idx_fivenet_documents_requests_document_id`,
    DROP INDEX `idx_fivenet_documents_requests_completed`;

ALTER TABLE `fivenet_job_conduct`
    DROP INDEX `type`,
    DROP INDEX `target_user_id`,
    DROP INDEX `idx_conduct_created_at`;

ALTER TABLE `fivenet_mailer_settings_blocked`
    DROP INDEX `idx_fivenet_mailer_settings_blocked`,
    DROP INDEX `idx_fivenet_mailer_settings_blocked_email_id`;

ALTER TABLE `fivenet_mailer_threads_recipients`
    DROP INDEX `idx_fivenet_mailer_threads_recipients_thread_id`;

ALTER TABLE `fivenet_owned_vehicles`
    DROP INDEX `idx_fivenet_owned_vehicles_user`;

ALTER TABLE `fivenet_qualifications`
    DROP INDEX `idx_fivenet_qualifications_job`;

ALTER TABLE `fivenet_qualifications_exam_responses`
    DROP INDEX `fk_fivenet_qualifications_exam_responses_user_id`;

ALTER TABLE `fivenet_qualifications_requirements`
    DROP INDEX `idx_fivenet_qualifications_requirements_qualification_id`,
    DROP INDEX `idx_fivenet_qualifications_requirements_qualification_ids`,
    DROP INDEX `fk_fivenet_qualifications_requirements_target_qualification_id`;

ALTER TABLE `fivenet_rbac_permissions`
    DROP INDEX `idx_category`,
    DROP INDEX `idx_fivenet_permissions_guard_name`;

COMMIT;
