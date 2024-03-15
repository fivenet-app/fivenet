<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationsResultsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsResultsListEntry from '~/components/jobs/qualifications/QualificationsResultsListEntry.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`qualifications-own`, () => listQualificationsResults());

async function listQualificationsResults(): Promise<ListQualificationsResultsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualificationsResults({
            pagination: {
                offset: offset.value,
            },
            status: [ResultStatus.SUCCESSFUL],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="border-b border-gray-200 bg-base-600 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                <div class="ml-4 mt-4">
                    <h3 class="text-base font-semibold leading-6 text-gray-200">
                        {{ $t('components.qualifications.user_qualifications') }}
                    </h3>
                    <p class="mt-1 text-sm text-gray-500"></p>
                </div>
            </div>
        </div>
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data?.results.length === 0" />

            <template v-else>
                <ul role="list" class="divide-y divide-gray-100">
                    <QualificationsResultsListEntry v-for="result in data?.results" :key="result.id" :qualification="result" />
                </ul>
            </template>
        </div>
        <div class="border-t border-gray-200 bg-base-600 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex items-center">
                <TablePagination
                    class="w-full"
                    :pagination="data?.pagination"
                    :show-border="false"
                    :refresh="refresh"
                    @offset-change="offset = $event"
                />
            </div>
        </div>
    </div>
</template>
