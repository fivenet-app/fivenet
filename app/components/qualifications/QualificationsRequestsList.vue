<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import QualificationsRequestsListEntry from '~/components/qualifications/QualificationsRequestsListEntry.vue';
import type { RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationRequestsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualificationId?: number;
        status?: RequestStatus[];
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
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
    `qualifications-requests-${sort.value.column}:${sort.value.direction}-${page.value}-${props.qualificationId}`,
    () => listQualificationsRequests(props.qualificationId),
);

async function listQualificationsRequests(
    qualificationId?: number,
    status?: RequestStatus[],
): Promise<ListQualificationRequestsResponse> {
    try {
        const call = $grpc.qualifications.qualifications.listQualificationRequests({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            sort: sort.value,
            qualificationId: qualificationId,
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
                    {{ $t('components.qualifications.user_requests') }}
                </h3>

                <SortButton v-model="sort" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <div>
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.request', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.requests.length === 0"
                :message="$t('common.not_found', [$t('common.request', 2)])"
                icon="i-mdi-account-school"
            />

            <ul v-else class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                <QualificationsRequestsListEntry
                    v-for="request in data?.requests"
                    :key="`${request.qualificationId}-${request.userId}`"
                    :request="request"
                />
            </ul>
        </div>

        <template #footer>
            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" disable-border />
        </template>
    </UCard>
</template>
