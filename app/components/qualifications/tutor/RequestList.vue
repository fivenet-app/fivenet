<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { checkQualificationAccess, requestStatusToTextColor } from '~/components/qualifications/helpers';
import RequestTutorModal from '~/components/qualifications/tutor/RequestTutorModal.vue';
import ResultTutorModal from '~/components/qualifications/tutor/ResultTutorModal.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import {
    type Qualification,
    QualificationExamMode,
    type QualificationRequest,
    RequestStatus,
} from '~~/gen/ts/resources/qualifications/qualifications';
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
        searchQuery?: {
            users: number[];
        };
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
        examMode: QualificationExamMode.DISABLED,
        searchQuery: () => ({ users: [] }),
    },
);

const emit = defineEmits<{
    (e: 'refresh'): void;
}>();

const { t } = useI18n();

const overlay = useOverlay();

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
                id: 'createdAt',
                desc: true,
            },
        ],
    },
    page: 1,
});

const { data, status, refresh, error } = useLazyAsyncData(
    `qualifications-requests-${JSON.stringify(query.sorting)}-${query.page}-${props.qualification.id}-${JSON.stringify(props.searchQuery)}`,
    () => listQualificationRequests(props.qualification.id),
    {
        watch: [query],
    },
);

watchDebounced(
    () => props.searchQuery,
    async () => refresh(),
    { deep: true, debounce: 250, maxWait: 1250 },
);

defineExpose({
    refresh,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

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

async function deleteQualificationRequest(qualificationId: number, userId: number): Promise<DeleteQualificationReqResponse> {
    try {
        const call = qualificationsQualificationsClient.deleteQualificationReq({
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

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', [
                        row.original.status !== RequestStatus.DENIED &&
                            h(UTooltip, { text: t('common.decline') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-close-thick',
                                    color: 'orange',
                                    onClick: () => {
                                        requestTutorModal.open({
                                            request: row.original,
                                            status: RequestStatus.DENIED,
                                            onRefresh: onRefresh,
                                        });
                                    },
                                }),
                            ),
                        row.original.status !== RequestStatus.ACCEPTED &&
                            row.original.status !== RequestStatus.EXAM_STARTED &&
                            row.original.status !== RequestStatus.EXAM_GRADING &&
                            h(UTooltip, { text: t('common.accept') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-check-bold',
                                    color: 'green',
                                    onClick: () => {
                                        requestTutorModal.open({
                                            request: row.original,
                                            status: RequestStatus.ACCEPTED,
                                            onRefresh: onRefresh,
                                        });
                                    },
                                }),
                            ),
                        (row.original.status === RequestStatus.ACCEPTED ||
                            row.original.status === RequestStatus.EXAM_GRADING) &&
                            h(UTooltip, { text: t('common.grade') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-star',
                                    color: 'warning',
                                    onClick: () => {
                                        (row.original.status === RequestStatus.EXAM_GRADING
                                            ? examViewResultModal
                                            : resultTutorModal
                                        ).open({
                                            qualificationId: row.original.qualificationId,
                                            examMode: props.examMode,
                                            userId: row.original.userId,
                                            onRefresh: onRefresh,
                                        });
                                    },
                                }),
                            ),
                        checkQualificationAccess(
                            props.qualification.access,
                            props.qualification.creator,
                            AccessLevel.EDIT,
                            undefined,
                            props.qualification?.creatorJob,
                        ) &&
                            h(UTooltip, { text: t('common.delete') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-delete',
                                    color: 'error',
                                    onClick: () => {
                                        confirmModal.open({
                                            confirm: async () =>
                                                deleteQualificationRequest(row.original.qualificationId, row.original.userId),
                                        });
                                    },
                                }),
                            ),
                    ]),
            },
            {
                accessorKey: 'citizen',
                header: t('common.citizen'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.user }),
            },
            {
                accessorKey: 'userComment',
                header: t('common.comment'),
                cell: ({ row }) => row.original.userComment,
            },
            {
                accessorKey: 'status',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.status'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) =>
                    h(
                        'span',
                        { class: `font-medium ${requestStatusToTextColor(row.original.status)}` },
                        h(
                            'span',
                            { class: 'font-semibold' },
                            t(`enums.qualifications.RequestStatus.${RequestStatus[row.original.status ?? 0]}`),
                        ),
                    ),
            },
            {
                accessorKey: 'createdAt',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.created_at'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
            {
                accessorKey: 'approvedAt',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.approved_at'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) => h(GenericTime, { value: row.original.approvedAt }),
            },
            {
                accessorKey: 'approver',
                header: t('common.approver'),
                cell: ({ row }) => (row.original.approver ? h(CitizenInfoPopover, { user: row.original.approver }) : null),
            },
        ] as TableColumn<QualificationRequest>[],
);

async function onRefresh(): Promise<void> {
    emit('refresh');
    return refresh();
}

const confirmModal = overlay.create(ConfirmModal);
const requestTutorModal = overlay.create(RequestTutorModal);
const resultTutorModal = overlay.create(ResultTutorModal);
const examViewResultModal = overlay.create(ExamViewResultModal);
</script>

<template>
    <div>
        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.request', 2)])"
            :error="error"
            :retry="refresh"
        />

        <template v-else>
            <UTable
                v-model:sorting="query.sorting.columns"
                :columns="columns"
                :data="data?.requests"
                :loading="isRequestPending(status)"
                :empty="$t('common.not_found', [$t('common.request', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
            >
                <template #citizen-cell="{ row }">
                    <CitizenInfoPopover :user="row.original.user" />
                </template>
                <template #status-cell="{ row }">
                    <span class="font-medium" :class="requestStatusToTextColor(row.original.status)">
                        <span class="font-semibold">{{
                            $t(`enums.qualifications.RequestStatus.${RequestStatus[row.original.status ?? 0]}`)
                        }}</span>
                    </span>
                </template>
                <template #createdAt-cell="{ row }">
                    <GenericTime :value="row.original.createdAt" />
                </template>
                <template #approvedAt-cell="{ row }">
                    <GenericTime :value="row.original.approvedAt" />
                </template>
                <template #approver-cell="{ row }">
                    <CitizenInfoPopover v-if="row.original.approver" :user="row.original.approver" />
                </template>
                <template #actions-cell="{ row }">
                    <UTooltip v-if="row.original.status !== RequestStatus.DENIED" :text="$t('common.decline')">
                        <UButton
                            variant="link"
                            icon="i-mdi-close-thick"
                            color="orange"
                            @click="
                                requestTutorModal.open({
                                    request: row.original,
                                    status: RequestStatus.DENIED,
                                    onRefresh: onRefresh,
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip
                        v-if="
                            row.original.status !== RequestStatus.ACCEPTED &&
                            row.original.status !== RequestStatus.EXAM_STARTED &&
                            row.original.status !== RequestStatus.EXAM_GRADING
                        "
                        :text="$t('common.accept')"
                    >
                        <UButton
                            variant="link"
                            icon="i-mdi-check-bold"
                            color="green"
                            @click="
                                requestTutorModal.open({
                                    request: row.original,
                                    status: RequestStatus.ACCEPTED,
                                    onRefresh: onRefresh,
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip
                        v-if="
                            row.original.status === RequestStatus.ACCEPTED || row.original.status === RequestStatus.EXAM_GRADING
                        "
                        :text="$t('common.grade')"
                    >
                        <UButton
                            variant="link"
                            icon="i-mdi-star"
                            color="warning"
                            @click="
                                (row.original.status === RequestStatus.EXAM_GRADING
                                    ? examViewResultModal
                                    : resultTutorModal
                                ).open({
                                    qualificationId: row.original.qualificationId,
                                    examMode: examMode,
                                    userId: row.original.userId,
                                    onRefresh: onRefresh,
                                })
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
                                confirmModal.open({
                                    confirm: async () =>
                                        deleteQualificationRequest(row.original.qualificationId, row.original.userId),
                                })
                            "
                        />
                    </UTooltip>
                </template>
            </UTable>

            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </div>
</template>
