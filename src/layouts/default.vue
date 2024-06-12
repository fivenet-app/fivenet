<script lang="ts" setup>
import type { Group } from '#ui/types';
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';
import BodyCheckupModal from '~/components/quickbuttons/bodycheckup/BodyCheckupModal.vue';
import PenaltyCalculatorModal from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorModal.vue';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';
import UserDropdown from '~/components/UserDropdown.vue';
import HelpSlideover from '~/components/HelpSlideover.vue';
import NotificationsSlideover from '~/components/NotificationsSlideover.vue';
import DashboardSidebarLinks from '~/components/partials/dashboard/DashboardSidebarLinks.vue';
import { messengerDB } from '~/store/messenger';

const authStore = useAuthStore();
const { activeChar, jobProps } = storeToRefs(authStore);

const { t } = useI18n();

const { isHelpSlideoverOpen } = useDashboard();

const modal = useModal();

const unreadThreadCount = useDexieLiveQuery(() => messengerDB.threads.filter((t) => !!t.userState?.unread).count(), {
    initialValue: 0,
});

const links = computed(() =>
    [
        {
            label: t('common.overview'),
            icon: 'i-mdi-home-outline',
            to: '/overview',
            tooltip: {
                text: t('common.overview'),
                shortcuts: ['G', 'H'],
            },
        },
        {
            label: t('common.messenger'),
            icon: 'i-mdi-conversation-outline',
            to: '/messenger',
            badge: unreadThreadCount.value > 0 ? unreadThreadCount.value.toString() : undefined,
            tooltip: {
                text: t('common.messenger'),
                shortcuts: ['G', 'I'],
            },
            permission: 'MessengerService.ListThreads' as Perms,
        },
        {
            label: t('common.citizen'),
            icon: 'i-mdi-account-multiple-outline',
            to: '/citizens',
            tooltip: {
                text: t('common.citizen'),
                shortcuts: ['G', 'C'],
            },
            permission: 'CitizenStoreService.ListCitizens' as Perms,
        },
        {
            label: t('common.vehicle', 2),
            icon: 'i-mdi-car-outline',
            to: '/vehicles',
            tooltip: {
                text: t('common.vehicle', 2),
                shortcuts: ['G', 'V'],
            },
            permission: 'DMVService.ListVehicles' as Perms,
        },
        {
            label: t('common.document', 2),
            icon: 'i-mdi-file-document-box-multiple-outline',
            to: '/documents',
            tooltip: {
                text: t('common.document', 2),
                shortcuts: ['G', 'D'],
            },
            permission: 'DocStoreService.ListDocuments' as Perms,
        },
        {
            label: t('common.job'),
            icon: 'i-mdi-briefcase-outline',
            to: '/jobs/overview',
            tooltip: {
                text: t('common.job'),
                shortcuts: ['G', 'J'],
            },
            defaultOpen: false,
            children: [
                {
                    label: t('common.overview'),
                    to: '/jobs/overview',
                },
                {
                    label: t('common.colleague', 2),
                    to: '/jobs/colleagues',
                    permission: 'JobsService.ListColleagues' as Perms,
                },
                {
                    label: t('common.activity'),
                    to: '/jobs/activity',
                    permission: 'JobsService.ListColleagueActivity' as Perms,
                },
                {
                    label: t('common.timeclock'),
                    to: '/jobs/timeclock',
                    permission: 'JobsTimeclockService.ListTimeclock' as Perms,
                },
                {
                    label: t('common.conduct_register', 2),
                    to: '/jobs/conduct',
                    permission: 'JobsConductService.ListConductEntries' as Perms,
                },
            ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
            permission: 'JobsService.ListColleagues' as Perms,
        },
        {
            label: t('common.calendar'),
            icon: 'i-mdi-calendar-outline',
            to: '/calendar',
            tooltip: {
                text: t('common.calendar'),
                shortcuts: ['G', 'K'],
            },
        },
        {
            label: t('common.qualification', 2),
            icon: 'i-mdi-school-outline',
            to: '/qualifications',
            tooltip: {
                text: t('common.qualification', 2),
                shortcuts: ['G', 'Q'],
            },
            permission: 'QualificationsService.ListQualifications' as Perms,
        },
        {
            label: t('common.livemap'),
            icon: 'i-mdi-map-outline',
            to: '/livemap',
            tooltip: {
                text: t('common.livemap'),
                shortcuts: ['G', 'M'],
            },
            permission: 'LivemapperService.Stream' as Perms,
        },
        {
            label: t('common.dispatch_center'),
            icon: 'i-mdi-car-emergency',
            to: '/centrum',
            tooltip: {
                text: t('common.dispatch_center'),
                shortcuts: ['G', 'W'],
            },
            permission: 'CentrumService.TakeControl' as Perms,
        },
        {
            label: t('common.control_panel'),
            icon: 'i-mdi-cog-outline',
            to: '/rector',
            tooltip: {
                text: t('common.control_panel'),
                shortcuts: ['G', 'P'],
            },
            permission: 'RectorService.GetJobProps' as Perms,
        },
    ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
);

const footerLinks = [
    {
        label: t('common.help'),
        icon: 'i-mdi-question-mark-circle-outline',
        click: () => (isHelpSlideoverOpen.value = true),
    },
    {
        label: t('common.about'),
        icon: 'i-mdi-about-circle-outline',
        to: '/about',
    },
];

const groups = computed(() => [
    {
        key: 'links',
        label: t('common.goto'),
        commands: links.value.map((link) => ({ ...link, shortcuts: link.tooltip?.shortcuts })),
    },
    {
        key: 'ids',
        label: t('common.id', 2),
        commands: [
            {
                id: 'cit',
                prefix: 'CIT-',
                icon: 'i-mdi-account-multiple-outline',
            },
            {
                id: 'doc',
                prefix: 'DOC-',
                icon: 'i-mdi-file-document-box-multiple-outline',
            },
        ],
        search: async (q?: string) => {
            const defaultCommands = [
                {
                    id: 'id-doc',
                    label: `DOC-...`,
                },
                {
                    id: 'id-citizen',
                    label: `CIT-...`,
                },
            ];

            if (!q || (!q.startsWith('CIT') && !q.startsWith('DOC'))) {
                if (q && (q.startsWith('@') || q.startsWith('#'))) {
                    return [];
                }

                return defaultCommands.filter((c) => !q || c.label.includes(q));
            }

            const prefix = q.substring(0, q.indexOf('-')).toUpperCase();
            const id = q.substring(q.indexOf('-') + 1).trim();
            if (id.length > 0 && isNumber(id)) {
                if (prefix === 'CIT') {
                    return [
                        {
                            id: 'id-citizen',
                            label: `CIT-${id}`,
                            to: `/citizens/${id}`,
                        },
                    ];
                } else if (prefix === 'DOC') {
                    return [
                        {
                            id: 'id-doc',
                            label: `DOC-${id}`,
                            to: `/documents/${id}`,
                        },
                    ];
                }
            }

            return defaultCommands;
        },
    },
    {
        key: 'search',
        label: t('common.search'),
        commands: [
            {
                id: 'cit',
                label: t('common.citizen', 2),
                prefix: '@',
                icon: 'i-mdi-account-multiple-outline',
            },
            {
                id: 'doc',
                label: t('common.document', 2),
                prefix: '#',
                icon: 'i-mdi-file-document-box-multiple-outline',
            },
        ],
        search: async (q?: string) => {
            if (!q || (!q.startsWith('@') && !q.startsWith('#'))) {
                return [
                    {
                        id: 'cit',
                        label: t('common.citizen', 2),
                        prefix: '@',
                        icon: 'i-mdi-account-multiple-outline',
                    },
                    {
                        id: 'doc',
                        label: t('common.document', 2),
                        prefix: '#',
                        icon: 'i-mdi-file-document-box-multiple-outline',
                    },
                ].filter((c) => !q || c.label.includes(q));
            }

            const searchType = q[0];
            const query = q.substring(1).trim();
            switch (searchType) {
                case '#': {
                    try {
                        const call = getGRPCDocStoreClient().listDocuments({
                            pagination: {
                                offset: 0,
                                pageSize: 10,
                            },
                            orderBy: [],
                            search: query,
                            categoryIds: [],
                            creatorIds: [],
                            documentIds: [],
                        });
                        const { response } = await call;

                        return response.documents.map((d) => ({
                            id: d.id,
                            label: d.title,
                            suffix: d.state,
                            to: `/documents/${d.id}`,
                        }));
                    } catch (e) {
                        handleGRPCError(e as RpcError);
                        throw e;
                    }
                }

                case '@':
                default: {
                    try {
                        const call = getGRPCCitizenStoreClient().listCitizens({
                            pagination: {
                                offset: 0,
                                pageSize: 10,
                            },
                            search: query,
                        });
                        const { response } = await call;

                        return response.users.map((u) => ({
                            id: u.userId,
                            label: `${u.firstname} ${u.lastname}`,
                            suffix: u.dateofbirth,
                            to: `/citizens/${u.userId}`,
                        }));
                    } catch (e) {
                        handleGRPCError(e as RpcError);
                        throw e;
                    }
                }
            }
        },
    },
]);

const clipboardLink = computed(() =>
    [
        activeChar.value &&
        can(['DocStoreService.CreateDocument', 'CitizenStoreService.GetUser', 'DMVService.ListVehicles'], 'oneof').value
            ? {
                  label: t('common.clipboard'),
                  icon: 'i-mdi-clipboard-list-outline',
                  click: () => modal.open(ClipboardModal, {}),
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const quickAccessButtons = computed(() =>
    [
        jobProps.value?.quickButtons?.penaltyCalculator
            ? {
                  label: t('components.penaltycalculator.title'),
                  icon: 'i-mdi-calculator',
                  click: () => modal.open(PenaltyCalculatorModal),
              }
            : undefined,
        jobProps.value?.quickButtons?.bodyCheckup
            ? {
                  label: t('components.bodycheckup.title'),
                  icon: 'i-mdi-human',
                  click: () => modal.open(BodyCheckupModal, {}),
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <UDashboardLayout>
        <UDashboardPanel id="mainleftsidebar" :width="225" :resizable="{ min: 175, max: 275 }" collapsible>
            <UDashboardNavbar class="!border-transparent" :ui="{ left: 'flex-1' }">
                <template #left>
                    <TopLogoDropdown />
                </template>
            </UDashboardNavbar>

            <UDashboardSidebar>
                <template #header>
                    <UDashboardSearchButton :label="$t('common.search_field')" />
                </template>

                <DashboardSidebarLinks :links="links" />

                <template v-if="clipboardLink.length > 0">
                    <UDivider />

                    <DashboardSidebarLinks :links="clipboardLink" />
                </template>

                <template v-if="quickAccessButtons">
                    <UDivider />

                    <DashboardSidebarLinks :links="quickAccessButtons" />
                </template>

                <div class="flex-1" />

                <DashboardSidebarLinks :links="footerLinks" />

                <UDivider class="sticky bottom-0" />

                <template #footer>
                    <UserDropdown />
                </template>
            </UDashboardSidebar>
        </UDashboardPanel>

        <slot />

        <ClientOnly>
            <!-- Events -->
            <LazyPartialsEventsLayer />
        </ClientOnly>

        <HelpSlideover />
        <NotificationsSlideover />

        <ClientOnly>
            <LazyUDashboardSearch
                v-if="activeChar"
                :empty-state="{
                    icon: 'i-mdi-globe-model',
                    label: $t('commandpalette.empty.title'),
                    queryLabel: $t('commandpalette.empty.title'),
                }"
                :placeholder="`${$t('common.search_field')} (${$t('commandpalette.footer', { key1: '@', key2: '#' })})`"
                :groups="groups as Group[]"
            />
        </ClientOnly>
    </UDashboardLayout>
</template>
