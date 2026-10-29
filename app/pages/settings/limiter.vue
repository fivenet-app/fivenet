<script lang="ts" setup>
import AttrJobList from '~/components/settings/attrs/AttrJobList.vue';

useHead({
    title: 'pages.settings.limiter.title',
});

definePageMeta({
    title: 'pages.settings.limiter.title',
    requiresAuth: true,
    permission: 'internal.Superuser/CanBeSuperuser',
});
</script>

<template>
    <AttrJobList />
</template>
