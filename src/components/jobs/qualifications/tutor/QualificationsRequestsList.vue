<script lang="ts" setup>
import { AccountSchoolIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { ListQualificationRequestsResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationsRequestsListEntry from '~/components/jobs/qualifications/tutor/QualificationsRequestsListEntry.vue';
import { RequestStatus, QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import GenericTable from '~/components/partials/elements/GenericTable.vue';
import QualificationRequestTutorModal from '~/components/jobs/qualifications/tutor/QualificationRequestTutorModal.vue';
import QualificationResultTutorModal from '~/components/jobs/qualifications/tutor/QualificationResultTutorModal.vue';

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
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

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

const selectedRequestStatus = ref<undefined | RequestStatus>();
const selectedRequest = ref<undefined | QualificationRequest>();

const openResultStatus = ref(false);
</script>

<template>
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.request', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.requests.length === 0"
                :message="$t('common.not_found', [$t('common.request', 2)])"
                :icon="markRaw(AccountSchoolIcon)"
            />

            <template v-else>
                <QualificationRequestTutorModal
                    v-if="selectedRequest && !openResultStatus"
                    :request="selectedRequest"
                    :status="selectedRequestStatus"
                    @close="selectedRequest = undefined"
                    @refresh="refresh()"
                />

                <QualificationResultTutorModal
                    v-if="selectedRequest"
                    :open="openResultStatus"
                    :qualification-id="selectedRequest.qualificationId"
                    :user-id="selectedRequest.userId"
                    @close="
                        selectedRequest = undefined;
                        openResultStatus = false;
                    "
                    @refresh="refresh()"
                />

                <GenericTable>
                    <template #thead>
                        <tr>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.citizen') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.comment') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.status') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.created_at') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.approved_at') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.approver') }}
                            </th>
                            <th scope="col" class="whitespace-nowrap p-1 text-left text-sm font-semibold text-gray-100">
                                {{ $t('common.action', 2) }}
                            </th>
                        </tr>
                    </template>
                    <template #tbody>
                        <QualificationsRequestsListEntry
                            v-for="request in data?.requests"
                            :key="`${request.qualificationId}-${request.userId}`"
                            :request="request"
                            @selected-request-status="
                                selectedRequestStatus = $event;
                                selectedRequest = request;
                            "
                            @grade-request="
                                selectedRequest = request;
                                openResultStatus = true;
                            "
                            @delete="refresh()"
                        />
                    </template>
                </GenericTable>
            </template>

            <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
                <UPagination
                    v-model="page"
                    :page-count="data?.pagination?.pageSize ?? 0"
                    :total="data?.pagination?.totalCount ?? 0"
                />
            </div>
        </div>
    </div>
</template>
