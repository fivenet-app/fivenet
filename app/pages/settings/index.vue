<script lang="ts" setup>
import CardsList from '~/components/partials/CardsList.vue';
import SystemStatus from '~/components/settings/SystemStatus.vue';
import type { CardElements } from '~/utils/types';

const { t } = useI18n();

useHead({
    title: 'common.control_panel',
});

definePageMeta({
    title: 'common.control_panel',
    requiresAuth: true,
    permission: 'settings.SettingsService/GetRoles',
});

const { isSuperuser } = useAuth();

const items = [
    {
        title: t('components.settings.job_props.job_properties'),
        description: t('pages.settings.features.properties'),
        to: { name: 'settings-props' },
        permission: 'settings.SettingsService/GetJobProps',
        icon: 'i-mdi-tune',
    },
    {
        title: t('common.role', 2),
        description: t('components.settings.role_view.add_permission'),
        to: { name: 'settings-roles' },
        permission: 'settings.SettingsService/GetRoles',
        icon: 'i-mdi-account-group',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.settings.features.audit_log'),
        to: { name: 'settings-audit' },
        permission: 'settings.SettingsService/ViewAuditLog',
        icon: 'i-mdi-math-log',
    },
    {
        title: t('common.dispatch_center_settings'),
        description: t('pages.settings.features.dispatch_center'),
        to: { name: 'centrum-settings' },
        permission: 'centrum.CentrumService/UpdateSettings',
        icon: 'i-mdi-car-emergency',
    },
    {
        title: t('pages.settings.laws.title'),
        description: t('pages.settings.features.laws'),
        to: { name: 'settings-laws' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-scale-balance',
    },
] as CardElements;

const superuserItems = [
    {
        title: t('pages.settings.limiter.title'),
        description: t('pages.settings.features.limiter'),
        to: { name: 'settings-limiter' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-car-speed-limiter',
    },
    {
        title: t('pages.settings.filestore.title'),
        description: t('pages.settings.features.filestore'),
        to: { name: 'settings-filestore' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-file-multiple',
    },
    {
        title: t('pages.settings.accounts.title'),
        description: t('pages.settings.features.accounts'),
        to: { name: 'settings-accounts' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-account-multiple',
    },
    {
        title: t('pages.settings.settings.title'),
        description: t('pages.settings.features.settings'),
        to: { name: 'settings-settings' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-office-building-cog',
    },
    {
        title: t('pages.settings.cron.title'),
        description: t('pages.settings.features.cron'),
        to: { name: 'settings-cron' },
        permission: 'Superuser/Superuser',
        icon: 'i-mdi-calendar-task',
    },
] as CardElements;
</script>

<template>
    <UDashboardPanel id="settings">
        <template #header>
            <UDashboardNavbar :title="$t('common.control_panel')" />
        </template>

        <template #body>
            <div class="flex flex-col gap-1">
                <div class="mb-4">
                    <CardsList :items="items" />
                </div>

                <UCard
                    v-if="isSuperuser"
                    class="mb-4"
                    :title="$t('components.settings.system_status.title')"
                    icon="i-mdi-server"
                >
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-server" class="size-5 text-primary" />
                            <h3 class="text-md font-semibold">{{ $t('components.settings.system_status.title') }}</h3>
                        </div>
                    </template>

                    <template #default>
                        <SystemStatus />
                    </template>
                </UCard>

                <UCard
                    v-if="isSuperuser"
                    class="mb-4"
                    :title="$t('components.settings.system_settings')"
                    icon="i-mdi-administrator"
                >
                    <template #header>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-server" class="size-5 text-primary" />
                            <h3 class="text-md font-semibold">{{ $t('components.settings.system_status.title') }}</h3>
                        </div>
                    </template>

                    <template #default>
                        <CardsList :items="superuserItems" />
                    </template>
                </UCard>
            </div>
        </template>
    </UDashboardPanel>
</template>
