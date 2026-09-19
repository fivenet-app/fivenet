<script lang="ts" setup>
import { UButton, UDropdownMenu, UTooltip } from '#components';
import type { DropdownMenuItem, TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DeletedAtBadge from '~/components/partials/DeletedAtBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import TableSortButton from '~/components/partials/TableSortButton.vue';
import { checkQualificationAccess, requestStatusToTextColor } from '~/components/qualifications/helpers';
import RequestTutorModal from '~/components/qualifications/tutor/RequestTutorModal.vue';
import ResultTutorModal from '~/components/qualifications/tutor/ResultTutorModal.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access/access';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/exam/exam';
import {
    type Qualification,
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
const selectedStatuses = ref<RequestStatus[]>(props.status ?? []);
const notifyUser = ref(true);

const statusOptions = computed(() =>
    [
        RequestStatus.PENDING,
        RequestStatus.ACCEPTED,
        RequestStatus.EXAM_STARTED,
        RequestStatus.EXAM_GRADING,
        RequestStatus.COMPLETED,
        RequestStatus.DENIED,
    ].map((status) => ({
        status,
        label: t(`enums.qualifications.RequestStatus.${RequestStatus[status]}`),
    })),
);

const { data, status, refresh, error } = useAuthedLazyAsyncData(
    'userState',
    `qualifications-requests:${query.page}-${JSON.stringify(query)}-${JSON.stringify(selectedStatuses.value)}-${props.qualification.id}-${JSON.stringify(props.searchQuery)}`,
    ({ signal }) => listQualificationRequests(props.qualification.id, undefined, signal),
);

useDebouncedRefresh([query, selectedStatuses, () => props.searchQuery], refresh, { debounce: 250, maxWait: 1250 });

defineExpose({
    refresh,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

async function listQualificationRequests(
    qualificationId?: number,
    status?: RequestStatus[],
    signal?: AbortSignal,
): Promise<ListQualificationRequestsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualificationRequests(
            {
                pagination: {
                    offset: calculateOffset(query.page, data.value?.pagination),
                },
                sort: query.sorting,
                qualificationId: qualificationId,
                status: selectedStatuses.value.length > 0 ? selectedStatuses.value : (status ?? []),
                userIds: props.searchQuery.users,
            },
            { abort: signal },
        );
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteOrRestoreQualificationRequest(
    qualificationId: number,
    userId: number,
    skipNotification: boolean,
): Promise<DeleteQualificationReqResponse> {
    try {
        const call = qualificationsQualificationsClient.deleteQualificationReq({
            qualificationId,
            userId,
            skipNotification,
        });
        const { response } = await call;

        onRefresh();

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function getRowActions(request: QualificationRequest): DropdownMenuItem[][] {
    const actions: DropdownMenuItem[] = [];

    if (request.deletedAt) {
        actions.push({
            label: t('common.restore'),
            icon: 'i-mdi-restore',
            color: 'success',
            onSelect: () =>
                confirmModal.open({
                    color: 'success',
                    icon: 'i-mdi-restore',
                    confirm: async () => deleteOrRestoreQualificationRequest(request.qualificationId, request.userId, true),
                }),
        });
        return [actions];
    }

    if (request.status === RequestStatus.PENDING || request.status === RequestStatus.ACCEPTED) {
        actions.push({
            label: t('common.decline'),
            icon: 'i-mdi-close-thick',
            color: 'warning',
            onSelect: () => requestTutorModal.open({ request, status: RequestStatus.DENIED, onRefresh }),
        });
    }

    if (request.status === RequestStatus.PENDING || request.status === RequestStatus.DENIED) {
        actions.push({
            label: t('common.accept'),
            icon: 'i-mdi-check-bold',
            color: 'success',
            onSelect: () => requestTutorModal.open({ request, status: RequestStatus.ACCEPTED, onRefresh }),
        });
    }

    if (request.status === RequestStatus.ACCEPTED || request.status === RequestStatus.EXAM_GRADING) {
        actions.push({
            label: t('common.grade'),
            icon: 'i-mdi-star',
            color: 'warning',
            onSelect: () =>
                (request.status === RequestStatus.EXAM_GRADING ? examViewResultModal : resultTutorModal).open({
                    qualificationId: request.qualificationId,
                    qualification: props.qualification,
                    examMode: props.examMode,
                    userId: request.userId,
                    onRefresh,
                }),
        });
    }

    if (
        checkQualificationAccess(
            props.qualification.access,
            props.qualification.creator,
            AccessLevel.EDIT,
            undefined,
            props.qualification.creatorJob,
        )
    ) {
        if (actions.length > 0) actions.push({ type: 'separator' as const });
        actions.push({
            label: t('common.delete'),
            icon: 'i-mdi-delete',
            color: 'error',
            onSelect: () =>
                confirmModal.open({
                    notifyUser: notifyUser.value,
                    onNotifyUserUpdate: (value) => (notifyUser.value = value),
                    confirm: async () =>
                        deleteOrRestoreQualificationRequest(request.qualificationId, request.userId, !notifyUser.value),
                }),
        });
    }

    return [actions];
}

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) => {
                    const items = getRowActions(row.original);
                    return h(
                        'div',
                        { class: 'flex min-h-7 items-center' },
                        items.length > 0
                            ? h(
                                  UDropdownMenu,
                                  { items: items, content: { align: 'end' } },
                                  {
                                      default: () =>
                                          h(
                                              UTooltip,
                                              { text: t('common.action', 2) },
                                              {
                                                  default: () =>
                                                      h(UButton, {
                                                          color: 'neutral',
                                                          variant: 'ghost',
                                                          icon: 'i-mdi-dots-vertical',
                                                          'aria-label': t('common.action', 2),
                                                      }),
                                              },
                                          ),
                                  },
                              )
                            : undefined,
                    );
                },
            },
            {
                accessorKey: 'citizen',
                header: t('common.citizen'),
                cell: ({ row }) =>
                    h('div', { class: 'flex items-center gap-2' }, [
                        h(CitizenInfoPopover, { user: row.original.user }),
                        row.original.deletedAt
                            ? h(DeletedAtBadge, {
                                  deletedAt: row.original.deletedAt,
                                  iconOnly: true,
                                  'aria-label': t('common.deleted'),
                                  title: t('common.deleted'),
                              })
                            : null,
                    ]),
            },
            {
                accessorKey: 'userComment',
                header: t('common.comment'),
                cell: ({ row }) =>
                    h('p', { class: 'line-clamp-1 whitespace-normal hover:line-clamp-3' }, row.original.userComment),
                meta: {
                    class: {
                        td: 'max-w-48 w-full min-w-0',
                    },
                },
            },
            {
                accessorKey: 'status',
                header: ({ column }) => {
                    return h(TableSortButton, {
                        column,
                        label: t('common.status'),
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
                    return h(TableSortButton, {
                        column,
                        label: t('common.created_at'),
                    });
                },
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
            {
                accessorKey: 'approvedAt',
                header: ({ column }) => {
                    return h(TableSortButton, {
                        column,
                        label: t('common.approved_at'),
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
    <div class="flex h-full min-h-0 flex-col overflow-x-hidden">
        <UDashboardToolbar>
            <div class="flex w-full items-center justify-between gap-2">
                <h2 class="text-lg font-semibold text-highlighted">{{ $t('common.request', 2) }}</h2>

                <USelectMenu
                    v-model="selectedStatuses"
                    class="w-full sm:w-64"
                    multiple
                    value-key="status"
                    :items="statusOptions"
                    :placeholder="$t('common.status', 2)"
                >
                    <template #default>
                        <span v-if="selectedStatuses.length === 0">{{ $t('common.all') }}</span>
                        <span v-else>{{ $t('common.selected', selectedStatuses.length) }}</span>
                    </template>
                    <template #item-label="{ item }">{{ item.label }}</template>
                </USelectMenu>
            </div>
        </UDashboardToolbar>

        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.request', 2)])"
            :error="error"
            :retry="refresh"
        />

        <template v-else>
            <UTable
                v-model:sorting="query.sorting.columns"
                class="my-0 min-h-0 flex-1 overflow-auto"
                :columns="columns"
                :data="data?.requests"
                :loading="isRequestPending(status)"
                :empty="$t('common.not_found', [$t('common.request', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
            />

            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </div>
</template>
