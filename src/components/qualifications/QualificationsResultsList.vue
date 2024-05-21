<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationsResultsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsResultsListEntry from '~/components/qualifications/QualificationsResultsListEntry.vue';
import Pagination from '~/components/partials/Pagination.vue';

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

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualifications-results-${page.value}-${props.qualificationId}-${props.userId}`, () =>
    listQualificationsResults(props.qualificationId, props.userId, props.status),
);

async function listQualificationsResults(
    qualificationId?: string,
    userId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = getGRPCQualificationsClient().listQualificationsResults({
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
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
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
            </div>
        </template>

        <div>
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.result', 2)])" />
            <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.result', 2)])" :retry="refresh" />
            <DataNoDataBlock
                v-else-if="data?.results.length === 0"
                :message="$t('common.not_found', [$t('common.result', 2)])"
                icon="i-mdi-sigma"
            />

            <template v-else>
                <ul role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
                    <QualificationsResultsListEntry v-for="result in data?.results" :key="result.id" :result="result" />
                </ul>
            </template>
        </div>

        <template #footer>
            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" disable-border />
        </template>
    </UCard>
</template>
