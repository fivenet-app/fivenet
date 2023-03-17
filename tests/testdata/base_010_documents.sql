-- Table data: apranet_documents_categories
INSERT INTO arpanet_documents_categories
(id, name, description, job)
VALUES(1, 'Patient / File', 'Patient files (e.g., reports, results)', 'ambulance');
INSERT INTO arpanet_documents_categories
(id, name, description, job)
VALUES(2, 'Criminal Record', 'Criminal record of a citizen', 'police');
INSERT INTO arpanet_documents_categories
(id, name, description, job)
VALUES(3, 'Non-Existant', 'Document Category for a non-existant job, no person should see it.', 'non-existant');
INSERT INTO arpanet_documents_categories
(id, name, description, job)
VALUES(4, 'Patient / Unused', 'Unused category for testing 🙃', 'ambulance');

-- Table data: arpanet_documents
INSERT INTO arpanet_documents
(id, created_at, updated_at, deleted_at, category_id, title, content_type, content, `data`, creator_id, state, closed, public)
VALUES(1, NULL, NULL, NULL, NULL, 'Public Document without category', 0, 'I''m a public Document without a category.', NULL, 1, 'Open', NULL, 1);
INSERT INTO arpanet_documents
(id, created_at, updated_at, deleted_at, category_id, title, content_type, content, `data`, creator_id, state, closed, public)
VALUES(2, NULL, NULL, NULL, NULL, 'Public Document with category (Closed State)', 0, 'I''m a public Document with a category that is closed.', NULL, 1, 'Closed', 1, 1);
