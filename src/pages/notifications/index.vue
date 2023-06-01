<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import NotificationsList from '~/components/partials/notification/NotificationsList.vue';

useHead({
    title: 'common.notification',
});
definePageMeta({
    title: 'pages.notifications',
    requiresAuth: true,
    showQuickButtons: false,
});
</script>

<template>
    <ContentWrapper>
        <NotificationsList />
    </ContentWrapper>
</template>
