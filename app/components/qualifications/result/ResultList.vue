<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import ResultListEntry from '~/components/qualifications/result/ResultListEntry.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
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

const _schema = z.object({
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'abbreviation',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'abbreviation', desc: true }] }),
    page: pageNumberSchema,
});

type Schema = z.output<typeof _schema>;

const query = reactive<Schema>({
    sorting: {
        columns: [
            {
                id: 'abbreviation',
                desc: true,
            },
        ],
    },
    page: 1,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const { data, status, refresh, error } = useLazyAsyncData(
    `qualifications-results-${JSON.stringify(query)}-${query.page}-${props.qualificationId}-${props.userId}`,
    () => listQualificationResults(props.qualificationId, props.userId, props.status),
    {
        watch: [query],
    },
);

async function listQualificationResults(
    qualificationId?: number,
    userId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualificationsResults({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
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
    <UCard :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-2xl leading-6 font-semibold">
                    {{ $t('common.qualification', 2) }}
                </h3>

                <SortButton v-model="query.sorting" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.result', 2)])" />
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

        <ul v-else class="divide-y divide-default" role="list">
            <ResultListEntry v-for="result in data?.results" :key="result.id" :result="result" />
        </ul>

        <template #footer>
            <Pagination
                v-model="query.page"
                :pagination="data?.pagination"
                :status="status"
                :refresh="refresh"
                disable-border
            />
        </template>
    </UCard>
</template>
