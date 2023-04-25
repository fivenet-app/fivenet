<script lang="ts" setup>
import { AdjustmentsVerticalIcon, KeyIcon, QueueListIcon } from '@heroicons/vue/24/outline';
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
});

const features = [
    {
        title: t('common.propertie', 2),
        description: t('pages.rector.features.properties'),
        href: { name: 'rector-props' },
        permission: 'RectorService.GetJobProps',
        icon: AdjustmentsVerticalIcon,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: t('common.role', 2),
        description: t('components.rector.role_view.add_permission'),
        href: { name: 'rector-roles' },
        permission: 'RectorService.GetRoles',
        icon: KeyIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: t('common.audit_log', 1),
        description: t('pages.rector.features.audit_log'),
        href: { name: 'rector-audit' },
        permission: 'RectorService.ViewAuditLog',
        icon: QueueListIcon,
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
