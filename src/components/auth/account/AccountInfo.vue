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
import SettingsPanel from '~/components/auth/account/SettingsPanel.vue';

const { $grpc } = useNuxtApp();

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
        <ChangeUsernameModal :open="changeUsernameModal" @close="changeUsernameModal = false" />
        <ChangePasswordModal :open="changePasswordModal" @close="changePasswordModal = false" />

        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])" />
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
            <div class="overflow-hidden bg-base-800 text-neutral shadow sm:rounded-lg">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-base font-semibold leading-6">
                        {{ $t('components.auth.account_info.title') }}
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm">
                        {{ $t('components.auth.account_info.subtitle') }}
                    </p>
                </div>
                <div class="border-t border-base-400 px-4 py-5 sm:p-0">
                    <dl class="sm:divide-y sm:divide-base-400">
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('common.username') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                {{ account.account?.username }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.license') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                {{ account.account?.license }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.change_username') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    class="rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="changeUsernameModal = true"
                                >
                                    {{ $t('components.auth.account_info.change_username_button') }}
                                </button>
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.change_password') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    class="rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="changePasswordModal = true"
                                >
                                    {{ $t('components.auth.account_info.change_password_button') }}
                                </button>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>

            <SettingsPanel />

            <OAuth2Connections
                v-if="account"
                :providers="account.oauth2Providers"
                :connections="account.oauth2Connections"
                @disconnected="removeOAuth2Connection($event)"
            />

            <DebugInfo />
        </template>
    </div>
</template>
