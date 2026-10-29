<script lang="ts" setup>
import type { Perms } from '~~/gen/ts/perms';
import CardsList from '~/components/partials/CardsList.vue';
import SystemStatus from '~/components/settings/SystemStatus.vue';
import type { CardElement } from '~/utils/types';

const { t } = useI18n();

useHead({
    title: 'common.control_panel',
});

definePageMeta({
    title: 'common.control_panel',
    requiresAuth: true,
    permission: [
        'internal.Superuser/ConfigAdmin',
        'settings.SettingsService/GetJobProps',
        'settings.SettingsService/GetRoles',
        'settings.SettingsService/ViewAuditLog',
        'settings.LawsService/CreateOrUpdateLawBook',
    ],
});

const { can } = useAuth();

const canJobAdmin = can('internal.Superuser/JobAdmin');
const canConfigAdmin = can('internal.Superuser/ConfigAdmin');

const items = computed<CardElement[]>(() => [
    {
        label: t('components.settings.job_props.job_properties'),
        description: t('pages.settings.features.properties'),
        to: '/settings/props',
        permission: 'settings.SettingsService/GetJobProps' as Perms,
        icon: 'i-mdi-tune',
    },
    {
        label: t('common.role', 2),
        description: t('components.settings.role_view.add_permission'),
        to: '/settings/roles',
        permission: 'settings.SettingsService/GetRoles' as Perms,
        icon: 'i-mdi-account-group',
    },
    {
        label: t('common.dispatch_center_settings'),
        description: t('pages.settings.features.dispatch_center'),
        to: '/dispatch/settings',
        permission: 'centrum.CentrumService/UpdateSettings' as Perms,
        icon: 'i-mdi-car-emergency',
    },
    {
        label: t('common.audit_log', 1),
        description: t('pages.settings.features.audit_log'),
        to: '/settings/audit',
        permission: 'settings.SettingsService/ViewAuditLog' as Perms,
        icon: 'i-mdi-math-log',
    },
    {
        label: t('pages.settings.job_assets.title'),
        description: t('pages.settings.job_assets.description'),
        to: '/settings/job-assets',
        permission: 'settings.JobsAssetsService/ListJobAssets' as Perms,
        icon: 'i-mdi-image-multiple-outline',
    },
    {
        label: t('pages.settings.laws.title'),
        description: t('pages.settings.features.laws'),
        to: '/settings/laws',
        permission: 'settings.LawsService/CreateOrUpdateLawBook' as Perms,
        icon: 'i-mdi-scale-balance',
    },
]);

const superuserItems = computed<CardElement[]>(() => [
    {
        label: t('pages.settings.settings.title'),
        description: t('pages.settings.features.settings'),
        to: '/settings/settings',
        permission: 'internal.Superuser/ConfigAdmin' as Perms,
        icon: 'i-mdi-office-building-cog',
    },
    {
        label: t('components.settings.property_definitions.title'),
        description: t('pages.settings.features.property_definitions'),
        to: '/settings/properties',
        permission: 'internal.Superuser/CanBeSuperuser' as Perms,
        icon: 'i-mdi-tune-variant',
    },
    {
        label: t('pages.settings.limiter.title'),
        description: t('pages.settings.features.limiter'),
        to: '/settings/limiter',
        permission: 'internal.Superuser/CanBeSuperuser' as Perms,
        icon: 'i-mdi-car-speed-limiter',
    },
    {
        label: t('pages.settings.jobs.title'),
        description: t('pages.settings.features.jobs'),
        to: '/settings/jobs',
        permission: 'internal.Superuser/CanBeSuperuser' as Perms,
        icon: 'i-mdi-briefcase',
    },
    {
        label: t('pages.settings.accounts.title'),
        description: t('pages.settings.features.accounts'),
        to: '/settings/accounts',
        permission: 'internal.Superuser/ConfigAdmin' as Perms,
        icon: 'i-mdi-account-multiple',
    },
    {
        label: t('pages.settings.filestore.title'),
        description: t('pages.settings.features.filestore'),
        to: '/settings/filestore',
        permission: 'internal.Superuser/ConfigAdmin' as Perms,
        icon: 'i-mdi-file-multiple',
    },
    {
        label: t('pages.settings.cron.title'),
        description: t('pages.settings.features.cron'),
        to: '/settings/cron',
        permission: 'internal.Superuser/ConfigAdmin' as Perms,
        icon: 'i-mdi-calendar-task',
    },
]);
</script>

<template>
    <UDashboardPanel id="settings" :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.control_panel')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="flex flex-col gap-4">
                <div>
                    <div data-tour="control-panel-properties">
                        <CardsList :items="items" />
                    </div>
                </div>

                <SystemStatus v-if="canConfigAdmin" />

                <template v-if="canJobAdmin || canConfigAdmin">
                    <UCard :title="$t('components.settings.system_settings')" icon="i-mdi-administrator" variant="outline">
                        <template #header>
                            <div class="flex items-center gap-2">
                                <UIcon class="size-5 text-primary" name="i-mdi-administrator" />
                                <h3 class="text-md font-semibold">{{ $t('components.settings.system_settings') }}</h3>
                            </div>
                        </template>
                    </UCard>

                    <div data-tour="control-panel-roles">
                        <CardsList :items="superuserItems" />
                    </div>
                </template>
            </div>
        </template>
    </UDashboardPanel>
</template>
