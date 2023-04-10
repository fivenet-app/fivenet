<script lang="ts" setup>
import { KeyIcon, QueueListIcon } from '@heroicons/vue/24/outline';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import Cards from '~/components/partials/Cards.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { CardElements } from '~/utils/types';

useHead({
    title: 'Control Panel',
});
definePageMeta({
    title: 'Control Panel',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
});

const features = [
    {
        title: 'Roles',
        description:
            'Setup your factions roles job rank permissions.',
        href: { name: 'rector-roles' },
        permission: 'RectorService.GetRoles',
        icon: KeyIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: 'Audit Log',
        description: 'An audit trail of everything your colleagues have done.',
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
