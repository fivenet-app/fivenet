<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationsResultsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsResultsListEntry from '~/components/jobs/qualifications/QualificationsResultsListEntry.vue';

const props = withDefaults(
    defineProps<{
        qualificationId?: string;
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

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(
    `qualifications-results-${page.value}-${props.qualificationId}-${props.userId}`,
    () => listQualificationsResults(props.qualificationId, props.userId, props.status),
);

async function listQualificationsResults(
    qualificationId?: string,
    userId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualificationsResults({
            pagination: {
                offset: offset.value,
            },
            qualificationId,
            userId,
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
</script>

<template>
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="bg-background border-b border-gray-200 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                <div class="ml-4 mt-4">
                    <h3 v-if="!userId" class="text-base font-semibold leading-6 text-gray-200">
                        {{ $t('components.qualifications.user_qualifications') }}
                    </h3>
                </div>
            </div>
        </div>
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.result', 2)])" />
            <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.result', 2)])" :retry="refresh" />
            <DataNoDataBlock
                v-else-if="data?.results.length === 0"
                :message="$t('common.not_found', [$t('common.result', 2)])"
                icon="i-mdi-sigma"
            />

            <template v-else>
                <ul role="list" class="divide-y divide-gray-100">
                    <QualificationsResultsListEntry v-for="result in data?.results" :key="result.id" :result="result" />
                </ul>
            </template>
        </div>
        <div class="bg-background border-t border-gray-200 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex items-center">
                <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                    <UPagination
                        v-model="page"
                        :page-count="data?.pagination?.pageSize ?? 0"
                        :total="data?.pagination?.totalCount ?? 0"
                    />
                </div>
            </div>
        </div>
    </div>
</template>
