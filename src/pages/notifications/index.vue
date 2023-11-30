<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import NotificationsList from '~/components/partials/notification/NotificationsList.vue';

useHead({
    title: 'components.notifications.title',
});
definePageMeta({
    title: 'components.notifications.title',
    requiresAuth: true,
    showQuickButtons: false,
});
</script>

<template>
    <ContentWrapper>
        <NotificationsList />
    </ContentWrapper>
</template>
