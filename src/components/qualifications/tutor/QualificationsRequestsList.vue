<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type {
    DeleteQualificationReqResponse,
    ListQualificationRequestsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import { RequestStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import QualificationRequestTutorModal from '~/components/qualifications/tutor/QualificationRequestTutorModal.vue';
import QualificationResultTutorModal from '~/components/qualifications/tutor/QualificationResultTutorModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { requestStatusToTextColor } from '../helpers';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

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

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { t } = useI18n();

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
        const call = getGRPCQualificationsClient().listQualificationRequests({
            pagination: {
                offset: offset.value,
            },
            qualificationId,
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

async function deleteQualificationRequest(qualificationId: string, userId: number): Promise<DeleteQualificationReqResponse> {
    try {
        const call = getGRPCQualificationsClient().deleteQualificationReq({
            qualificationId,
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
        sortable: false,
    },
];

async function onRefresh(): Promise<void> {
    emits('refresh');
    return refresh();
}

defineExpose({
    refresh,
});
</script>

<template>
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2">
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
                        <div :key="`${request.userId}-${request.approverId}-${request.status}`">
                            <UButton
                                v-if="request.status !== RequestStatus.DENIED"
                                variant="link"
                                icon="i-mdi-close-thick"
                                color="red"
                                @click="
                                    modal.open(QualificationRequestTutorModal, {
                                        request: request,
                                        status: RequestStatus.DENIED,
                                        onRefresh: onRefresh,
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
                                        onRefresh: onRefresh,
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
                                        onRefresh: onRefresh,
                                    })
                                "
                            />

                            <UButton
                                v-if="can('QualificationsService.DeleteQualificationReq')"
                                variant="link"
                                icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () =>
                                            deleteQualificationRequest(request.qualificationId, request.userId),
                                    })
                                "
                            />
                        </div>
                    </template>
                </UTable>

                <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
            </template>
        </div>
    </div>
</template>
