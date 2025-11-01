<script lang="ts" setup>
import List from '~/components/notifications/List.vue';

useHead({
    title: 'components.notifications.title',
});

definePageMeta({
    title: 'components.notifications.title',
    requiresAuth: true,
});
</script>

<template>
    <List />
</template>
