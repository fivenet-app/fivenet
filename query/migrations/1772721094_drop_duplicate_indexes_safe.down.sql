BEGIN;

ALTER TABLE `fivenet_centrum_dispatches_asgmts`
    ADD KEY `idx_fivenet_centrum_dispatches_asgmts_dispatch_id` (`dispatch_id`);

ALTER TABLE `fivenet_centrum_job_access`
    ADD KEY `idx_fivenet_centrum_job_access_unit_id` (`source_job`);

ALTER TABLE `fivenet_centrum_markers`
    ADD KEY `fk_fivenet_centrum_markers_user_id` (`creator_id`);

ALTER TABLE `fivenet_centrum_units_users`
    ADD KEY `idx_fivenet_centrum_units_users_unit_id` (`unit_id`);

ALTER TABLE `fivenet_centrum_user_locations`
    ADD KEY `idx_job` (`job`);

ALTER TABLE `fivenet_documents_pins`
    ADD KEY `idx_fivenet_documents_pins_document_id` (`document_id`);

ALTER TABLE `fivenet_documents_references`
    ADD KEY `idx_fivenet_documents_references_source_document_id` (`source_document_id`);

ALTER TABLE `fivenet_documents_relations`
    ADD KEY `idx_fivenet_documents_relations_document_id` (`document_id`);

ALTER TABLE `fivenet_documents_requests`
    ADD KEY `idx_fivenet_documents_requests_document_id` (`document_id`),
    ADD KEY `idx_fivenet_documents_requests_completed` (`accepted`);

ALTER TABLE `fivenet_job_conduct`
    ADD KEY `type` (`type`),
    ADD KEY `target_user_id` (`target_user_id`),
    ADD KEY `idx_conduct_created_at` (`created_at`);

ALTER TABLE `fivenet_mailer_settings_blocked`
    ADD UNIQUE KEY `idx_fivenet_mailer_settings_blocked` (`email_id`, `target_email`),
    ADD KEY `idx_fivenet_mailer_settings_blocked_email_id` (`email_id`);

ALTER TABLE `fivenet_mailer_threads_recipients`
    ADD KEY `idx_fivenet_mailer_threads_recipients_thread_id` (`thread_id`);

ALTER TABLE `fivenet_owned_vehicles`
    ADD KEY `idx_fivenet_owned_vehicles_user` (`user_id`);

ALTER TABLE `fivenet_qualifications`
    ADD KEY `idx_fivenet_qualifications_job` (`job`);

ALTER TABLE `fivenet_qualifications_exam_responses`
    ADD KEY `fk_fivenet_qualifications_exam_responses_user_id` (`user_id`);

ALTER TABLE `fivenet_qualifications_requirements`
    ADD KEY `idx_fivenet_qualifications_requirements_qualification_id` (`qualification_id`),
    ADD KEY `idx_fivenet_qualifications_requirements_qualification_ids` (`qualification_id`, `target_qualification_id`),
    ADD KEY `fk_fivenet_qualifications_requirements_target_qualification_id` (`target_qualification_id`);

ALTER TABLE `fivenet_rbac_permissions`
    ADD KEY `idx_category` (`category`),
    ADD UNIQUE KEY `idx_fivenet_permissions_guard_name` (`guard_name`);

COMMIT;
