<script lang="ts" setup>
import LawList from '~/components/rector/laws/LawList.vue';

useHead({
    title: 'pages.rector.laws.title',
});

definePageMeta({
    title: 'pages.rector.laws.title',
    requiresAuth: true,
    permission: 'RectorLawsService.CreateOrUpdateLawBook',
});
</script>

<template>

        <UDashboardPanel grow>
            <LawList />
        </UDashboardPanel>

</template>
