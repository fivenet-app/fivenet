<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationsResultsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsResultsListEntry from '~/components/jobs/qualifications/tutor/QualificationsResultsListEntry.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const props = withDefaults(
    defineProps<{
        qualificationId?: string;
        status?: ResultStatus[];
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
    },
);

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`qualifications-results-${props.qualificationId}`, () =>
    listQualificationsResults(props.qualificationId, props.status),
);

async function listQualificationsResults(
    qualificationId?: string,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualificationsResults({
            pagination: {
                offset: offset.value,
            },
            qualificationId,
            status: status ?? [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

async function deleteQualificationResult(resultId: string): Promise<void> {
    // Remove result from list
    const idx = data.value?.results.findIndex((r) => r.id === resultId);
    if (idx !== undefined && idx > -1) {
        delete data.value?.results[idx];
    }
}
</script>

<template>
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.qualifications', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.results.length === 0"
                :message="$t('common.not_found', [$t('common.result', 2)])"
            />

            <template v-else>
                <GenericTable>
                    <template #thead>
                        <tr>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.citizen') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.status') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.score') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.summary') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.created_at') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.creator') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.action', 2) }}
                            </th>
                        </tr>
                    </template>
                    <template #tbody>
                        <QualificationsResultsListEntry
                            v-for="result in data?.results"
                            :key="`${result.qualificationId}-${result.userId}`"
                            :result="result"
                            @delete="deleteQualificationResult(result.id)"
                        />
                    </template>
                </GenericTable>
            </template>

            <TablePagination
                class="w-full"
                :pagination="data?.pagination"
                :show-border="false"
                :refresh="refresh"
                @offset-change="offset = $event"
            />
        </div>
    </div>
</template>
