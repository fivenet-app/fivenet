<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import ChangeUsernameModal from '~/components/auth/account/ChangeUsernameModal.vue';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';
import DebugInfo from './DebugInfo.vue';
import OAuth2Connections from './OAuth2Connections.vue';

const { t } = useI18n();

const modal = useOverlay();

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

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);
    }
}

const items = computed<TabsItem[]>(() =>
    [
        {
            slot: 'accountInfo' as const,
            label: t('components.auth.AccountInfo.title'),
            icon: 'i-mdi-information-outline',
            value: 'accountInfo',
        },
        account.value?.oauth2Providers && account.value.oauth2Providers.length > 0
            ? {
                  slot: 'oauth2Connections' as const,
                  label: t('components.auth.OAuth2Connections.title'),
                  icon: 'i-simple-icons-discord',
                  value: 'oauth2Connections',
              }
            : undefined,
        {
            slot: 'debugInfo' as const,
            label: t('components.debug_info.title'),
            icon: 'i-mdi-connection',
            value: 'debugInfo',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'accountInfo';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const changeUsernameModal = modal.create(ChangeUsernameModal);
const changePasswordModal = modal.create(ChangePasswordModal);
</script>

<template>
    <DataPendingBlock
        v-if="isRequestPending(status)"
        :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])"
    />

    <UDashboardToolbar v-else>
        <UNavigationMenu
            :items="items.map((i) => ({ ...i, onClick: () => (selectedTab = i.value) }))"
            highlight
            class="-mx-1 flex-1"
        />
    </UDashboardToolbar>

    <UTabs v-model="selectedTab" :items="items">
        <template #accountInfo>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" icon="i-mdi-account" />

            <UPageCard
                v-else
                :title="$t('components.auth.AccountInfo.title')"
                :description="$t('components.auth.AccountInfo.subtitle')"
            >
                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="username"
                    :label="$t('common.username')"
                    :ui="{ container: '' }"
                >
                    <div class="inline-flex w-full justify-between gap-2">
                        <span class="truncate">
                            {{ account.account?.username }}
                        </span>
                        <CopyToClipboardButton v-if="account.account?.username" :value="account.account?.username" />
                    </div>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="license"
                    :label="$t('components.auth.AccountInfo.license')"
                    :ui="{ container: '' }"
                >
                    <div class="inline-flex w-full justify-between gap-2">
                        <span class="truncate">
                            {{ account.account?.license }}
                        </span>

                        <CopyToClipboardButton v-if="account.account?.license" :value="account.account?.license" />
                    </div>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="change_username"
                    :label="$t('components.auth.AccountInfo.change_username')"
                    :ui="{ container: '' }"
                >
                    <UButton @click="changeUsernameModal.open()">
                        {{ $t('components.auth.AccountInfo.change_username_button') }}
                    </UButton>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="change_password"
                    :label="$t('components.auth.AccountInfo.change_password')"
                    :ui="{ container: '' }"
                >
                    <UButton @click="changePasswordModal.open()">
                        {{ $t('components.auth.AccountInfo.change_password_button') }}
                    </UButton>
                </UFormField>
            </UPageCard>
        </template>

        <template v-if="account?.oauth2Providers && account.oauth2Providers.length > 0" #oauth2Connections>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" icon="i-mdi-account" />

            <template v-else>
                <UPageCard
                    :title="$t('components.auth.OAuth2Connections.title')"
                    :description="$t('components.auth.OAuth2Connections.subtitle')"
                >
                    <OAuth2Connections
                        :providers="account.oauth2Providers"
                        :connections="account.oauth2Connections"
                        @disconnected="removeOAuth2Connection($event)"
                    />
                </UPageCard>
            </template>
        </template>

        <template #debugInfo>
            <DebugInfo />
        </template>
    </UTabs>
</template>
