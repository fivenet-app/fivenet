<script lang="ts" setup>
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import ChangeUsernameModal from '~/components/auth/account/ChangeUsernameModal.vue';
import DebugInfo from '~/components/auth/account/DebugInfo.vue';
import OAuth2Connections from '~/components/auth/account/OAuth2Connections.vue';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const { data: account, pending: loading, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse> {
    try {
        const call = $grpc.auth.auth.getAccountInfo({});
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

const items = computed(() =>
    [
        {
            slot: 'accountInfo',
            label: t('components.auth.AccountInfo.title'),
            icon: 'i-mdi-information-outline',
        },
        account.value?.oauth2Providers && account.value.oauth2Providers.length > 0
            ? {
                  slot: 'oauth2Connections',
                  label: t('components.auth.OAuth2Connections.title'),
                  icon: 'i-simple-icons-discord',
              }
            : undefined,
        { slot: 'debugInfo', label: t('components.debug_info.title'), icon: 'i-mdi-connection' },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.value.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items.value[value]?.slot }, hash: '#' });
    },
});
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardPanelContent>
            <StreamerModeAlert />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <UDashboardPanelContent class="p-0 sm:pb-0">
            <DataPendingBlock
                v-if="loading"
                :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" icon="i-mdi-account" />

            <UTabs v-else v-model="selectedTab" class="w-full" :items="items" :ui="{ list: { rounded: '' } }">
                <template #accountInfo>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.auth.AccountInfo.title')"
                            :description="$t('components.auth.AccountInfo.subtitle')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="username"
                                :label="$t('common.username')"
                                :ui="{ container: '' }"
                            >
                                <div class="inline-flex w-full justify-between gap-2">
                                    <span class="truncate">
                                        {{ account.account?.username }}
                                    </span>
                                    <CopyToClipboardButton
                                        v-if="account.account?.username"
                                        :value="account.account?.username"
                                    />
                                </div>
                            </UFormGroup>

                            <UFormGroup
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
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="change_username"
                                :label="$t('components.auth.AccountInfo.change_username')"
                                :ui="{ container: '' }"
                            >
                                <UButton @click="modal.open(ChangeUsernameModal, {})">
                                    {{ $t('components.auth.AccountInfo.change_username_button') }}
                                </UButton>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="change_password"
                                :label="$t('components.auth.AccountInfo.change_password')"
                                :ui="{ container: '' }"
                            >
                                <UButton @click="modal.open(ChangePasswordModal, {})">
                                    {{ $t('components.auth.AccountInfo.change_password_button') }}
                                </UButton>
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>

                <template v-if="account.oauth2Providers.length > 0" #oauth2Connections>
                    <OAuth2Connections
                        :providers="account.oauth2Providers"
                        :connections="account.oauth2Connections"
                        @disconnected="removeOAuth2Connection($event)"
                    />
                </template>

                <template #debugInfo>
                    <DebugInfo />
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </template>
</template>
