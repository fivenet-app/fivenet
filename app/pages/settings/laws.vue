<script lang="ts" setup>
import LawList from '~/components/settings/laws/LawList.vue';

useHead({
    title: 'pages.settings.laws.title',
});

definePageMeta({
    title: 'pages.settings.laws.title',
    requiresAuth: true,
    permission: 'settings.LawsService/CreateOrUpdateLawBook',
});
</script>

<template>
    <UDashboardPanel>
        <LawList />
    </UDashboardPanel>
</template>
