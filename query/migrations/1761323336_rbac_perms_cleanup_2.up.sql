BEGIN;

DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'wiki.WikiService' AND `name` = 'UploadFile' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'wiki.WikiService' AND `name` = 'GetPage' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'UploadJobLogo' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'ListUserGuilds' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'ListDiscordChannels' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'GetRole' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'GetPermissions' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'GetEffectivePermissions' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.SettingsService' AND `name` = 'DeleteJobLogo' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.LawsService' AND `name` = 'DeleteLaw' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'settings.LawsService' AND `name` = 'CreateOrUpdateLaw' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'UploadFile' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'TakeExam' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'SubmitExam' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'ListQualificationsResults' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'ListQualificationRequests' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'GetUserExam' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'GetQualification' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'GetExamInfo' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'DeleteQualificationResult' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'DeleteQualificationReq' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'CreateQualification' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'CreateOrUpdateQualificationResult' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'qualifications.QualificationsService' AND `name` = 'CreateOrUpdateQualificationRequest' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'SetThreadState' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'SetEmailSettings' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'SearchThreads' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'PostMessage' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'ListThreads' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'ListThreadMessages' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'ListTemplates' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetThreadState' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetThread' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetTemplate' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetEmailSettings' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetEmailProposals' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'GetEmail' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'DeleteTemplate' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'CreateThread' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'mailer.MailerService' AND `name` = 'CreateOrUpdateTemplate' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'jobs.TimeclockService' AND `name` = 'GetTimeclockStats' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'jobs.JobsService' AND `name` = 'GetSelf' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'jobs.JobsService' AND `name` = 'GetColleagueLabelsStats' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'jobs.JobsService' AND `name` = 'GetColleagueLabels' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'UploadFile' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'UpdateTemplate' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'UpdateDocumentReq' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'SetDocumentAccess' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'RemoveDocumentRelation' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'RemoveDocumentReference' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'PostComment' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'ListDocumentPins' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetTemplate' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetDocumentRelations' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetDocumentReferences' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetDocumentAccess' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetDocument' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'GetComments' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'EditComment' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.DocumentsService' AND `name` = 'CreateDocument' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.ApprovalService' AND `name` = 'ReopenApprovalTask' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.ApprovalService' AND `name` = 'RecomputeApprovalPolicyCounters' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'citizens.CitizensService' AND `name` = 'UploadMugshot' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'citizens.CitizensService' AND `name` = 'DeleteMugshot' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'UpdateUnitStatus' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'UpdateDispatchStatus' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'ListUnits' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'ListUnitActivity' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'ListDispatches' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'ListDispatchTargetJobs' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'ListDispatchActivity' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'JoinUnit' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'GetSettings' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'GetDispatchHeatmap' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'GetDispatch' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'AssignUnit' LIMIT 1;
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'centrum.CentrumService' AND `name` = 'AssignDispatch' LIMIT 1;

-- Ensure wiki CreatePage permission exists
INSERT IGNORE INTO fivenet_rbac_permissions
(`category`, `name`, `guard_name`, `order`, `icon`)
SELECT 'wiki.WikiService', 'CreatePage', 'wiki-wikiservice-createpage', 11000, 'i-mdi-brain'
WHERE EXISTS (
  SELECT NULL FROM fivenet_rbac_permissions
);

-- Assign CreatePage permission to job perms that currently have UpdatePage permission
INSERT INTO `fivenet_rbac_job_permissions` (job, permission_id, val)
    SELECT
	fr.job,
	(
	SELECT
		id
	FROM
		`fivenet_rbac_permissions`
	WHERE
		`category` = 'wiki.WikiService'
		AND `name` = 'CreatePage'
	LIMIT 1) AS `permission_id`,
	1 AS `val`
FROM
	`fivenet_rbac_job_permissions` fr,
	`fivenet_rbac_permissions` fp
WHERE
	fr.permission_id = (
	SELECT
		id
	FROM
		`fivenet_rbac_permissions`
	WHERE
		`category` = 'wiki.WikiService'
		AND `name` = 'UpdatePage'
	LIMIT 1)
	AND fp.category = 'wiki.WikiService'
	AND fp.name = 'UpdatePage'
ON DUPLICATE KEY UPDATE val = 1;

COMMIT;
