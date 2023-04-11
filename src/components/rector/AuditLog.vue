<script lang="ts" setup>
import { PaginationRequest, PaginationResponse } from '@fivenet/gen/resources/common/database/database_pb';
import { AuditEntry } from '@fivenet/gen/resources/rector/audit_pb';
import { ViewAuditLogRequest } from '@fivenet/gen/services/rector/rector_pb';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/solid';
import AuditLogEntry from './AuditLogEntry.vue';
import TablePagination from '../partials/TablePagination.vue';

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(0);

const { data: logs, pending, refresh, error } = await useLazyAsyncData('rector-roles', () => getAuditLog());

async function getAuditLog(): Promise<Array<AuditEntry>> {
    return new Promise(async (res, rej) => {
        const req = new ViewAuditLogRequest();
        req.setPagination((new PaginationRequest()).setOffset(offset.value));

        try {
            const resp = await $grpc.getRectorClient().
                viewAuditLog(req, null);

            pagination.value = resp.getPagination();
            return res(resp.getLogsList());
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
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <!-- TODO -->
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" message="Loading audit logs..." />
                        <DataErrorBlock v-else-if="error" title="Unable to load audit logs!" :retry="refresh" />
                        <button v-else-if="logs && logs.length == 0" type="button"
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
                                        <th scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            ID
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Time
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            User
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Service
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            State
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Data
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <AuditLogEntry v-for="log in logs" :key="log.getId()" :log="log"
                                        class="transition-colors hover:bg-neutral/5" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            ID
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Time
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            User
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Service
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            State
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Data
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
