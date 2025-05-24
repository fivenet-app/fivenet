<script lang="ts" setup>
import CardsList from '~/components/partials/CardsList.vue';
import type { CardElements } from '~/utils/types';

const { t } = useI18n();

useHead({
    title: 'common.control_panel',
});

definePageMeta({
    title: 'common.control_panel',
    requiresAuth: true,
    permission: 'settings.SettingsService.GetRoles',
});

const items = [
    {
        title: t('components.settings.job_props.job_properties'),
        description: t('pages.settings.features.properties'),
        to: { name: 'settings-props' },
        permission: 'settings.SettingsService.GetJobProps',
        icon: 'i-mdi-tune',
    },
    {
        title: t('common.role', 2),
        description: t('components.settings.role_view.add_permission'),
        to: { name: 'settings-roles' },
        permission: 'settings.SettingsService.GetRoles',
        icon: 'i-mdi-account-group',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.settings.features.audit_log'),
        to: { name: 'settings-audit' },
        permission: 'settings.SettingsService.ViewAuditLog',
        icon: 'i-mdi-math-log',
    },
    {
        title: t('pages.settings.limiter.title'),
        description: t('pages.settings.features.limiter'),
        to: { name: 'settings-limiter' },
        permission: 'Superuser',
        icon: 'i-mdi-car-speed-limiter',
    },
    {
        title: t('pages.settings.laws.title'),
        description: t('pages.settings.features.laws'),
        to: { name: 'settings-laws' },
        permission: 'Superuser',
        icon: 'i-mdi-scale-balance',
    },
    {
        title: t('pages.settings.filestore.title'),
        description: t('pages.settings.features.filestore'),
        to: { name: 'settings-filestore' },
        permission: 'Superuser',
        icon: 'i-mdi-file-multiple',
    },
    {
        title: t('pages.settings.settings.title'),
        description: t('pages.settings.features.settings'),
        to: { name: 'settings-settings' },
        permission: 'Superuser',
        icon: 'i-mdi-office-building-cog',
    },
    {
        title: t('pages.settings.accounts.title'),
        description: t('pages.settings.features.accounts'),
        to: { name: 'settings-accounts' },
        permission: 'Superuser',
        icon: 'i-mdi-account-multiple',
    },
    {
        title: t('pages.settings.cron.title'),
        description: t('pages.settings.features.cron'),
        to: { name: 'settings-cron' },
        permission: 'Superuser',
        icon: 'i-mdi-calendar-task',
    },
] as CardElements;
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.control_panel')" />

            <UDashboardPanelContent>
                <CardsList :items="items" />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
