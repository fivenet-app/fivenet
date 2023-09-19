<script lang="ts" setup>
import AccountInfo from '~/components/auth/account/AccountInfo.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import { useNotificatorStore } from '~/store/notificator';

useHead({
    title: 'pages.auth.account_info.title',
});
definePageMeta({
    title: 'pages.auth.account_info.title',
    requiresAuth: true,
    authOnlyToken: true,
    showQuickButtons: false,
    showCookieOptions: true,
});

const notifications = useNotificatorStore();
const route = useRoute();

// `oauth2Connect` can be `failed` (with `reason`) or `success`
const query = route.query;
if (query.oauth2Connect) {
    if (query.oauth2Connect === 'success') {
        notifications.dispatchNotification({
            title: { key: 'notifications.auth.oauth2_connect.success.title', parameters: [] },
            content: { key: 'notifications.auth.oauth2_connect.success.content', parameters: [] },
            type: 'info',
        });
    } else if (query.oauth2Connect === 'failed') {
        const reason = query.reason ?? 'N/A';

        notifications.dispatchNotification({
            title: { key: 'notifications.auth.oauth2_connect.failed.title', parameters: [] },
            content: { key: 'notifications.auth.oauth2_connect.failed.content', parameters: [reason.toString()] },
            type: 'error',
        });
    }
}
</script>

<template>
    <ContentWrapper>
        <AccountInfo />
    </ContentWrapper>
</template>
