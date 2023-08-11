<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { AccountIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';
import ChangePasswordModal from './ChangePasswordModal.vue';
import DebugInfo from './DebugInfo.vue';
import OAuth2Connections from './OAuth2Connections.vue';

const { $grpc } = useNuxtApp();

const { data: account, pending, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse | undefined> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getAuthClient().getAccountInfo({});

            return res(await call.response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const changePasswordModal = ref(false);

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);

        await refresh();
    }
}
</script>

<template>
    <div class="py-2 mt-5 max-w-5xl mx-auto">
        <ChangePasswordModal :open="changePasswordModal" @close="changePasswordModal = false" />
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" :icon="AccountIcon" />
        <div v-else>
            <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral">
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
                                {{ $t('components.auth.account_info.change_password') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    @click="changePasswordModal = true"
                                    class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                >
                                    {{ $t('components.auth.account_info.change_password_button') }}
                                </button>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>

            <OAuth2Connections
                v-if="account"
                @click="removeOAuth2Connection($event)"
                :providers="account.oauth2Providers"
                :connections="account.oauth2Connections"
            />

            <DebugInfo />
        </div>
    </div>
</template>
