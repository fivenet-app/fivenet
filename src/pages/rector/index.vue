<script lang="ts" setup>
import {
    AccountGroupIcon,
    CarSpeedLimiterIcon,
    FileMultipleIcon,
    MathLogIcon,
    OfficeBuildingCogIcon,
    ScaleBalanceIcon,
    TuneIcon,
} from 'mdi-vue3';
import CardsList from '~/components/partials/CardsList.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { type CardElements } from '~/utils/types';

const { t } = useI18n();

useHead({
    title: 'common.control_panel',
});
definePageMeta({
    title: 'common.control_panel',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
    showQuickButtons: false,
});

const features = [
    {
        title: t('components.rector.job_props.job_properties'),
        description: t('pages.rector.features.properties'),
        href: { name: 'rector-props' },
        permission: 'RectorService.GetJobProps',
        icon: markRaw(TuneIcon),
        iconForeground: 'text-warn-900',
        iconBackground: 'bg-warn-100',
    },
    {
        title: t('common.role', 2),
        description: t('components.rector.role_view.add_permission'),
        href: { name: 'rector-roles' },
        permission: 'RectorService.GetRoles',
        icon: markRaw(AccountGroupIcon),
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-100',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.rector.features.audit_log'),
        href: { name: 'rector-audit' },
        permission: 'RectorService.ViewAuditLog',
        icon: markRaw(MathLogIcon),
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-100',
    },
    {
        title: t('pages.rector.limiter.title'),
        description: t('pages.rector.features.limiter'),
        href: { name: 'rector-limiter' },
        permission: 'SuperUser',
        icon: markRaw(CarSpeedLimiterIcon),
        iconForeground: 'text-warn-900',
        iconBackground: 'bg-warn-100',
    },
    {
        title: t('pages.rector.laws.title'),
        description: t('pages.rector.features.laws'),
        href: { name: 'rector-laws' },
        permission: 'SuperUser',
        icon: markRaw(ScaleBalanceIcon),
        iconForeground: 'text-error-900',
        iconBackground: 'bg-error-100',
    },
    {
        title: t('pages.rector.filestore.title'),
        description: t('pages.rector.features.filestore'),
        href: { name: 'rector-filestore' },
        permission: 'SuperUser',
        icon: markRaw(FileMultipleIcon),
        iconForeground: 'text-blue-900',
        iconBackground: 'bg-blue-100',
    },
    {
        title: t('pages.rector.settings.title'),
        description: t('pages.rector.features.settings'),
        href: { name: 'rector-settings' },
        permission: 'SuperUser',
        icon: markRaw(OfficeBuildingCogIcon),
        iconForeground: 'text-red-900',
        iconBackground: 'bg-red-100',
    },
] as CardElements;
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.control_panel')"> </UDashboardNavbar>

            <CardsList :items="features" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
