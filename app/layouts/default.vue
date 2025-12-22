<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';
import SuperuserJobToggle from '~/components/partials/SuperuserJobToggle.vue';
import MathCalculatorDrawer from '~/components/quickbuttons/mathcalculator/MathCalculatorDrawer.vue';
import PenaltyCalculatorDrawer from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorDrawer.vue';
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import UserMenu from '~/components/UserMenu.vue';
import { useMailerStore } from '~/stores/mailer';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can, activeChar, jobProps, isSuperuser } = useAuth();

const { isDashboardSidebarSlideoverOpen, isHelpSlideoverOpen } = useDashboard();

const overlay = useOverlay();

const { website } = useAppConfig();

const mailerStore = useMailerStore();
const { unreadCount } = storeToRefs(mailerStore);

const route = useRoute();

const open = ref(false);

const links = computed<NavigationMenuItem[]>(() =>
    [
        {
            label: t('common.overview'),
            icon: 'i-mdi-home-outline',
            to: '/overview',
            tooltip: {
                text: t('common.overview'),
                kbds: ['G', 'H'],
            },
        },
        {
            label: t('common.mail'),
            icon: 'i-mdi-inbox-full-outline',
            to: '/mail',
            badge: unreadCount.value > 0 ? (unreadCount.value <= 9 ? unreadCount.value.toString() : '9+') : undefined,
            tooltip: {
                text: t('common.mail'),
                kbds: ['G', 'E'],
            },
            permission: 'mailer.MailerService/ListEmails' as Perms,
            active: route.name.startsWith('mail'),
        },
        {
            label: t('common.citizen', 1),
            icon: 'i-mdi-account-multiple-outline',
            to: '/citizens',
            tooltip: {
                text: t('common.citizen', 1),
                kbds: ['G', 'C'],
            },
            permission: 'citizens.CitizensService/ListCitizens' as Perms,
            active: route.name.startsWith('citizens'),
        },
        {
            label: t('common.vehicle', 2),
            icon: 'i-mdi-car-outline',
            to: '/vehicles',
            tooltip: {
                text: t('common.vehicle', 2),
                kbds: ['G', 'V'],
            },
            permission: 'vehicles.VehiclesService/ListVehicles' as Perms,
        },
        {
            label: t('common.document', 2),
            icon: 'i-mdi-file-document-box-multiple-outline',
            to: '/documents',
            tooltip: {
                text: t('common.document', 2),
                kbds: ['G', 'D'],
            },
            defaultOpen: false,
            children: [
                {
                    label: t('common.approvals', 2),
                    icon: 'i-mdi-approval',
                    to: '/documents/approvals',
                },
            ],
            permission: 'documents.DocumentsService/ListDocuments' as Perms,
            active: route.name.startsWith('documents'),
        },
        {
            label: t('common.job'),
            icon: 'i-mdi-briefcase-outline',
            to: '/jobs/overview',
            tooltip: {
                text: t('common.job'),
                kbds: ['G', 'J'],
            },
            defaultOpen: false,
            children: [
                {
                    label: t('common.overview'),
                    icon: 'i-mdi-briefcase-outline',
                    to: '/jobs/overview',
                },
                {
                    label: t('common.colleague', 2),
                    icon: 'i-mdi-account-group',
                    to: '/jobs/colleagues',
                    permission: 'jobs.JobsService/ListColleagues' as Perms,
                },
                {
                    label: t('common.activity'),
                    icon: 'i-mdi-pulse',
                    to: '/jobs/activity',
                    permission: 'jobs.JobsService/ListColleagueActivity' as Perms,
                },
                {
                    label: t('common.timeclock'),
                    icon: 'i-mdi-timeline-clock',
                    to: '/jobs/timeclock',
                    permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
                },
                {
                    label: t('common.conduct_register', 2),
                    icon: 'i-mdi-list-status',
                    to: '/jobs/conduct',
                    permission: 'jobs.ConductService/ListConductEntries' as Perms,
                },
            ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
            permission: 'jobs.JobsService/ListColleagues' as Perms,
            active: route.name.startsWith('jobs'),
        },
        {
            label: t('common.calendar'),
            icon: 'i-mdi-calendar-outline',
            to: '/calendar',
            tooltip: {
                text: t('common.calendar'),
                kbds: ['G', 'K'],
            },
            active: route.name.startsWith('calendar'),
        },
        {
            label: t('common.qualification', 2),
            icon: 'i-mdi-school-outline',
            to: '/qualifications',
            tooltip: {
                text: t('common.qualification', 2),
                kbds: ['G', 'Q'],
            },
            permission: 'qualifications.QualificationsService/ListQualifications' as Perms,
            active: route.name.startsWith('qualifications'),
        },
        {
            label: t('common.livemap'),
            icon: 'i-mdi-map-outline',
            to: '/livemap',
            tooltip: {
                text: t('common.livemap'),
                kbds: ['G', 'M'],
            },
            permission: 'livemap.LivemapService/Stream' as Perms,
        },
        {
            label: t('common.dispatch_center'),
            icon: 'i-mdi-car-emergency',
            to: '/centrum',
            tooltip: {
                text: t('common.dispatch_center'),
                kbds: ['G', 'W'],
            },
            permission: 'centrum.CentrumService/TakeControl' as Perms,
            active: route.name.startsWith('centrum'),
        },
        {
            label: t('common.wiki'),
            icon: 'i-mdi-brain',
            to: '/wiki',
            tooltip: {
                text: t('common.wiki'),
                kbds: ['G', 'L'],
            },
            permission: 'wiki.WikiService/ListPages' as Perms,
            active: route.name.startsWith('wiki'),
        },
        {
            label: t('common.control_panel'),
            icon: 'i-mdi-cog-outline',
            to: '/settings',
            tooltip: {
                text: t('common.control_panel'),
                kbds: ['G', 'P'],
            },
            defaultOpen: false,
            children: [
                {
                    label: t('components.settings.job_props.job_properties'),
                    icon: 'i-mdi-tune',
                    to: '/settings/props',
                    permission: 'settings.SettingsService/SetJobProps' as Perms,
                },
                {
                    label: t('common.role', 2),
                    icon: 'i-mdi-account-group',
                    to: '/settings/roles',
                    permission: 'settings.SettingsService/GetRoles' as Perms,
                },
            ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
            permission: 'settings.SettingsService/GetJobProps' as Perms,
            active: route.name.startsWith('settings'),
        },
    ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
);

const footerLinks = computed(() =>
    [
        website.statsPage
            ? {
                  label: t('pages.stats.title'),
                  icon: 'i-mdi-analytics',
                  to: '/stats',
              }
            : undefined,
        {
            label: t('common.help'),
            icon: 'i-mdi-question-mark-circle-outline',
            tooltip: {
                kbds: ['?'],
            },
            onClick: () => (isHelpSlideoverOpen.value = true),
        },
        {
            label: t('common.about'),
            icon: 'i-mdi-about-circle-outline',
            to: '/about',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const clipboardModal = overlay.create(ClipboardModal);

const clipboardLink = computed(() =>
    [
        activeChar.value &&
        can([
            'documents.DocumentsService/UpdateDocument',
            'citizens.CitizensService/GetUser',
            'vehicles.VehiclesService/ListVehicles',
        ]).value
            ? {
                  label: t('common.clipboard'),
                  icon: 'i-mdi-clipboard-list-outline',
                  onClick: () => clipboardModal.open(),
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const penaltyCalculatorDrawer = overlay.create(PenaltyCalculatorDrawer);
const mathCalculatorDrawer = overlay.create(MathCalculatorDrawer);

const quickAccessButtons = computed(() =>
    [
        jobProps.value?.quickButtons?.penaltyCalculator || isSuperuser.value
            ? {
                  label: t('components.penaltycalculator.title'),
                  icon: 'i-mdi-gavel',
                  onClick: () => {
                      isDashboardSidebarSlideoverOpen.value = false;
                      penaltyCalculatorDrawer.open();
                  },
              }
            : undefined,
        jobProps.value?.quickButtons?.mathCalculator || isSuperuser.value
            ? {
                  label: t('components.mathcalculator.title'),
                  icon: 'i-mdi-calculator',
                  onClick: () => {
                      isDashboardSidebarSlideoverOpen.value = false;
                      mathCalculatorDrawer.open();
                  },
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts(extractShortcuts(links.value));
</script>

<template>
    <UDashboardGroup unit="%">
        <UDashboardSidebar
            id="default"
            v-model:open="open"
            :default-size="15"
            :min-size="10"
            :max-size="25"
            collapsible
            resizable
            class="bg-elevated/25"
            :ui="{ footer: 'lg:border-t lg:border-default' }"
        >
            <template #header="{ collapsed }">
                <TopLogoDropdown :collapsed="collapsed" />
            </template>

            <template #default="{ collapsed }">
                <UDashboardSearchButton :collapsed="collapsed" :label="$t('common.search_field')" />

                <UNavigationMenu orientation="vertical" tooltip popover :items="links" :collapsed="collapsed" />

                <template v-if="clipboardLink.length > 0">
                    <USeparator />

                    <UNavigationMenu orientation="vertical" tooltip popover :items="clipboardLink" :collapsed="collapsed" />
                </template>

                <template v-if="quickAccessButtons">
                    <USeparator />

                    <UNavigationMenu
                        orientation="vertical"
                        tooltip
                        popover
                        :items="quickAccessButtons"
                        :collapsed="collapsed"
                    />
                </template>

                <div class="flex-1" />

                <template v-if="can(['Superuser/CanBeSuperuser', 'Superuser/Superuser']).value">
                    <SuperuserJobToggle :collapsed="collapsed" />

                    <USeparator />
                </template>

                <UNavigationMenu orientation="vertical" tooltip popover :items="footerLinks" :collapsed="collapsed" />
            </template>

            <template #footer="{ collapsed }">
                <UserMenu :collapsed="collapsed" />
            </template>
        </UDashboardSidebar>

        <slot />

        <ClientOnly>
            <LazyPartialsCommandSearch v-if="activeChar" :links="links" />

            <LazyPartialsWebSocketStatusOverlay />

            <LazyPartialsEventsLayer />
        </ClientOnly>
    </UDashboardGroup>
</template>
