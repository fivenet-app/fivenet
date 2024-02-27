<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { AccountIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import ChangeUsernameModal from '~/components/auth/account/ChangeUsernameModal.vue';
import DebugInfo from '~/components/auth/account/DebugInfo.vue';
import OAuth2Connections from '~/components/auth/account/OAuth2Connections.vue';
import GenericContainerPanel from '~/components/partials/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/GenericContainerPanelEntry.vue';
import { useSettingsStore } from '~/store/settings';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const { data: account, pending, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse> {
    try {
        const call = $grpc.getAuthClient().getAccountInfo({});

        return call.response!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const changePasswordModal = ref(false);
const changeUsernameModal = ref(false);

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);

        await refresh();
    }
}
</script>

<template>
    <div class="mx-auto max-w-5xl py-2">
        <template v-if="streamerMode">
            <GenericContainerPanel>
                <template #title>
                    {{ $t('system.streamer_mode.title') }}
                </template>
                <template #description>
                    {{ $t('system.streamer_mode.description') }}
                </template>
            </GenericContainerPanel>
        </template>
        <template v-else>
            <ChangeUsernameModal :open="changeUsernameModal" @close="changeUsernameModal = false" />
            <ChangePasswordModal :open="changePasswordModal" @close="changePasswordModal = false" />

            <DataPendingBlock
                v-if="pending"
                :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="account === null"
                :type="`${$t('common.account')} ${$t('common.data')}`"
                :icon="AccountIcon"
            />

            <template v-else>
                <GenericContainerPanel>
                    <template #title>
                        {{ $t('components.auth.account_info.title') }}
                    </template>
                    <template #description>
                        {{ $t('components.auth.account_info.subtitle') }}
                    </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('common.username') }}
                            </template>
                            <template #default>
                                {{ account.account?.username }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.auth.account_info.license') }}
                            </template>
                            <template #default>
                                {{ account.account?.license }}
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.auth.account_info.change_username') }}
                            </template>
                            <template #default>
                                <button
                                    type="button"
                                    class="rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    @click="changeUsernameModal = true"
                                >
                                    {{ $t('components.auth.account_info.change_username_button') }}
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.auth.account_info.change_password') }}
                            </template>
                            <template #default>
                                <button
                                    type="button"
                                    class="rounded-md bg-primary-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-primary-400"
                                    @click="changePasswordModal = true"
                                >
                                    {{ $t('components.auth.account_info.change_password_button') }}
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>

                <OAuth2Connections
                    :providers="account.oauth2Providers"
                    :connections="account.oauth2Connections"
                    @disconnected="removeOAuth2Connection($event)"
                />

                <DebugInfo />
            </template>
        </template>
    </div>
</template>
