<script lang="ts" setup>
import { mdiAccountGroup, mdiCarSpeedLimiter, mdiMathLog, mdiTune } from '@mdi/js';
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
        icon: mdiTune,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: t('common.role', 2),
        description: t('components.rector.role_view.add_permission'),
        href: { name: 'rector-roles' },
        permission: 'RectorService.GetRoles',
        icon: mdiAccountGroup,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.rector.features.audit_log'),
        href: { name: 'rector-audit' },
        permission: 'RectorService.ViewAuditLog',
        icon: mdiMathLog,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-50',
    },
    {
        title: t('pages.rector.limiter.title'),
        description: t('pages.rector.features.limiter'),
        href: { name: 'rector-limiter' },
        permission: 'SuperUser',
        icon: mdiCarSpeedLimiter,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-50',
    },
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>
