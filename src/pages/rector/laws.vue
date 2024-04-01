<script lang="ts" setup>
import LawList from '~/components/rector/laws/LawList.vue';

useHead({
    title: 'pages.rector.laws.title',
});
definePageMeta({
    title: 'pages.rector.laws.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.laws.title')"> </UDashboardNavbar>

            <LawList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
