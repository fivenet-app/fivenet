<script lang="ts" setup>
import AttrRolesList from '~/components/rector/attrs/AttrRolesList.vue';

useHead({
    title: 'pages.rector.limiter.title',
});
definePageMeta({
    title: 'pages.rector.limiter.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>

<template>
    <AttrRolesList :all="true" />
</template>
