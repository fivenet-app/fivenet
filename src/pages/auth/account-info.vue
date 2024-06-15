<script lang="ts" setup>
import AccountInfo from '~/components/auth/account/AccountInfo.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

useHead({
    title: 'components.auth.AccountInfo.title',
});
definePageMeta({
    title: 'components.auth.AccountInfo.title',
    requiresAuth: true,
    authOnlyToken: true,
    showCookieOptions: true,
});

const notifications = useNotificatorStore();

// `oauth2Connect` can be `failed` (with `reason`) or `success`
const oauth2Connect = useRouteQuery('oauth2Connect');
if (oauth2Connect.value) {
    const status = oauth2Connect.value;
    if (status === 'success') {
        notifications.add({
            title: { key: 'notifications.auth.oauth2_connect.success.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_connect.success.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } else {
        const reason = useRouteQuery('reason', 'N/A');

        notifications.add({
            title: { key: 'notifications.auth.oauth2_connect.failed.title', parameters: {} },
            description: { key: 'notifications.auth.oauth2_connect.failed.content', parameters: { msg: reason.toString() } },
            type: NotificationType.ERROR,
        });
    }
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.auth.AccountInfo.title')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/overview">
                        {{ $t('common.back') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <AccountInfo />
        </UDashboardPanel>
    </UDashboardPage>
</template>
