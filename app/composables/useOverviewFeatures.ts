import type { CardElement } from '~/utils/types';
import type { Perms } from '~~/gen/ts/perms';

export type OverviewFeature = CardElement & {
    to: string;
};

export const useOverviewFeatures = () => {
    const { t } = useI18n();
    const { can } = useAuth();

    return computed<OverviewFeature[]>(() =>
        [
            {
                title: t('common.mail'),
                description: t('pages.overview.features.mailer'),
                to: '/mail',
                permission: 'mailer.MailerService/ListEmails' as Perms,
                icon: 'i-mdi-inbox-full-outline',
            },
            {
                title: t('common.citizen', 2),
                description: t('pages.overview.features.citizens'),
                to: '/citizens',
                permission: 'citizens.CitizensService/ListCitizens' as Perms,
                icon: 'i-mdi-account-multiple-outline',
            },
            {
                title: t('common.vehicle', 2),
                description: t('pages.overview.features.vehicles'),
                to: '/vehicles',
                permission: 'vehicles.VehiclesService/ListVehicles' as Perms,
                icon: 'i-mdi-car-outline',
            },
            {
                title: t('common.document', 2),
                description: t('pages.overview.features.documents'),
                to: '/documents',
                permission: 'documents.DocumentsService/ListDocuments' as Perms,
                icon: 'i-mdi-file-document-box-multiple-outline',
            },
            {
                title: t('common.job'),
                description: t('pages.overview.features.jobs'),
                to: '/jobs/overview',
                permission: 'jobs.ColleaguesService/ListColleagues' as Perms,
                icon: 'i-mdi-briefcase-outline',
            },
            {
                title: t('common.calendar'),
                description: t('pages.overview.features.calendar'),
                to: '/calendar',
                icon: 'i-mdi-calendar-outline',
            },
            {
                title: t('common.qualification', 2),
                description: t('pages.overview.features.qualifications'),
                to: '/qualifications',
                permission: 'qualifications.QualificationsService/ListQualifications' as Perms,
                icon: 'i-mdi-school-outline',
            },
            {
                title: t('common.livemap'),
                description: t('pages.overview.features.livemap'),
                to: '/livemap',
                permission: 'livemap.LivemapService/Stream' as Perms,
                icon: 'i-mdi-map-outline',
            },
            {
                title: t('common.dispatch_center'),
                description: t('pages.overview.features.centrum'),
                to: '/dispatch',
                permission: 'centrum.CentrumService/TakeControl' as Perms,
                icon: 'i-mdi-car-emergency',
            },
            {
                title: t('common.wiki'),
                description: t('pages.overview.features.wiki'),
                to: '/wiki',
                permission: 'wiki.WikiService/ListPages' as Perms,
                icon: 'i-mdi-brain',
            },
        ].filter((item) => item.permission === undefined || can(item.permission).value),
    );
};
