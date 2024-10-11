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
import { useSettingsStore } from '~/store/settings';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

const { t } = useI18n();

const modal = useModal();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const { data: account, pending: loading, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse> {
    try {
        const call = getGRPCAuthClient().getAccountInfo({});
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

const items = [
    {
        slot: 'accountInfo',
        label: t('components.auth.AccountInfo.title'),
        icon: 'i-mdi-information-slab-circle',
    },
    {
        slot: 'oauth2Connections',
        label: t('components.auth.OAuth2Connections.title'),
        icon: 'i-simple-icons-discord',
    },
    { slot: 'debugInfo', label: t('components.debug_info.title'), icon: 'i-mdi-connection' },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});
</script>

<template>
    <div>
        <template v-if="streamerMode">
            <UDashboardPanelContent class="pb-24">
                <StreamerModeAlert />
            </UDashboardPanelContent>
        </template>
        <template v-else>
            <DataPendingBlock
                v-if="loading"
                :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" icon="i-mdi-account" />

            <template v-else>
                <UTabs v-model="selectedTab" :items="items" class="w-full" :ui="{ list: { rounded: '' } }">
                    <template #accountInfo>
                        <UDashboardPanelContent>
                            <UDashboardSection
                                :title="$t('components.auth.AccountInfo.title')"
                                :description="$t('components.auth.AccountInfo.subtitle')"
                            >
                                <UFormGroup
                                    name="username"
                                    :label="$t('common.username')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <div class="inline-flex w-full justify-between">
                                        <span>
                                            {{ account.account?.username }}
                                        </span>
                                        <CopyToClipboardButton
                                            v-if="account.account?.username"
                                            :value="account.account?.username"
                                        />
                                    </div>
                                </UFormGroup>

                                <UFormGroup
                                    name="license"
                                    :label="$t('components.auth.AccountInfo.license')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <div class="inline-flex w-full justify-between">
                                        <span>
                                            {{ account.account?.license }}
                                        </span>
                                        <CopyToClipboardButton
                                            v-if="account.account?.license"
                                            :value="account.account?.license"
                                        />
                                    </div>
                                </UFormGroup>

                                <UFormGroup
                                    name="change_username"
                                    :label="$t('components.auth.AccountInfo.change_username')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UButton @click="modal.open(ChangeUsernameModal, {})">
                                        {{ $t('components.auth.AccountInfo.change_username_button') }}
                                    </UButton>
                                </UFormGroup>

                                <UFormGroup
                                    name="change_password"
                                    :label="$t('components.auth.AccountInfo.change_password')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UButton @click="modal.open(ChangePasswordModal, {})">
                                        {{ $t('components.auth.AccountInfo.change_password_button') }}
                                    </UButton>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #oauth2Connections>
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
            </template>
        </template>
    </div>
</template>
