<script lang="ts" setup>
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import CommandPalette from '~/components/partials/CommandPalette.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import { useAuthStore } from '~/store/auth';

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
        id: 'home',
        label: t('common.overview'),
        icon: 'i-mdi-home-outline',
        to: '/overview',
        tooltip: {
            text: t('common.overview'),
            shortcuts: ['G', 'H'],
        },
    },
    {
        id: 'citizens',
        label: t('common.citizen'),
        icon: 'i-mdi-account-multiple-outline',
        to: '/citizens',
        badge: '4',
        tooltip: {
            text: t('common.citizen'),
            shortcuts: ['G', 'C'],
        },
    },
    {
        id: 'vehicles',
        label: t('common.vehicle'),
        icon: 'i-mdi-car-outline',
        to: '/vehicles',
        tooltip: {
            text: t('common.vehicle'),
            shortcuts: ['G', 'V'],
        },
    },
    {
        id: 'documents',
        label: t('common.document'),
        icon: 'i-mdi-file-document-box-multiple-outline',
        to: '/documents',
        tooltip: {
            text: t('common.document'),
            shortcuts: ['G', 'D'],
        },
    },
    {
        id: 'job',
        label: t('common.job'),
        icon: 'i-mdi-briefcase-outline',
        to: '/jobs/overview',
        tooltip: {
            text: t('common.job'),
            shortcuts: ['G', 'J'],
        },
    },
    {
        id: 'livemap',
        label: t('common.livemap'),
        icon: 'i-mdi-map-outline',
        to: '/livemap',
        tooltip: {
            text: t('common.livemap'),
            shortcuts: ['G', 'M'],
        },
    },
    {
        id: 'centrum',
        label: t('common.dispatch_center'),
        icon: 'i-mdi-car-emergency',
        to: '/centrum',
        tooltip: {
            text: t('common.dispatch_center'),
            shortcuts: ['G', 'W'],
        },
    },
    {
        id: 'settings',
        label: 'Settings',
        to: '/settings',
        icon: 'i-heroicons-cog-8-tooth',
        children: [
            {
                label: 'General',
                to: '/settings',
                exact: true,
            },
            {
                label: 'Members',
                to: '/settings/members',
            },
            {
                label: 'Notifications',
                to: '/settings/notifications',
            },
        ],
        tooltip: {
            text: 'Settings',
            shortcuts: ['G', 'S'],
        },
    },
];

const footerLinks = [
    {
        label: 'Invite people',
        icon: 'i-heroicons-plus',
        to: '/settings/members',
    },
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
    {
        key: 'code',
        label: 'Code',
        commands: [
            {
                id: 'source',
                label: 'View page source',
                icon: 'i-simple-icons-github',
                click: () => {
                    window.open(
                        `https://github.com/nuxt-ui-pro/dashboard/blob/main/pages${route.path === '/' ? '/index' : route.path}.vue`,
                        '_blank',
                    );
                },
            },
        ],
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

                <UDashboardSidebarLinks :links="links" />

                <UDivider />

                <!-- <UDashboardSidebarLinks
                    :links="[{ label: 'Colors', draggable: true, children: colors }]"
                    @update:links="(colors) => (defaultColors = colors)"
                /> -->

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
