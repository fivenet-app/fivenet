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
    permission: 'RectorService.GetRoles',
});

const items = [
    {
        title: t('components.rector.job_props.job_properties'),
        description: t('pages.rector.features.properties'),
        to: { name: 'rector-props' },
        permission: 'RectorService.GetJobProps',
        icon: 'i-mdi-tune',
    },
    {
        title: t('common.role', 2),
        description: t('components.rector.role_view.add_permission'),
        to: { name: 'rector-roles' },
        permission: 'RectorService.GetRoles',
        icon: 'i-mdi-account-group',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.rector.features.audit_log'),
        to: { name: 'rector-audit' },
        permission: 'RectorService.ViewAuditLog',
        icon: 'i-mdi-math-log',
    },
    {
        title: t('pages.rector.limiter.title'),
        description: t('pages.rector.features.limiter'),
        to: { name: 'rector-limiter' },
        permission: 'SuperUser',
        icon: 'i-mdi-car-speed-limiter',
    },
    {
        title: t('pages.rector.laws.title'),
        description: t('pages.rector.features.laws'),
        to: { name: 'rector-laws' },
        permission: 'SuperUser',
        icon: 'i-mdi-scale-balance',
    },
    {
        title: t('pages.rector.filestore.title'),
        description: t('pages.rector.features.filestore'),
        to: { name: 'rector-filestore' },
        permission: 'SuperUser',
        icon: 'i-mdi-file-multiple',
    },
    {
        title: t('pages.rector.settings.title'),
        description: t('pages.rector.features.settings'),
        to: { name: 'rector-settings' },
        permission: 'SuperUser',
        icon: 'i-mdi-office-building-cog',
    },
] as CardElements;
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.control_panel')" />

            <CardsList class="m-2 pb-2" :items="items" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
