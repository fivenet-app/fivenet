<script lang="ts" setup>
import { Account } from '@fivenet/gen/resources/accounts/accounts_pb';
import { GetAccountInfoRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { UserIcon } from '@heroicons/vue/24/outline';
import ChangePasswordModal from './ChangePasswordModal.vue';

const { $grpc } = useNuxtApp();

const { data: account, pending, refresh, error } = await useLazyAsyncData(`accounmt`, () => getAccountInfo());

async function getAccountInfo(): Promise<Account> {
    return new Promise(async (res, rej) => {
        const req = new GetAccountInfoRequest();

        try {
            const resp = await $grpc.getAuthClient().
                getAccountInfo(req, null);

            if (resp.hasAccount()) {
                return res(resp.getAccount()!);
            }
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const changePasswordModal = ref(false);
</script>

<template>
    <div class="py-2">
        <ChangePasswordModal :open="changePasswordModal" @close="changePasswordModal = false" />
        <DataPendingBlock v-if="pending" message="Loading documents..." />
        <DataErrorBlock v-else-if="error" title="Unable to load account info!" :retry="refresh" />
        <button v-else-if="!account" type="button"
            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <UserIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                No account data found!
            </span>
        </button>
        <div v-else>
            <div class="overflow-hidden bg-gray-400 shadow sm:rounded-lg">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-base font-semibold leading-6 text-white">
                        Account Info
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm text-gray-900">
                        Your FiveNet account details.
                    </p>
                </div>
                <div class="border-t border-gray-800 px-4 py-5 sm:p-0">
                    <dl class="sm:divide-y sm:divide-gray-800">
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium text-gray-900">
                                Username
                            </dt>
                            <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                                {{ account.getUsername() }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium text-gray-900">
                                License
                            </dt>
                            <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                                {{ account.getLicense() }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium text-gray-900">
                                Change Password
                            </dt>
                            <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                                <button type="button" @click="changePasswordModal = true"
                                    class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
                                    Change account Password
                                </button>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>
        </div>
    </div>
</template>
