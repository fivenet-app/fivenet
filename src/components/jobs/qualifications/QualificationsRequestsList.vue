<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationRequestsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsRequestsListEntry from '~/components/jobs/qualifications/QualificationsRequestsListEntry.vue';
import type { RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import Pagination from '~/components/partials/Pagination.vue';

const props = withDefaults(
    defineProps<{
        qualificationId?: string;
        status?: RequestStatus[];
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
    },
);

const { $grpc } = useNuxtApp();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(
    `qualifications-requests-${page.value}-${props.qualificationId}`,
    () => listQualificationsRequests(props.qualificationId),
);

async function listQualificationsRequests(
    qualificationId?: string,
    status?: RequestStatus[],
): Promise<ListQualificationRequestsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualificationRequests({
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
</script>

<template>
    <div class="overflow-hidden rounded-lg bg-base-700 shadow">
        <div class="bg-background border-b border-gray-200 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                <div class="ml-4 mt-4">
                    <h3 class="text-base font-semibold leading-6 text-gray-200">
                        {{ $t('components.qualifications.user_requests') }}
                    </h3>
                </div>
            </div>
        </div>
        <div class="px-1 sm:px-2">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.request', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.requests.length === 0"
                :message="$t('common.not_found', [$t('common.request', 2)])"
                icon="i-mdi-account-school"
            />

            <template v-else>
                <ul role="list" class="divide-y divide-gray-100">
                    <QualificationsRequestsListEntry
                        v-for="request in data?.requests"
                        :key="`${request.qualificationId}-${request.userId}`"
                        :request="request"
                    />
                </ul>
            </template>
        </div>
        <div class="bg-background border-t border-gray-200 px-4 py-5 sm:p-6">
            <div class="-ml-4 -mt-4 flex items-center">
                <Pagination v-model="page" :pagination="data?.pagination" />
            </div>
        </div>
    </div>
</template>
