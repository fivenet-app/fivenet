<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import ChangeUsernameModal from '~/components/auth/account/ChangeUsernameModal.vue';
import DebugInfo from '~/components/auth/account/DebugInfo.vue';
import OAuth2Connections from '~/components/auth/account/OAuth2Connections.vue';
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

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);

        await refresh();
    }
}

const modal = useModal();
</script>

<template>
    <div>
        <template v-if="streamerMode">
            <UDashboardPanelContent class="pb-2">
                <UDashboardSection
                    :title="$t('system.streamer_mode.title')"
                    :description="$t('system.streamer_mode.description')"
                />
            </UDashboardPanelContent>
        </template>
        <template v-else>
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
                icon="i-mdi-account"
            />

            <template v-else>
                <UDashboardPanelContent class="pb-2">
                    <UDashboardSection
                        :title="$t('components.auth.AccountInfo.title')"
                        :description="$t('components.auth.AccountInfo.subtitle')"
                    >
                        <UFormGroup
                            name="version"
                            :label="$t('common.username')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            {{ account.account?.username }}
                        </UFormGroup>

                        <UFormGroup
                            name="version"
                            :label="$t('components.auth.AccountInfo.license')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            {{ account.account?.license }}
                        </UFormGroup>

                        <UFormGroup
                            name="version"
                            :label="$t('components.auth.AccountInfo.change_username')"
                            class="grid grid-cols-2 items-center gap-2"
                            :ui="{ container: '' }"
                        >
                            <UButton @click="modal.open(ChangeUsernameModal, {})">
                                {{ $t('components.auth.AccountInfo.change_username_button') }}
                            </UButton>
                        </UFormGroup>

                        <UFormGroup
                            name="version"
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
