<script lang="ts" setup>
import { KeyIcon } from '@heroicons/vue/24/outline';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import Cards from '~~/src/components/partials/Cards.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { CardElements } from '~~/src/utils/types';

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
] as CardElements;
</script>

<template>
    <ContentCenterWrapper>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>
