<script lang="ts" setup>
import LawList from '~/components/rector/laws/LawList.vue';

useHead({
    title: 'pages.rector.laws.title',
});

definePageMeta({
    title: 'pages.rector.laws.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <LawList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
