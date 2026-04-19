import type { AttributeValues } from '~~/gen/ts/resources/permissions/attributes/attributes';

export type Template = {
    permissions: PermissionTemplate[];
    attributes: AttributeTemplate[];
};

export type PermissionTemplate = {
    category: string;
    name: string;
};

export type AttributeTemplate = {
    category: string;
    name: string;
    key: string;
    type: string;
    validValues?: AttributeValues;
};

export const policeJobTemplate: Template = {
    permissions: [
        { category: 'auth.AuthService', name: 'ChooseCharacter' },
        { category: 'citizens.CitizensService', name: 'GetUser' },
        { category: 'citizens.CitizensService', name: 'ListCitizens' },
        { category: 'citizens.CitizensService', name: 'ListUserActivity' },
        { category: 'citizens.CitizensService', name: 'SetUserProps' },
        { category: 'vehicles.VehiclesService', name: 'ListVehicles' },
        { category: 'documents.DocumentsService', name: 'AddDocumentReference' },
        { category: 'documents.DocumentsService', name: 'AddDocumentRelation' },
        { category: 'documents.DocumentsService', name: 'CreateTemplate' },
        { category: 'documents.DocumentsService', name: 'DeleteDocument' },
        { category: 'documents.DocumentsService', name: 'DeleteCategory' },
        { category: 'documents.DocumentsService', name: 'DeleteComment' },
        { category: 'documents.DocumentsService', name: 'DeleteTemplate' },
        { category: 'documents.DocumentsService', name: 'ListCategories' },
        { category: 'documents.DocumentsService', name: 'ListDocuments' },
        { category: 'documents.DocumentsService', name: 'ListTemplates' },
        { category: 'documents.DocumentsService', name: 'ListUserDocuments' },
        { category: 'documents.DocumentsService', name: 'UpdateDocument' },
        { category: 'livemap.LivemapService', name: 'Stream' },
        { category: 'settings.SettingsService', name: 'CreateRole' },
        { category: 'settings.SettingsService', name: 'DeleteRole' },
        { category: 'settings.SettingsService', name: 'GetJobProps' },
        { category: 'settings.SettingsService', name: 'GetRoles' },
        { category: 'settings.SettingsService', name: 'SetJobProps' },
        { category: 'settings.SettingsService', name: 'UpdateRolePerms' },
        { category: 'settings.SettingsService', name: 'ViewAuditLog' },
        { category: 'documents.DocumentsService', name: 'ToggleDocument' },
        { category: 'centrum.CentrumService', name: 'CreateDispatch' },
        { category: 'centrum.CentrumService', name: 'Stream' },
        { category: 'centrum.CentrumService', name: 'UpdateDispatch' },
        { category: 'centrum.CentrumService', name: 'CreateOrUpdateUnit' },
        { category: 'centrum.CentrumService', name: 'DeleteUnit' },
        { category: 'centrum.CentrumService', name: 'TakeDispatch' },
        { category: 'centrum.CentrumService', name: 'TakeControl' },
        { category: 'centrum.CentrumService', name: 'UpdateSettings' },
        { category: 'livemap.LivemapService', name: 'CreateOrUpdateMarker' },
        { category: 'livemap.LivemapService', name: 'DeleteMarker' },
        { category: 'centrum.CentrumService', name: 'DeleteDispatch' },
        { category: 'documents.DocumentsService', name: 'ListDocumentActivity' },
        { category: 'documents.DocumentsService', name: 'ChangeDocumentOwner' },
        { category: 'documents.DocumentsService', name: 'CreateDocumentReq' },
        { category: 'documents.DocumentsService', name: 'DeleteDocumentReq' },
        { category: 'documents.DocumentsService', name: 'ListDocumentReqs' },
        { category: 'jobs.ConductService', name: 'CreateConductEntry' },
        { category: 'jobs.ConductService', name: 'DeleteConductEntry' },
        { category: 'jobs.ConductService', name: 'ListConductEntries' },
        { category: 'jobs.ConductService', name: 'UpdateConductEntry' },
        { category: 'jobs.JobsService', name: 'ListColleagues' },
        { category: 'jobs.TimeclockService', name: 'ListTimeclock' },
        { category: 'jobs.JobsService', name: 'SetMOTD' },
        { category: 'jobs.JobsService', name: 'GetColleague' },
        { category: 'jobs.JobsService', name: 'SetColleagueProps' },
        { category: 'jobs.TimeclockService', name: 'ListInactiveEmployees' },
        { category: 'jobs.JobsService', name: 'ListColleagueActivity' },
        { category: 'qualifications.QualificationsService', name: 'DeleteQualification' },
        { category: 'qualifications.QualificationsService', name: 'ListQualifications' },
        { category: 'qualifications.QualificationsService', name: 'UpdateQualification' },
        { category: 'completor.CompletorService', name: 'CompleteCitizenLabels' },
        { category: 'calendar.CalendarService', name: 'CreateCalendar' },
        { category: 'documents.DocumentsService', name: 'ToggleDocumentPin' },
        { category: 'wiki.WikiService', name: 'UpdatePage' },
        { category: 'wiki.WikiService', name: 'DeletePage' },
        { category: 'wiki.WikiService', name: 'ListPageActivity' },
        { category: 'wiki.WikiService', name: 'ListPages' },
        { category: 'mailer.MailerService', name: 'CreateOrUpdateEmail' },
        { category: 'mailer.MailerService', name: 'DeleteEmail' },
        { category: 'mailer.MailerService', name: 'ListEmails' },
        { category: 'jobs.JobsService', name: 'ManageLabels' },
        { category: 'documents.DocumentsService', name: 'SetDocumentReminder' },
        { category: 'documents.DocumentsService', name: 'CreateOrUpdateCategory' },
        { category: 'citizens.CitizensService', name: 'ManageLabels' },
        { category: 'centrum.CentrumService', name: 'UpdateDispatchers' },
        { category: 'vehicles.VehiclesService', name: 'SetVehicleProps' },
        { category: 'documents.ApprovalService', name: 'DeleteApprovalTasks' },
        { category: 'documents.ApprovalService', name: 'RevokeApproval' },
        { category: 'documents.ApprovalService', name: 'UpsertApprovalPolicy' },
        { category: 'documents.ApprovalService', name: 'UpsertApprovalTasks' },
        { category: 'documents.StampsService', name: 'DeleteStamp' },
        { category: 'documents.StampsService', name: 'ListUsableStamps' },
        { category: 'documents.StampsService', name: 'UpsertStamp' },
        { category: 'wiki.WikiService', name: 'CreatePage' },
    ],
    attributes: [
        {
            category: 'citizens.CitizensService',
            name: 'GetUser',
            key: 'Jobs',
            type: 'JobGradeList',
            validValues: {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: { fineGrained: false, jobs: {}, grades: {} },
                },
            },
        },
        {
            category: 'citizens.CitizensService',
            name: 'ListCitizens',
            key: 'Fields',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [
                            'PhoneNumber',
                            'Licenses',
                            'UserProps.Wanted',
                            'UserProps.Job',
                            'UserProps.TrafficInfractionPoints',
                            'UserProps.OpenFines',
                            'UserProps.Mugshot',
                            'UserProps.Labels',
                            'UserProps.Email',
                        ],
                    },
                },
            },
        },
        {
            category: 'citizens.CitizensService',
            name: 'ListUserActivity',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['SourceUser', 'Own'] } } },
        },
        {
            category: 'citizens.CitizensService',
            name: 'SetUserProps',
            key: 'Fields',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Wanted', 'Job', 'TrafficInfractionPoints', 'Mugshot', 'Labels'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'Stream',
            key: 'Players',
            type: 'JobGradeList',
            validValues: {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: { fineGrained: false, jobs: {}, grades: {} },
                },
            },
        },
        {
            category: 'documents.DocumentsService',
            name: 'DeleteDocument',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'documents.DocumentsService',
            name: 'DeleteComment',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'documents.DocumentsService',
            name: 'ListCategories',
            key: 'Jobs',
            type: 'JobList',
            validValues: { validValues: { oneofKind: 'jobList', jobList: { strings: [] } } },
        },
        {
            category: 'documents.DocumentsService',
            name: 'UpdateDocument',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'documents.DocumentsService',
            name: 'ToggleDocument',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'CreateOrUpdateMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'DeleteMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'Stream',
            key: 'Markers',
            type: 'JobList',
            validValues: { validValues: { oneofKind: 'jobList', jobList: { strings: [] } } },
        },
        {
            category: 'documents.DocumentsService',
            name: 'ChangeDocumentOwner',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'documents.DocumentsService',
            name: 'CreateDocumentReq',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Access', 'Closure', 'Update', 'Deletion', 'OwnerChange'] },
                },
            },
        },
        {
            category: 'jobs.ConductService',
            name: 'ListConductEntries',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Own', 'All'] } } },
        },
        {
            category: 'jobs.TimeclockService',
            name: 'ListTimeclock',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['All'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'ListColleagueActivity',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: ['HIRED', 'FIRED', 'PROMOTED', 'DEMOTED', 'ABSENCE_DATE', 'NOTE', 'LABELS', 'NAME'],
                    },
                },
            },
        },
        {
            category: 'qualifications.QualificationsService',
            name: 'DeleteQualification',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'qualifications.QualificationsService',
            name: 'UpdateQualification',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'completor.CompletorService',
            name: 'CompleteCitizenLabels',
            key: 'Jobs',
            type: 'JobList',
            validValues: { validValues: { oneofKind: 'jobList', jobList: { strings: [] } } },
        },
        {
            category: 'calendar.CalendarService',
            name: 'CreateCalendar',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job', 'Public'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Types',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Note', 'Labels'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['AbsenceDate', 'Note', 'Labels', 'Name'] },
                },
            },
        },
        {
            category: 'wiki.WikiService',
            name: 'UpdatePage',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Public'] } } },
        },
        {
            category: 'mailer.MailerService',
            name: 'CreateOrUpdateEmail',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job'] } } },
        },
        {
            category: 'documents.DocumentsService',
            name: 'ToggleDocumentPin',
            key: 'Types',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['JobWide'] } } },
        },
        {
            category: 'qualifications.QualificationsService',
            name: 'UpdateQualification',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Public'] } } },
        },
        {
            category: 'centrum.CentrumService',
            name: 'UpdateSettings',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Shared', 'Public'] } } },
        },
        {
            category: 'vehicles.VehiclesService',
            name: 'ListVehicles',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Wanted'] } } },
        },
        {
            category: 'vehicles.VehiclesService',
            name: 'SetVehicleProps',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Wanted'] } } },
        },
    ],
};

export const neutralJobTemplate: Template = {
    permissions: [
        { category: 'auth.AuthService', name: 'ChooseCharacter' },
        { category: 'documents.DocumentsService', name: 'ListCategories' },
        { category: 'documents.DocumentsService', name: 'ListDocuments' },
        { category: 'livemap.LivemapService', name: 'Stream' },
        { category: 'settings.SettingsService', name: 'CreateRole' },
        { category: 'settings.SettingsService', name: 'DeleteRole' },
        { category: 'settings.SettingsService', name: 'GetJobProps' },
        { category: 'settings.SettingsService', name: 'GetRoles' },
        { category: 'settings.SettingsService', name: 'SetJobProps' },
        { category: 'settings.SettingsService', name: 'UpdateRolePerms' },
        { category: 'livemap.LivemapService', name: 'CreateOrUpdateMarker' },
        { category: 'livemap.LivemapService', name: 'DeleteMarker' },
        { category: 'jobs.ConductService', name: 'CreateConductEntry' },
        { category: 'jobs.ConductService', name: 'DeleteConductEntry' },
        { category: 'jobs.ConductService', name: 'ListConductEntries' },
        { category: 'jobs.ConductService', name: 'UpdateConductEntry' },
        { category: 'jobs.JobsService', name: 'ListColleagues' },
        { category: 'jobs.TimeclockService', name: 'ListTimeclock' },
        { category: 'jobs.JobsService', name: 'SetMOTD' },
        { category: 'jobs.JobsService', name: 'GetColleague' },
        { category: 'jobs.JobsService', name: 'SetColleagueProps' },
        { category: 'jobs.TimeclockService', name: 'ListInactiveEmployees' },
        { category: 'jobs.JobsService', name: 'ListColleagueActivity' },
        { category: 'calendar.CalendarService', name: 'CreateCalendar' },
        { category: 'wiki.WikiService', name: 'UpdatePage' },
        { category: 'wiki.WikiService', name: 'DeletePage' },
        { category: 'wiki.WikiService', name: 'ListPageActivity' },
        { category: 'wiki.WikiService', name: 'ListPages' },
        { category: 'jobs.JobsService', name: 'ManageLabels' },
        { category: 'wiki.WikiService', name: 'CreatePage' },
    ],
    attributes: [
        {
            category: 'documents.DocumentsService',
            name: 'ListCategories',
            key: 'Jobs',
            type: 'JobList',
            validValues: { validValues: { oneofKind: 'jobList', jobList: { strings: [] } } },
        },
        {
            category: 'livemap.LivemapService',
            name: 'CreateOrUpdateMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'DeleteMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'Stream',
            key: 'Markers',
            type: 'JobList',
            validValues: undefined,
        },
        {
            category: 'jobs.ConductService',
            name: 'ListConductEntries',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Own', 'All'] } } },
        },
        {
            category: 'jobs.TimeclockService',
            name: 'ListTimeclock',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['All'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'ListColleagueActivity',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: ['HIRED', 'FIRED', 'PROMOTED', 'DEMOTED', 'ABSENCE_DATE', 'NOTE', 'LABELS', 'NAME'],
                    },
                },
            },
        },
        {
            category: 'qualifications.QualificationsService',
            name: 'UpdateQualification',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'calendar.CalendarService',
            name: 'CreateCalendar',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job', 'Public'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Types',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Note', 'Labels'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['AbsenceDate', 'Note', 'Labels', 'Name'] },
                },
            },
        },
        {
            category: 'wiki.WikiService',
            name: 'UpdatePage',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Public'] } } },
        },
        {
            category: 'mailer.MailerService',
            name: 'CreateOrUpdateEmail',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job'] } } },
        },
    ],
};

export const gangJobTemplate: Template = {
    permissions: [
        { category: 'auth.AuthService', name: 'ChooseCharacter' },
        { category: 'documents.DocumentsService', name: 'ListDocuments' },
        { category: 'livemap.LivemapService', name: 'Stream' },
        { category: 'settings.SettingsService', name: 'CreateRole' },
        { category: 'settings.SettingsService', name: 'DeleteRole' },
        { category: 'settings.SettingsService', name: 'GetJobProps' },
        { category: 'settings.SettingsService', name: 'GetRoles' },
        { category: 'settings.SettingsService', name: 'SetJobProps' },
        { category: 'settings.SettingsService', name: 'UpdateRolePerms' },
        { category: 'livemap.LivemapService', name: 'CreateOrUpdateMarker' },
        { category: 'livemap.LivemapService', name: 'DeleteMarker' },
        { category: 'jobs.ConductService', name: 'CreateConductEntry' },
        { category: 'jobs.ConductService', name: 'DeleteConductEntry' },
        { category: 'jobs.ConductService', name: 'ListConductEntries' },
        { category: 'jobs.ConductService', name: 'UpdateConductEntry' },
        { category: 'jobs.JobsService', name: 'ListColleagues' },
        { category: 'jobs.TimeclockService', name: 'ListTimeclock' },
        { category: 'jobs.JobsService', name: 'SetMOTD' },
        { category: 'jobs.JobsService', name: 'GetColleague' },
        { category: 'jobs.JobsService', name: 'SetColleagueProps' },
        { category: 'jobs.TimeclockService', name: 'ListInactiveEmployees' },
        { category: 'jobs.JobsService', name: 'ListColleagueActivity' },
        { category: 'calendar.CalendarService', name: 'CreateCalendar' },
        { category: 'wiki.WikiService', name: 'DeletePage' },
        { category: 'mailer.MailerService', name: 'CreateOrUpdateEmail' },
        { category: 'mailer.MailerService', name: 'DeleteEmail' },
        { category: 'mailer.MailerService', name: 'ListEmails' },
        { category: 'jobs.JobsService', name: 'ManageLabels' },
    ],
    attributes: [
        {
            category: 'livemap.LivemapService',
            name: 'CreateOrUpdateMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'DeleteMarker',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'livemap.LivemapService',
            name: 'Stream',
            key: 'Markers',
            type: 'JobList',
            validValues: undefined,
        },
        {
            category: 'jobs.ConductService',
            name: 'ListConductEntries',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Own', 'All'] } } },
        },
        {
            category: 'jobs.TimeclockService',
            name: 'ListTimeclock',
            key: 'Access',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['All'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'jobs.JobsService',
            name: 'ListColleagueActivity',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: ['HIRED', 'FIRED', 'PROMOTED', 'DEMOTED', 'ABSENCE_DATE', 'NOTE', 'LABELS', 'NAME'],
                    },
                },
            },
        },
        {
            category: 'qualifications.QualificationsService',
            name: 'UpdateQualification',
            key: 'Access',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['Own', 'Lower_Rank', 'Same_Rank', 'Any'] },
                },
            },
        },
        {
            category: 'calendar.CalendarService',
            name: 'CreateCalendar',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'GetColleague',
            key: 'Types',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Note', 'Labels'] } } },
        },
        {
            category: 'jobs.JobsService',
            name: 'SetColleagueProps',
            key: 'Types',
            type: 'StringList',
            validValues: {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: { strings: ['AbsenceDate', 'Note', 'Labels', 'Name'] },
                },
            },
        },
        {
            category: 'wiki.WikiService',
            name: 'UpdatePage',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: [] } } },
        },
        {
            category: 'mailer.MailerService',
            name: 'CreateOrUpdateEmail',
            key: 'Fields',
            type: 'StringList',
            validValues: { validValues: { oneofKind: 'stringList', stringList: { strings: ['Job'] } } },
        },
    ],
};
