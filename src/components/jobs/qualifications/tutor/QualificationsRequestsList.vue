<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type {
    DeleteQualificationReqResponse,
    ListQualificationRequestsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import { RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import QualificationRequestTutorModal from '~/components/jobs/qualifications/tutor/QualificationRequestTutorModal.vue';
import QualificationResultTutorModal from '~/components/jobs/qualifications/tutor/QualificationResultTutorModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { requestStatusToTextColor } from '../helpers';

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

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const modal = useModal();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualifications-requests-${page.value}-${props.qualificationId}`, () =>
    listQualificationsRequests(props.qualificationId),
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

async function deleteQualificationRequest(qualificationId: string, userId: number): Promise<DeleteQualificationReqResponse> {
    try {
        const call = $grpc.getQualificationsClient().deleteQualificationReq({
            qualificationId,
            userId,
        });
        const { response } = await call;

        refresh();

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'citizen',
        label: t('common.citizen'),
    },
    {
        key: 'comment',
        label: t('common.comment'),
    },
    {
        key: 'status',
        label: t('common.status'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'approvedAt',
        label: t('common.approved_at'),
    },
    {
        key: 'approver',
        label: t('common.approver'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
    },
];
</script>

<template>
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.request', 2)])" :retry="refresh" />

            <template v-else>
                <UTable
                    :loading="loading"
                    :columns="columns"
                    :rows="data?.requests"
                    :empty-state="{ icon: 'i-mdi-account-school', label: $t('common.not_found', [$t('common.request', 2)]) }"
                >
                    <template #citizen-data="{ row: request }">
                        <CitizenInfoPopover :user="request.user" />
                    </template>
                    <template #status-data="{ row: request }">
                        <span class="font-medium" :class="requestStatusToTextColor(request.status)">
                            <span class="font-semibold">{{
                                $t(`enums.qualifications.RequestStatus.${RequestStatus[request.status]}`)
                            }}</span>
                        </span>
                    </template>
                    <template #createdAt-data="{ row: request }">
                        <GenericTime :value="request.createdAt" />
                    </template>
                    <template #approvedAt-data="{ row: request }">
                        <GenericTime :value="request.approvedAt" />
                    </template>
                    <template #approver-data="{ row: request }">
                        <CitizenInfoPopover v-if="request.approver" :user="request.approver" />
                    </template>
                    <template #actions-data="{ row: request }">
                        <UButton
                            v-if="request.status !== RequestStatus.DENIED"
                            variant="link"
                            icon="i-mdi-close-thick"
                            color="red"
                            @click="
                                modal.open(QualificationRequestTutorModal, {
                                    request: request,
                                    status: RequestStatus.DENIED,
                                    onRefresh: refresh,
                                })
                            "
                        />
                        <UButton
                            v-if="request.status !== RequestStatus.ACCEPTED"
                            variant="link"
                            icon="i-mdi-check-bold"
                            color="green"
                            @click="
                                modal.open(QualificationRequestTutorModal, {
                                    request: request,
                                    status: RequestStatus.ACCEPTED,
                                    onRefresh: refresh,
                                })
                            "
                        />
                        <UButton
                            v-if="request.status === RequestStatus.ACCEPTED"
                            variant="link"
                            icon="i-mdi-star"
                            color="amber"
                            @click="
                                modal.open(QualificationResultTutorModal, {
                                    qualificationId: request.qualificationId,
                                    userId: request.userId,
                                    onRefresh: refresh,
                                })
                            "
                        />
                        <UButton
                            v-if="can('QualificationsService.DeleteQualificationReq')"
                            variant="link"
                            icon="i-mdi-trash-can"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteQualificationRequest(request.qualificationId, request.userId),
                                })
                            "
                        />
                    </template>
                </UTable>

                <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                    <UPagination
                        v-model="page"
                        :page-count="data?.pagination?.pageSize ?? 0"
                        :total="data?.pagination?.totalCount ?? 0"
                    />
                </div>
            </template>
        </div>
    </div>
</template>
