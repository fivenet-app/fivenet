<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

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

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const authAuthClient = await getAuthAuthClient();

const { data: account, status, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse> {
    try {
        const call = authAuthClient.getAccountInfo({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const items = computed<NavigationMenuItem[]>(() =>
    [
        {
            label: t('components.auth.AccountInfo.title'),
            icon: 'i-mdi-information-outline',
            to: '/auth/account-info',
            exact: true,
        },
        account.value?.oauth2Providers && account.value.oauth2Providers.length > 0
            ? {
                  label: t('components.auth.OAuth2Connections.title'),
                  icon: 'i-simple-icons-discord',
                  to: '/auth/account-info/oauth2',
              }
            : undefined,
        {
            label: t('components.debug_info.title'),
            icon: 'i-mdi-connection',
            to: '/auth/account-info/debug',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

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
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('components.auth.AccountInfo.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/overview" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UNavigationMenu :items="items" highlight class="-mx-1 flex-1" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />

            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :error="error"
                :retry="refresh"
            />

            <DataPendingBlock
                v-else-if="isRequestPending(status)"
                :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])"
            />

            <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" icon="i-mdi-account" />

            <NuxtPage v-else :account="account" />
        </template>
    </UDashboardPanel>
</template>
