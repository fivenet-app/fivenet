<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { checkQualificationAccess, requestStatusToTextColor } from '~/components/qualifications/helpers';
import QualificationRequestTutorModal from '~/components/qualifications/tutor/QualificationRequestTutorModal.vue';
import QualificationResultTutorModal from '~/components/qualifications/tutor/QualificationResultTutorModal.vue';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import { type Qualification, QualificationExamMode, RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type {
    DeleteQualificationReqResponse,
    ListQualificationRequestsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import ExamViewResultModal from './ExamViewResultModal.vue';

const props = withDefaults(
    defineProps<{
        qualification: Qualification;
        status?: RequestStatus[];
        examMode?: QualificationExamMode;
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
        examMode: QualificationExamMode.DISABLED,
    },
);

const emit = defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const page = useRouteQuery('page', '1', { transform: Number });

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const { data, status, refresh, error } = useLazyAsyncData(
    `qualifications-requests-${sort.value.column}:${sort.value.direction}-${page.value}-${props.qualification.id}`,
    () => listQualificationRequests(props.qualification.id),
    {
        watch: [sort],
    },
);

async function listQualificationRequests(
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

async function deleteQualificationRequest(qualificationId: number, userId: number): Promise<DeleteQualificationReqResponse> {
    try {
        const call = $grpc.qualifications.qualifications.deleteQualificationReq({
            qualificationId: qualificationId,
            userId,
        });
        const { response } = await call;

        onRefresh();

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'citizen',
        label: t('common.citizen'),
    },
    {
        key: 'userComment',
        label: t('common.comment'),
    },
    {
        key: 'status',
        label: t('common.status'),
        sortable: true,
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
        sortable: true,
    },
    {
        key: 'approvedAt',
        label: t('common.approved_at'),
        sortable: true,
    },
    {
        key: 'approver',
        label: t('common.approver'),
    },
];

async function onRefresh(): Promise<void> {
    emit('refresh');
    return refresh();
}

defineExpose({
    refresh,
});
</script>

<template>
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2">
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                :error="error"
                :retry="refresh"
            />

            <template v-else>
                <UTable
                    v-model:sort="sort"
                    :loading="isRequestPending(status)"
                    :columns="columns"
                    :rows="data?.requests"
                    :empty-state="{ icon: 'i-mdi-account-school', label: $t('common.not_found', [$t('common.request', 2)]) }"
                    sort-mode="manual"
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
                        <div :key="`${request.userId}-${request.approverId}-${request.status}`">
                            <UTooltip v-if="request.status !== RequestStatus.DENIED" :text="$t('common.decline')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-close-thick"
                                    color="orange"
                                    @click="
                                        modal.open(QualificationRequestTutorModal, {
                                            request: request,
                                            status: RequestStatus.DENIED,
                                            onRefresh: onRefresh,
                                        })
                                    "
                                />
                            </UTooltip>

                            <UTooltip
                                v-if="
                                    request.status !== RequestStatus.ACCEPTED &&
                                    request.status !== RequestStatus.EXAM_STARTED &&
                                    request.status !== RequestStatus.EXAM_GRADING
                                "
                                :text="$t('common.accept')"
                            >
                                <UButton
                                    variant="link"
                                    icon="i-mdi-check-bold"
                                    color="green"
                                    @click="
                                        modal.open(QualificationRequestTutorModal, {
                                            request: request,
                                            status: RequestStatus.ACCEPTED,
                                            onRefresh: onRefresh,
                                        })
                                    "
                                />
                            </UTooltip>

                            <UTooltip
                                v-if="
                                    request.status === RequestStatus.ACCEPTED || request.status === RequestStatus.EXAM_GRADING
                                "
                                :text="$t('common.grade')"
                            >
                                <UButton
                                    variant="link"
                                    icon="i-mdi-star"
                                    color="amber"
                                    @click="
                                        modal.open(
                                            request.status === RequestStatus.EXAM_GRADING
                                                ? ExamViewResultModal
                                                : QualificationResultTutorModal,
                                            {
                                                qualificationId: request.qualificationId,
                                                examMode: examMode,
                                                userId: request.userId,
                                                onRefresh: onRefresh,
                                            },
                                        )
                                    "
                                />
                            </UTooltip>

                            <UTooltip
                                v-if="
                                    checkQualificationAccess(
                                        qualification.access,
                                        qualification.creator,
                                        AccessLevel.EDIT,
                                        undefined,
                                        qualification?.creatorJob,
                                    )
                                "
                                :text="$t('common.delete')"
                            >
                                <UButton
                                    variant="link"
                                    icon="i-mdi-delete"
                                    color="error"
                                    @click="
                                        modal.open(ConfirmModal, {
                                            confirm: async () =>
                                                deleteQualificationRequest(request.qualificationId, request.userId),
                                        })
                                    "
                                />
                            </UTooltip>
                        </div>
                    </template>
                </UTable>

                <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
            </template>
        </div>
    </div>
</template>
