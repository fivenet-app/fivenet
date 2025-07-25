<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import QualificationsResultsListEntry from '~/components/qualifications/QualificationsResultsListEntry.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationsResultsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualificationId?: number;
        userId?: number;
        status?: ResultStatus[];
    }>(),
    {
        qualificationId: undefined,
        userId: undefined,
        status: () => [ResultStatus.SUCCESSFUL],
    },
);

const { $grpc } = useNuxtApp();

const page = useRouteQuery('page', '1', { transform: Number });

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'abbreviation',
    direction: 'desc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `qualifications-results-${sort.value.column}:${sort.value.direction}-${page.value}-${props.qualificationId}-${props.userId}`,
    () => listQualificationsResults(props.qualificationId, props.userId, props.status),
    {
        watch: [sort],
    },
);

async function listQualificationsResults(
    qualificationId?: number,
    userId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = $grpc.qualifications.qualifications.listQualificationsResults({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            sort: sort.value,
            qualificationId: qualificationId,
            userId,
            status: status ?? [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UCard
        :ui="{
            body: { padding: '' },
        }"
    >
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-2xl font-semibold leading-6">
                    {{ $t('common.qualification', 2) }}
                </h3>

                <SortButton v-model="sort" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <div>
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.result', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.result', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.results.length === 0"
                :message="$t('common.not_found', [$t('common.result', 2)])"
                icon="i-mdi-sigma"
            />

            <ul v-else class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                <QualificationsResultsListEntry v-for="result in data?.results" :key="result.id" :result="result" />
            </ul>
        </div>

        <template #footer>
            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" disable-border />
        </template>
    </UCard>
</template>
