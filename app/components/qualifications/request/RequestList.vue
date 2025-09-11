<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import RequestListEntry from '~/components/qualifications/request/RequestListEntry.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
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
    `qualifications-requests-${JSON.stringify(query.sorting)}-${query.page}-${props.qualificationId}`,
    () => listQualificationRequests(props.qualificationId),
);

async function listQualificationRequests(
    qualificationId?: number,
    status?: RequestStatus[],
): Promise<ListQualificationRequestsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualificationRequests({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
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
    <UCard :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-2xl leading-6 font-semibold">
                    {{ $t('components.qualifications.user_requests') }}
                </h3>

                <SortButton v-model="query.sorting" :fields="[{ label: $t('common.id'), value: 'id' }]" />
            </div>
        </template>

        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.request', 2)])" />
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

        <ul v-else class="divide-y divide-default" role="list">
            <RequestListEntry
                v-for="request in data?.requests"
                :key="`${request.qualificationId}-${request.userId}`"
                :request="request"
            />
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
