<script lang="ts" setup>
import { AccountGroupIcon, CarSpeedLimiterIcon, GroupIcon, MathLogIcon, ScaleBalanceIcon, TuneIcon } from 'mdi-vue3';
import Cards from '~/components/partials/Cards.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { CardElements } from '~/utils/types';

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
        title: t('common.propertie', 2),
        description: t('pages.rector.features.properties'),
        href: { name: 'rector-props' },
        permission: 'RectorService.GetJobProps',
        icon: markRaw(TuneIcon),
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-100',
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
        title: t('common.units'),
        description: t('pages.rector.features.units'),
        href: { name: 'rector-units' },
        permission: 'CentrumService.ListUnits',
        icon: markRaw(GroupIcon),
        iconForeground: 'text-blue-900',
        iconBackground: 'bg-blue-100',
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
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>
