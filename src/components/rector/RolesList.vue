<script lang="ts" setup>
import { Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { GetRolesRequest } from '@fivenet/gen/services/rector/rector_pb';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';
import RolesListEntry from './RolesListEntry.vue';

const { $grpc } = useNuxtApp();

const { data: roles, pending, refresh, error } = await useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Array<Role>> {
    return new Promise(async (res, rej) => {
        const req = new GetRolesRequest();
        req.setRank(1);

        try {
            const resp = await $grpc.getRectorClient().
                getRoles(req, null);

            return res(resp.getRolesList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" message="Loading roles..." />
                        <DataErrorBlock v-else-if="error" title="Unable to load roles!" :retry="refresh" />
                        <button v-else-if="roles && roles.length == 0" type="button"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                Use the search field above to search or update your query.
                            </span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Name
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Updated At
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <RolesListEntry v-for="role in roles" :role="role" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Name
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Updated At
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
