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
    authTokenOnly: true,
    showCookieOptions: true,
});

const { t } = useI18n();

const notifications = useNotificatorStore();

// `oauth2Connect` can be `failed` (with `reason`) or `success`
const oauth2Connect = useRouteQuery('oauth2Connect');

onBeforeMount(async () => {
    if (oauth2Connect.value && typeof oauth2Connect.value === 'string') {
        const status = oauth2Connect.value.toLowerCase();

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
                description: {
                    key: 'notifications.auth.oauth2_connect.failed.content',
                    parameters: {
                        msg: t(`notifications.auth.oauth2_connect.failed.reasons.${reason.value.toString()}`, {
                            reason: reason.value.toString(),
                        }),
                    },
                },
                type: NotificationType.ERROR,
            });
        }
    }
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.auth.AccountInfo.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/overview" />
                </template>
            </UDashboardNavbar>

            <AccountInfo />
        </UDashboardPanel>
    </UDashboardPage>
</template>
