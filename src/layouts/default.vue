<script lang="ts" setup>
import type { DashboardSidebarLink } from '@nuxt/ui-pro/types';
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import CommandPalette from '~/components/partials/CommandPalette.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import QuickButtons from '~/components/partials/quickbuttons/QuickButtons.vue';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

// Use client date to show any event overlays
const now = new Date();
const showSnowflakes = now.getMonth() + 1 === 12 && now.getDate() >= 21 && now.getDate() <= 26;

const { t } = useI18n();

const route = useRoute();

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
        label: 'Help & Support',
        icon: 'i-heroicons-question-mark-circle',
        click: () => (isHelpSlideoverOpen.value = true),
    },
];

const groups = [
    {
        key: 'links',
        label: 'Go to',
        commands: links.map((link) => ({ ...link, shortcuts: link.tooltip?.shortcuts })),
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
</script>

<template>
    <UDashboardLayout>
        <UDashboardPanel :width="250" :resizable="{ min: 200, max: 300 }" collapsible>
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

        <NotificationProvider />

        <!-- Events -->
        <LazyPartialsEventsSnowflakesContainer v-if="showSnowflakes" />

        <div class="w-full overflow-y-auto">
            <slot />

            <QuickButtons v-if="activeChar && (route.meta.showQuickButtons === undefined || route.meta.showQuickButtons)" />
        </div>

        <!-- ~/components/HelpSlideover.vue -->
        <HelpSlideover />
        <!-- ~/components/NotificationsSlideover.vue -->
        <NotificationsSlideover />

        <CommandPalette v-if="activeChar" />

        <ClientOnly>
            <LazyUDashboardSearch :groups="groups" />
        </ClientOnly>
    </UDashboardLayout>
</template>
