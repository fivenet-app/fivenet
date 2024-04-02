<script lang="ts" setup>
import type { DashboardSidebarLink } from '@nuxt/ui-pro/types';
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import BodyCheckupModal from '~/components/quickbuttons/bodycheckup/BodyCheckupModal.vue';
import PenaltyCalculatorModal from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorModal.vue';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';

const authStore = useAuthStore();
const { activeChar, jobProps } = storeToRefs(authStore);

// Use client date to show any event overlays
const now = new Date();
const showSnowflakes = now.getMonth() + 1 === 12 && now.getDate() >= 21 && now.getDate() <= 26;

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const appConfig = useAppConfig();
const { isHelpSlideoverOpen } = useDashboard();

const links = [
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
        label: t('common.citizen'),
        icon: 'i-mdi-account-multiple-outline',
        to: '/citizens',
        badge: '4',
        tooltip: {
            text: t('common.citizen'),
            shortcuts: ['G', 'C'],
        },
        permission: 'CitizenStoreService.ListCitizens',
    },
    {
        label: t('common.vehicle'),
        icon: 'i-mdi-car-outline',
        to: '/vehicles',
        tooltip: {
            text: t('common.vehicle'),
            shortcuts: ['G', 'V'],
        },
        permission: 'DMVService.ListVehicles',
    },
    {
        label: t('common.document'),
        icon: 'i-mdi-file-document-box-multiple-outline',
        to: '/documents',
        tooltip: {
            text: t('common.document'),
            shortcuts: ['G', 'D'],
        },
        permission: 'DocStoreService.ListDocuments',
    },
    {
        label: t('common.job'),
        icon: 'i-mdi-briefcase-outline',
        to: '/jobs/overview',
        tooltip: {
            text: t('common.job'),
            shortcuts: ['G', 'J'],
        },
        permission: 'JobsService.ListColleagues',
    },
    {
        label: t('common.livemap'),
        icon: 'i-mdi-map-outline',
        to: '/livemap',
        tooltip: {
            text: t('common.livemap'),
            shortcuts: ['G', 'M'],
        },
        permission: 'LivemapperService.Stream',
    },
    {
        label: t('common.dispatch_center'),
        icon: 'i-mdi-car-emergency',
        to: '/centrum',
        tooltip: {
            text: t('common.dispatch_center'),
            shortcuts: ['G', 'W'],
        },
        permission: 'CentrumService.TakeControl',
    },
    {
        label: t('common.control_panel'),
        icon: 'i-mdi-cog',
        to: '/rector',
        tooltip: {
            text: t('common.control_panel'),
            shortcuts: ['G', 'P'],
        },
        permission: 'RectorService.GetJobProps',
    },
] as (DashboardSidebarLink & { permission?: Perms | Perms[] })[];

const footerLinks = [
    {
        label: t('common.help'),
        icon: 'i-heroicons-question-mark-circle',
        click: () => (isHelpSlideoverOpen.value = true),
    },
    {
        label: t('common.about'),
        icon: 'i-mdi-about-circle-outline',
        to: '/about',
    },
];

const groups = [
    {
        key: 'links',
        label: t('common.goto'),
        commands: links.map((link) => ({ ...link, shortcuts: link.tooltip?.shortcuts })),
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
                console.log(q);
                if (q && (q.startsWith('@') || q.startsWith('#'))) {
                    return [];
                }

                return defaultCommands.filter((c) => !q || c.label.includes(q));
            }

            const prefix = q.substring(0, q.indexOf('-')).toUpperCase();
            const id = q.substring(q.indexOf('-') + 1).trim();
            console.log(prefix, id, id.length > 0 && isNumber(id));
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
                        const call = $grpc.getDocStoreClient().listDocuments({
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
                        $grpc.handleError(e as RpcError);
                        throw e;
                    }
                }

                case '@':
                default: {
                    try {
                        const call = $grpc.getCitizenStoreClient().listCitizens({
                            pagination: {
                                offset: 0,
                                pageSize: 10,
                            },
                            searchName: query,
                        });
                        const { response } = await call;

                        return response.users.map((u) => ({
                            id: u.userId,
                            label: `${u.firstname} ${u.lastname}`,
                            suffix: u.dateofbirth,
                            to: `/citizens/${u.userId}`,
                        }));
                    } catch (e) {
                        $grpc.handleError(e as RpcError);
                        throw e;
                    }
                }
            }
        },
    },
];

const defaultColors = ref(
    ['green', 'teal', 'cyan', 'sky', 'blue', 'indigo', 'violet'].map((color) => ({
        label: color,
        chip: color,
        click: () => (appConfig.ui.primary = color),
    })),
);
const colors = computed(() => defaultColors.value.map((color) => ({ ...color, active: appConfig.ui.primary === color.label })));

const modal = useModal();

const quickAccessButtons = computed(() =>
    [
        jobProps.value?.quickButtons?.penaltyCalculator
            ? {
                  label: t('components.penaltycalculator.title'),
                  icon: 'i-mdi-calculator',
                  click: () => modal.open(PenaltyCalculatorModal, {}),
              }
            : undefined,
        jobProps.value?.quickButtons?.bodyCheckup
            ? {
                  label: t('components.bodycheckup.title'),
                  icon: 'i-mdi-human',
                  click: () => modal.open(BodyCheckupModal, {}),
              }
            : undefined,
    ].filter((c) => c !== undefined),
);
</script>

<template>
    <UDashboardLayout>
        <UDashboardPanel :width="225" :resizable="{ min: 175, max: 275 }" collapsible>
            <UDashboardNavbar class="!border-transparent" :ui="{ left: 'flex-1' }">
                <template #left>
                    <TopLogoDropdown />
                </template>
            </UDashboardNavbar>

            <UDashboardSidebar>
                <template #header>
                    <UDashboardSearchButton />
                </template>

                <UDashboardSidebarLinks :links="links.filter((l) => l.permission === undefined || can(l.permission))" />

                <UDivider />

                <UDashboardSidebarLinks
                    v-if="quickAccessButtons.length > 0"
                    :links="[{ label: t('components.rector.job_props.quick_buttons'), children: quickAccessButtons }]"
                />

                <UDivider />

                <UDashboardSidebarLinks
                    :links="[{ label: 'Colors', draggable: true, children: colors }]"
                    @update:links="(colors) => (defaultColors = colors)"
                />

                <div class="flex-1" />

                <UDashboardSidebarLinks :links="footerLinks" />

                <UDivider class="sticky bottom-0" />

                <template #footer>
                    <!-- ~/components/UserDropdown.vue -->
                    <UserDropdown />
                </template>
            </UDashboardSidebar>
        </UDashboardPanel>

        <!-- Events -->
        <LazyPartialsEventsSnowflakesContainer v-if="showSnowflakes" />

        <div class="w-full max-w-full overflow-y-auto">
            <slot />
        </div>

        <!-- ~/components/HelpSlideover.vue -->
        <HelpSlideover />
        <!-- ~/components/NotificationsSlideover.vue -->
        <NotificationsSlideover />

        <NotificationProvider />

        <ClientOnly>
            <LazyUDashboardSearch
                v-if="activeChar"
                :groups="groups"
                :empty-state="{
                    icon: 'i-mdi-globe-model',
                    label: t('commandpalette.empty.title'),
                    queryLabel: t('commandpalette.empty.title'),
                }"
                :placeholder="`${$t('common.search')}... (${$t('commandpalette.footer', { key1: '#', key2: '@' })})`"
            />
        </ClientOnly>
    </UDashboardLayout>
</template>
