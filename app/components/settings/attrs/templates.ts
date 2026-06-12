import { PermAttributes, type PermAttrDescriptor, type PermAttrKey, type PermAttrPerm, type Perms } from '~~/gen/ts/perms';
import type { AttributeValues } from '~~/gen/ts/resources/permissions/attributes/attributes';

export type Template = {
    permissions: Perms[];
    attributes: TemplateAttribute[];
    grantAllPermissions?: boolean;
};

export type TemplateAttribute = {
    permission: Perms;
    key: string;
    validValues?: AttributeValues;
};

const templateAttribute = <P extends PermAttrPerm, K extends PermAttrKey<P>>(
    permission: P,
    key: K,
    validValues?: AttributeValues,
): TemplateAttribute => {
    const defaultValidValues = (): AttributeValues | undefined => {
        const descriptor = PermAttributes[permission][key] as PermAttrDescriptor<P, K>;

        switch (descriptor.type) {
            case 'stringList':
                return {
                    validValues: {
                        oneofKind: 'stringList',
                        stringList: {
                            strings: [...descriptor.values],
                        },
                    },
                };
            case 'jobList':
                return {
                    validValues: {
                        oneofKind: 'jobList',
                        jobList: {
                            strings: [],
                        },
                    },
                };
            case 'jobGradeList':
                return {
                    validValues: {
                        oneofKind: 'jobGradeList',
                        jobGradeList: {
                            fineGrained: false,
                            jobs: {},
                            grades: {},
                        },
                    },
                };
        }
    };

    return {
        permission,
        key,
        validValues: validValues ?? defaultValidValues(),
    };
};

const policePermissions = [
    'auth.AuthService/ChooseCharacter',
    'calendar.CalendarService/CreateCalendar',
    'centrum.CentrumService/Stream',
    'centrum.CentrumService/TakeControl',
    'centrum.CentrumService/UpdateDispatchers',
    'centrum.CentrumService/UpdateSettings',
    'centrum.DispatchesService/CreateDispatch',
    'centrum.DispatchesService/DeleteDispatch',
    'centrum.DispatchesService/TakeDispatch',
    'centrum.DispatchesService/UpdateDispatch',
    'centrum.UnitsService/CreateOrUpdateUnit',
    'centrum.UnitsService/DeleteUnit',
    'citizens.CitizensService/GetUser',
    'citizens.CitizensService/ListCitizens',
    'citizens.CitizensService/ListUserActivity',
    'citizens.CitizensService/SetUserProps',
    'citizens.LabelsService/CreateOrUpdateLabel',
    'citizens.LabelsService/DeleteLabel',
    'documents.ApprovalService/DeleteApprovalTasks',
    'documents.ApprovalService/RevokeApproval',
    'documents.ApprovalService/UpsertApprovalPolicy',
    'documents.ApprovalService/UpsertApprovalTasks',
    'documents.CategoriesService/CreateOrUpdateCategory',
    'documents.CategoriesService/DeleteCategory',
    'documents.CategoriesService/ListCategories',
    'documents.CommentsService/DeleteComment',
    'documents.DocumentsService/AddDocumentReference',
    'documents.DocumentsService/AddDocumentRelation',
    'documents.DocumentsService/ChangeDocumentOwner',
    'documents.DocumentsService/CreateDocumentReq',
    'documents.DocumentsService/DeleteDocument',
    'documents.DocumentsService/DeleteDocumentReq',
    'documents.DocumentsService/ListDocumentActivity',
    'documents.DocumentsService/ListDocumentReqs',
    'documents.DocumentsService/ListDocuments',
    'documents.DocumentsService/ListUserDocuments',
    'documents.DocumentsService/SetDocumentReminder',
    'documents.DocumentsService/ToggleDocument',
    'documents.DocumentsService/ToggleDocumentPin',
    'documents.DocumentsService/UpdateDocument',
    'documents.StampsService/DeleteStamp',
    'documents.StampsService/ListUsableStamps',
    'documents.StampsService/UpsertStamp',
    'documents.StatsService/GetStats',
    'documents.TemplatesService/CreateTemplate',
    'documents.TemplatesService/DeleteTemplate',
    'documents.TemplatesService/ListTemplates',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/ListColleagueActivity',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.ColleaguesService/ManageLabels',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.JobsService/SetMOTD',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.TimeclockService/ListTimeclock',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'livemap.LivemapService/Stream',
    'mailer.MailerService/CreateOrUpdateEmail',
    'mailer.MailerService/DeleteEmail',
    'mailer.MailerService/ListEmails',
    'qualifications.QualificationsService/DeleteQualification',
    'qualifications.QualificationsService/ListQualifications',
    'qualifications.QualificationsService/UpdateQualification',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'settings.SettingsService/ViewAuditLog',
    'vehicles.VehiclesService/ListVehicles',
    'vehicles.VehiclesService/SetVehicleProps',
    'wiki.WikiService/CreatePage',
    'wiki.WikiService/DeletePage',
    'wiki.WikiService/ListPageActivity',
    'wiki.WikiService/ListPages',
    'wiki.WikiService/MovePage',
    'wiki.WikiService/UpdatePage',
] as const satisfies readonly Perms[];

const neutralPermissions = [
    'auth.AuthService/ChooseCharacter',
    'calendar.CalendarService/CreateCalendar',
    'documents.CategoriesService/ListCategories',
    'documents.DocumentsService/ListDocuments',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/ListColleagueActivity',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.ColleaguesService/ManageLabels',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.JobsService/SetMOTD',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.TimeclockService/ListTimeclock',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'livemap.LivemapService/Stream',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'wiki.WikiService/CreatePage',
    'wiki.WikiService/DeletePage',
    'wiki.WikiService/ListPageActivity',
    'wiki.WikiService/ListPages',
    'wiki.WikiService/MovePage',
    'wiki.WikiService/UpdatePage',
] as const satisfies readonly Perms[];

const gangPermissions = [
    'auth.AuthService/ChooseCharacter',
    'calendar.CalendarService/CreateCalendar',
    'documents.DocumentsService/ListDocuments',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/ListColleagueActivity',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.ColleaguesService/ManageLabels',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.JobsService/SetMOTD',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.TimeclockService/ListTimeclock',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'livemap.LivemapService/Stream',
    'mailer.MailerService/CreateOrUpdateEmail',
    'mailer.MailerService/DeleteEmail',
    'mailer.MailerService/ListEmails',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'wiki.WikiService/DeletePage',
    'wiki.WikiService/ListPageActivity',
    'wiki.WikiService/ListPages',
    'wiki.WikiService/MovePage',
    'wiki.WikiService/UpdatePage',
] as const satisfies readonly Perms[];

const policeAttributes = [
    templateAttribute('citizens.CitizensService/GetUser', 'Jobs'),
    templateAttribute('citizens.CitizensService/ListCitizens', 'Fields'),
    templateAttribute('citizens.CitizensService/ListUserActivity', 'Fields'),
    templateAttribute('citizens.CitizensService/SetUserProps', 'Fields'),
    templateAttribute('livemap.LivemapService/Stream', 'Players'),
    templateAttribute('documents.DocumentsService/DeleteDocument', 'Access'),
    templateAttribute('documents.CommentsService/DeleteComment', 'Access'),
    templateAttribute('documents.CategoriesService/ListCategories', 'Jobs'),
    templateAttribute('documents.DocumentsService/UpdateDocument', 'Access'),
    templateAttribute('documents.DocumentsService/ToggleDocument', 'Access'),
    templateAttribute('livemap.LivemapService/CreateOrUpdateMarker', 'Access'),
    templateAttribute('livemap.LivemapService/DeleteMarker', 'Access'),
    templateAttribute('documents.DocumentsService/ChangeDocumentOwner', 'Access'),
    templateAttribute('documents.DocumentsService/CreateDocumentReq', 'Types'),
    templateAttribute('documents.StatsService/GetStats', 'Categories', {
        validValues: { oneofKind: 'stringList', stringList: { strings: ['PenaltyCalculator'] } },
    }),
    templateAttribute('documents.StatsService/GetStats', 'Jobs'),
    templateAttribute('jobs.ConductService/ListConductEntries', 'Access'),
    templateAttribute('jobs.TimeclockService/ListTimeclock', 'Access'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Access'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Access'),
    templateAttribute('jobs.ColleaguesService/ListColleagueActivity', 'Types'),
    templateAttribute('qualifications.QualificationsService/DeleteQualification', 'Access'),
    templateAttribute('qualifications.QualificationsService/UpdateQualification', 'Access'),
    templateAttribute('calendar.CalendarService/CreateCalendar', 'Fields'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Types'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Types'),
    templateAttribute('wiki.WikiService/UpdatePage', 'Fields', {
        validValues: { oneofKind: 'stringList', stringList: { strings: [] } },
    }),
    templateAttribute('mailer.MailerService/CreateOrUpdateEmail', 'Fields'),
    templateAttribute('documents.DocumentsService/ToggleDocumentPin', 'Types'),
    templateAttribute('qualifications.QualificationsService/UpdateQualification', 'Fields'),
    templateAttribute('centrum.CentrumService/UpdateSettings', 'Access'),
    templateAttribute('vehicles.VehiclesService/ListVehicles', 'Fields'),
    templateAttribute('vehicles.VehiclesService/SetVehicleProps', 'Fields'),
] as const satisfies readonly TemplateAttribute[];

const neutralAttributes = [
    templateAttribute('documents.CategoriesService/ListCategories', 'Jobs'),
    templateAttribute('livemap.LivemapService/Stream', 'Markers'),
    templateAttribute('jobs.ConductService/ListConductEntries', 'Access'),
    templateAttribute('jobs.TimeclockService/ListTimeclock', 'Access'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Access'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Access'),
    templateAttribute('jobs.ColleaguesService/ListColleagueActivity', 'Types'),
    templateAttribute('qualifications.QualificationsService/UpdateQualification', 'Access'),
    templateAttribute('calendar.CalendarService/CreateCalendar', 'Fields'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Types'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Types'),
    templateAttribute('wiki.WikiService/UpdatePage', 'Fields'),
    templateAttribute('mailer.MailerService/CreateOrUpdateEmail', 'Fields'),
] as const satisfies readonly TemplateAttribute[];

const gangAttributes = [
    templateAttribute('livemap.LivemapService/CreateOrUpdateMarker', 'Access'),
    templateAttribute('livemap.LivemapService/DeleteMarker', 'Access'),
    templateAttribute('livemap.LivemapService/Stream', 'Markers'),
    templateAttribute('jobs.ConductService/ListConductEntries', 'Access'),
    templateAttribute('jobs.TimeclockService/ListTimeclock', 'Access'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Access'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Access'),
    templateAttribute('jobs.ColleaguesService/ListColleagueActivity', 'Types'),
    templateAttribute('qualifications.QualificationsService/UpdateQualification', 'Access'),
    templateAttribute('calendar.CalendarService/CreateCalendar', 'Fields'),
    templateAttribute('jobs.ColleaguesService/GetColleague', 'Types'),
    templateAttribute('jobs.ColleaguesService/SetColleagueProps', 'Types'),
    templateAttribute('wiki.WikiService/UpdatePage', 'Fields'),
    templateAttribute('mailer.MailerService/CreateOrUpdateEmail', 'Fields'),
] as const satisfies readonly TemplateAttribute[];

export const policeJobTemplate: Template = {
    permissions: [...policePermissions],
    attributes: [...policeAttributes],
};

export const neutralJobTemplate: Template = {
    permissions: [...neutralPermissions],
    attributes: [...neutralAttributes],
};

export const gangJobTemplate: Template = {
    permissions: [...gangPermissions],
    attributes: [...gangAttributes],
};

export const fullPermsTemplate: Template = {
    permissions: [],
    attributes: [],
    grantAllPermissions: true,
};
