import type { NavigationMenuItem } from '@nuxt/ui';
import type { RoutePathSchema } from '@typed-router';
import { isRoute } from '~/utils/route';
import type { CardElement } from '~/utils/types';
import type { Perms } from '~~/gen/ts/perms';

export type AppFeature = {
    id: string;
    label: string;
    labelCount?: number;
    description?: string;
    to: string;
    startpagePath?: RoutePathSchema;
    icon: string;
    permission?: Perms | Perms[];
    shortcut?: string[];
    activePaths?: string[];
    overview?: boolean;
    startpage?: boolean;
    children?: AppFeature[];
};

export type AppOverviewFeature = Omit<CardElement, 'to'> & { to: string };

const appFeatures: AppFeature[] = [
    {
        id: 'overview',
        label: 'common.overview',
        to: '/overview',
        icon: 'i-mdi-home-outline',
        shortcut: ['G', 'H'],
        startpage: true,
        startpagePath: '/overview',
    },
    {
        id: 'mail',
        label: 'common.mail',
        to: '/mail',
        icon: 'i-mdi-inbox-full-outline',
        shortcut: ['G', 'E'],
        permission: 'mailer.MailerService/ListEmails' as Perms,
        activePaths: ['/mail'],
        overview: true,
        description: 'pages.overview.features.mailer',
        startpage: true,
        startpagePath: '/mail/:thread?',
    },
    {
        id: 'citizens',
        label: 'common.citizen',
        labelCount: 2,
        to: '/citizens',
        icon: 'i-mdi-account-multiple-outline',
        shortcut: ['G', 'C'],
        permission: 'citizens.CitizensService/ListCitizens' as Perms,
        activePaths: ['/citizens'],
        overview: true,
        description: 'pages.overview.features.citizens',
        startpage: true,
    },
    {
        id: 'vehicles',
        label: 'common.vehicle',
        labelCount: 2,
        to: '/vehicles',
        icon: 'i-mdi-car-outline',
        shortcut: ['G', 'V'],
        permission: 'vehicles.VehiclesService/ListVehicles' as Perms,
        overview: true,
        description: 'pages.overview.features.vehicles',
        startpage: true,
    },
    {
        id: 'documents',
        label: 'common.document',
        labelCount: 2,
        to: '/documents',
        icon: 'i-mdi-file-document-box-multiple-outline',
        shortcut: ['G', 'D'],
        permission: 'documents.DocumentsService/ListDocuments' as Perms,
        activePaths: ['/documents'],
        overview: true,
        description: 'pages.overview.features.documents',
        startpage: true,
        startpagePath: '/documents/',
        children: [
            {
                id: 'document-approvals',
                label: 'common.approvals',
                labelCount: 2,
                to: '/documents/approvals',
                icon: 'i-mdi-approval',
            },
            {
                id: 'document-stats',
                label: 'common.stats',
                to: '/documents/stats',
                icon: 'i-mdi-graph-box-outline',
                permission: 'documents.StatsService/GetStats' as Perms,
            },
        ],
    },
    {
        id: 'jobs',
        label: 'common.job',
        to: '/jobs/overview',
        icon: 'i-mdi-briefcase-outline',
        shortcut: ['G', 'J'],
        permission: 'jobs.ColleaguesService/ListColleagues' as Perms,
        activePaths: ['/jobs'],
        overview: true,
        description: 'pages.overview.features.jobs',
        startpage: true,
        children: [
            { id: 'jobs-overview', label: 'common.overview', to: '/jobs/overview', icon: 'i-mdi-briefcase-outline' },
            {
                id: 'jobs-colleagues',
                label: 'common.colleague',
                labelCount: 2,
                to: '/jobs/colleagues',
                icon: 'i-mdi-account-group',
                permission: 'jobs.ColleaguesService/ListColleagues' as Perms,
                activePaths: ['/jobs/colleagues'],
                children: [
                    {
                        id: 'jobs-colleagues-stats',
                        label: 'pages.jobs.colleagues.stats.title',
                        to: '/jobs/colleagues/stats',
                        icon: 'i-mdi-chart-timeline-variant-shimmer',
                        permission: 'jobs.StatsService/GetStats' as Perms,
                    },
                    {
                        id: 'jobs-colleagues-labels',
                        label: 'pages.jobs.colleagues.labels.title',
                        to: '/jobs/colleagues/labels',
                        icon: 'i-mdi-label-multiple',
                        permission: ['jobs.ColleaguesService/CreateOrUpdateLabel'] as Perms[],
                    },
                ],
            },
            {
                id: 'jobs-activity',
                label: 'common.activity',
                to: '/jobs/activity',
                icon: 'i-mdi-pulse',
                permission: 'jobs.ColleaguesService/ListColleagueActivity' as Perms,
            },
            {
                id: 'jobs-timeclock',
                label: 'common.timeclock',
                to: '/jobs/timeclock',
                icon: 'i-mdi-timeline-clock',
                permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
                activePaths: ['/jobs/timeclock'],
                children: [
                    {
                        id: 'jobs-timeclock-inactive',
                        label: 'common.inactive_colleagues',
                        to: '/jobs/timeclock/inactive',
                        icon: 'i-mdi-account-remove',
                        permission: 'jobs.TimeclockService/ListInactiveEmployees' as Perms,
                    },
                ],
            },
            {
                id: 'jobs-conduct',
                label: 'common.conduct_register',
                labelCount: 2,
                to: '/jobs/conduct',
                icon: 'i-mdi-list-status',
                permission: 'jobs.ConductService/ListConductEntries' as Perms,
            },
            {
                id: 'jobs-groups',
                label: 'common.group',
                labelCount: 2,
                to: '/jobs/groups',
                icon: 'i-mdi-users-group-outline',
                permission: 'jobs.GroupsService/ListGroups' as Perms,
            },
        ],
    },
    {
        id: 'calendar',
        label: 'common.calendar',
        to: '/calendar',
        icon: 'i-mdi-calendar-outline',
        shortcut: ['G', 'K'],
        activePaths: ['/calendar'],
        overview: true,
        description: 'pages.overview.features.calendar',
        startpage: true,
    },
    {
        id: 'qualifications',
        label: 'common.qualification',
        labelCount: 2,
        to: '/qualifications',
        icon: 'i-mdi-school-outline',
        shortcut: ['G', 'Q'],
        permission: 'qualifications.QualificationsService/ListQualifications' as Perms,
        activePaths: ['/qualifications'],
        overview: true,
        description: 'pages.overview.features.qualifications',
        startpage: true,
    },
    {
        id: 'livemap',
        label: 'common.livemap',
        to: '/livemap',
        icon: 'i-mdi-map-outline',
        shortcut: ['G', 'M'],
        permission: 'livemap.LivemapService/Stream' as Perms,
        overview: true,
        description: 'pages.overview.features.livemap',
        startpage: true,
    },
    {
        id: 'dispatch',
        label: 'common.dispatch_center',
        to: '/dispatch',
        icon: 'i-mdi-car-emergency',
        shortcut: ['G', 'W'],
        permission: 'centrum.CentrumService/TakeControl' as Perms,
        activePaths: ['/dispatch', '/centrum'],
        overview: true,
        description: 'pages.overview.features.centrum',
        startpage: true,
    },
    {
        id: 'wiki',
        label: 'common.wiki',
        to: '/wiki',
        icon: 'i-mdi-brain',
        shortcut: ['G', 'L'],
        permission: 'wiki.WikiService/ListPages' as Perms,
        activePaths: ['/wiki'],
        overview: true,
        description: 'pages.overview.features.wiki',
        startpage: true,
    },
    {
        id: 'settings',
        label: 'common.control_panel',
        to: '/settings',
        icon: 'i-mdi-cog-outline',
        shortcut: ['G', 'P'],
        permission: 'settings.SettingsService/GetJobProps' as Perms,
        activePaths: ['/settings'],
        children: [
            {
                id: 'settings-props',
                label: 'components.settings.job_props.job_properties',
                to: '/settings/props',
                icon: 'i-mdi-tune',
                permission: 'settings.SettingsService/SetJobProps' as Perms,
            },
            {
                id: 'settings-roles',
                label: 'common.role',
                labelCount: 2,
                to: '/settings/roles',
                icon: 'i-mdi-account-group',
                permission: 'settings.SettingsService/GetRoles' as Perms,
            },
        ],
    },
];

export const useAppFeatures = () => {
    const { t } = useI18n();
    const { can } = useAuth();
    const route = useRoute();

    const translate = (feature: AppFeature) =>
        feature.labelCount === undefined ? t(feature.label) : t(feature.label, feature.labelCount);

    const visibleFeatures = (features: AppFeature[]): AppFeature[] =>
        features
            .filter((feature) => feature.permission === undefined || can(feature.permission).value)
            .map((feature) => ({
                ...feature,
                children: feature.children ? visibleFeatures(feature.children) : undefined,
            }));

    const toNavigationItem = (feature: AppFeature, root = false): NavigationMenuItem => ({
        id: feature.id,
        label: translate(feature),
        icon: feature.icon,
        to: feature.to,
        permission: feature.permission,
        active: feature.activePaths?.some((path) => isRoute(route.path, path)),
        ...(root
            ? {
                  tooltip: feature.shortcut ? { text: translate(feature), kbds: feature.shortcut } : undefined,
                  kbds: feature.shortcut,
                  defaultOpen: feature.children ? false : undefined,
              }
            : {}),
        children: feature.children?.map((child) => toNavigationItem(child)),
    });

    const navigationItems = computed<NavigationMenuItem[]>(() =>
        visibleFeatures(appFeatures).map((feature) => toNavigationItem(feature, true)),
    );

    const overviewItems = computed<AppOverviewFeature[]>(() =>
        visibleFeatures(appFeatures)
            .filter((feature) => feature.overview)
            .map((feature) => ({
                label: translate(feature),
                description: feature.description ? t(feature.description) : undefined,
                to: feature.to,
                icon: feature.icon,
                permission: feature.permission,
            })),
    );

    const startpageItems = computed(() =>
        visibleFeatures(appFeatures)
            .filter((feature) => feature.startpage)
            .map((feature) => ({
                label: translate(feature),
                path: (feature.startpagePath ?? feature.to) as RoutePathSchema,
            })),
    );

    return { navigationItems, overviewItems, startpageItems };
};
