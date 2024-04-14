<script lang="ts" setup>
import AttrRolesList from '~/components/rector/attrs/AttrRolesList.vue';

useHead({
    title: 'pages.rector.limiter.title',
});
definePageMeta({
    title: 'pages.rector.limiter.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.limiter.title')" />

            <AttrRolesList :all="true" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
