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
    'citizens.CitizensService/GetUser',
    'citizens.CitizensService/ListCitizens',
    'citizens.CitizensService/ListUserActivity',
    'citizens.CitizensService/SetUserProps',
    'vehicles.VehiclesService/ListVehicles',
    'documents.DocumentsService/AddDocumentReference',
    'documents.DocumentsService/AddDocumentRelation',
    'documents.TemplatesService/CreateTemplate',
    'documents.DocumentsService/DeleteDocument',
    'documents.CategoriesService/DeleteCategory',
    'documents.CommentsService/DeleteComment',
    'documents.TemplatesService/DeleteTemplate',
    'documents.CategoriesService/ListCategories',
    'documents.DocumentsService/ListDocuments',
    'documents.TemplatesService/ListTemplates',
    'documents.DocumentsService/ListUserDocuments',
    'documents.DocumentsService/UpdateDocument',
    'livemap.LivemapService/Stream',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'settings.SettingsService/ViewAuditLog',
    'documents.DocumentsService/ToggleDocument',
    'centrum.DispatchesService/CreateDispatch',
    'centrum.CentrumService/Stream',
    'centrum.DispatchesService/UpdateDispatch',
    'centrum.UnitsService/CreateOrUpdateUnit',
    'centrum.UnitsService/DeleteUnit',
    'centrum.DispatchesService/TakeDispatch',
    'centrum.CentrumService/TakeControl',
    'centrum.CentrumService/UpdateSettings',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'centrum.DispatchesService/DeleteDispatch',
    'documents.DocumentsService/ListDocumentActivity',
    'documents.DocumentsService/ChangeDocumentOwner',
    'documents.DocumentsService/CreateDocumentReq',
    'documents.DocumentsService/DeleteDocumentReq',
    'documents.DocumentsService/ListDocumentReqs',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.TimeclockService/ListTimeclock',
    'jobs.JobsService/SetMOTD',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.ColleaguesService/ListColleagueActivity',
    'qualifications.QualificationsService/DeleteQualification',
    'qualifications.QualificationsService/ListQualifications',
    'qualifications.QualificationsService/UpdateQualification',
    'calendar.CalendarService/CreateCalendar',
    'documents.DocumentsService/ToggleDocumentPin',
    'wiki.WikiService/UpdatePage',
    'wiki.WikiService/DeletePage',
    'wiki.WikiService/ListPageActivity',
    'wiki.WikiService/ListPages',
    'mailer.MailerService/CreateOrUpdateEmail',
    'mailer.MailerService/DeleteEmail',
    'mailer.MailerService/ListEmails',
    'jobs.ColleaguesService/ManageLabels',
    'documents.DocumentsService/SetDocumentReminder',
    'documents.CategoriesService/CreateOrUpdateCategory',
    'citizens.LabelsService/CreateOrUpdateLabel',
    'citizens.LabelsService/DeleteLabel',
    'centrum.CentrumService/UpdateDispatchers',
    'vehicles.VehiclesService/SetVehicleProps',
    'documents.ApprovalService/DeleteApprovalTasks',
    'documents.ApprovalService/RevokeApproval',
    'documents.ApprovalService/UpsertApprovalPolicy',
    'documents.ApprovalService/UpsertApprovalTasks',
    'documents.StampsService/DeleteStamp',
    'documents.StampsService/ListUsableStamps',
    'documents.StampsService/UpsertStamp',
    'wiki.WikiService/CreatePage',
] as const satisfies readonly Perms[];

const neutralPermissions = [
    'auth.AuthService/ChooseCharacter',
    'documents.CategoriesService/ListCategories',
    'documents.DocumentsService/ListDocuments',
    'livemap.LivemapService/Stream',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.TimeclockService/ListTimeclock',
    'jobs.JobsService/SetMOTD',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.ColleaguesService/ListColleagueActivity',
    'calendar.CalendarService/CreateCalendar',
    'wiki.WikiService/UpdatePage',
    'wiki.WikiService/DeletePage',
    'wiki.WikiService/ListPageActivity',
    'wiki.WikiService/ListPages',
    'jobs.ColleaguesService/ManageLabels',
    'wiki.WikiService/CreatePage',
] as const satisfies readonly Perms[];

const gangPermissions = [
    'auth.AuthService/ChooseCharacter',
    'documents.DocumentsService/ListDocuments',
    'livemap.LivemapService/Stream',
    'settings.SettingsService/CreateRole',
    'settings.SettingsService/DeleteRole',
    'settings.SettingsService/GetJobProps',
    'settings.SettingsService/GetRoles',
    'settings.SettingsService/SetJobProps',
    'settings.SettingsService/UpdateRolePerms',
    'livemap.LivemapService/CreateOrUpdateMarker',
    'livemap.LivemapService/DeleteMarker',
    'jobs.ConductService/CreateConductEntry',
    'jobs.ConductService/DeleteConductEntry',
    'jobs.ConductService/ListConductEntries',
    'jobs.ConductService/UpdateConductEntry',
    'jobs.ColleaguesService/ListColleagues',
    'jobs.TimeclockService/ListTimeclock',
    'jobs.JobsService/SetMOTD',
    'jobs.ColleaguesService/GetColleague',
    'jobs.ColleaguesService/SetColleagueProps',
    'jobs.TimeclockService/ListInactiveEmployees',
    'jobs.ColleaguesService/ListColleagueActivity',
    'calendar.CalendarService/CreateCalendar',
    'wiki.WikiService/DeletePage',
    'mailer.MailerService/CreateOrUpdateEmail',
    'mailer.MailerService/DeleteEmail',
    'mailer.MailerService/ListEmails',
    'jobs.ColleaguesService/ManageLabels',
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
