<script lang="ts" setup>
import PropertyDefinitions from '~/components/settings/PropertyDefinitions.vue';

useHead({
    title: 'components.settings.property_definitions.title',
});

definePageMeta({
    title: 'components.settings.property_definitions.title',
    requiresAuth: true,
    permission: 'internal.Superuser/CanBeSuperuser',
});
</script>

<template>
    <PropertyDefinitions />
</template>
