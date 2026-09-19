<script lang="ts" setup>
import { UBadge, UButton, UDropdownMenu, UTooltip } from '#components';
import type { DropdownMenuItem, TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DeletedAtBadge from '~/components/partials/DeletedAtBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import TableSortButton from '~/components/partials/TableSortButton.vue';
import { checkQualificationAccess, resultStatusToTextColor } from '~/components/qualifications/helpers';
import ExamViewResultModal from '~/components/qualifications/tutor/ExamViewResultModal.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access/access';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/exam/exam';
import { type Qualification, type QualificationResult, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type {
    DeleteQualificationResultResponse,
    ListQualificationsResultsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualification: Qualification;
        status?: ResultStatus[];
        examMode?: QualificationExamMode;
        searchQuery?: {
            users: number[];
        };
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
        examMode: QualificationExamMode.UNSPECIFIED,
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
const selectedStatuses = ref<ResultStatus[]>(props.status ?? []);
const notifyUser = ref(true);

const statusOptions = computed(() =>
    [ResultStatus.PENDING, ResultStatus.SUCCESSFUL, ResultStatus.FAILED].map((status) => ({
        status,
        label: t(`enums.qualifications.ResultStatus.${ResultStatus[status]}`),
    })),
);

const { data, status, refresh, error } = useAuthedLazyAsyncData(
    'userState',
    () =>
        `qualifications-results:${query.page}-${JSON.stringify(query)}-${JSON.stringify(selectedStatuses.value)}-${props.qualification.id}-${JSON.stringify(props.searchQuery)}`,
    ({ signal }) => listQualificationResults(props.qualification.id, selectedStatuses.value, signal),
);

useDebouncedRefresh([query, selectedStatuses, () => props.searchQuery], refresh, { debounce: 250, maxWait: 1250 });

defineExpose({
    refresh,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

async function listQualificationResults(
    qualificationId?: number,
    status?: ResultStatus[],
    signal?: AbortSignal,
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualificationsResults(
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

async function deleteQualificationResult(
    resultId: number,
    skipNotification: boolean,
): Promise<DeleteQualificationResultResponse> {
    try {
        const call = qualificationsQualificationsClient.deleteQualificationResult({
            resultId,
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

function getRowActions(result: QualificationResult): DropdownMenuItem[][] {
    const actions: DropdownMenuItem[] = [];

    if (result.status === ResultStatus.PENDING && !result.deletedAt) {
        actions.push({
            label: t('common.grade'),
            icon: 'i-mdi-star',
            onSelect: () =>
                examViewResultModal.open({
                    qualificationId: result.qualificationId,
                    qualification: props.qualification,
                    userId: result.userId,
                    resultId: result.id,
                    examMode: props.examMode,
                    onRefresh,
                }),
        });
    } else if (props.examMode > QualificationExamMode.DISABLED) {
        actions.push({
            label: t('common.show'),
            icon: 'i-mdi-eye',
            onSelect: () =>
                examViewResultModal.open({
                    qualificationId: result.qualificationId,
                    qualification: props.qualification,
                    userId: result.userId,
                    resultId: result.id,
                    examMode: props.examMode,
                    viewOnly: true,
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
        actions.push({ type: 'separator' as const });
        actions.push({
            label: t(result.deletedAt ? 'common.restore' : 'common.delete'),
            icon: result.deletedAt ? 'i-mdi-restore' : 'i-mdi-delete',
            color: result.deletedAt ? 'success' : 'error',
            onSelect: () => {
                confirmModal.open({
                    color: result.deletedAt ? 'success' : 'error',
                    icon: result.deletedAt ? 'i-mdi-restore' : 'i-mdi-warning-circle',
                    notifyUser: result.deletedAt ? undefined : notifyUser.value,
                    onNotifyUserUpdate: result.deletedAt ? undefined : (value) => (notifyUser.value = value),
                    confirm: async () =>
                        deleteQualificationResult(result.id, result.deletedAt !== undefined || !notifyUser.value),
                });
            },
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
                                  { items, content: { align: 'end' } },
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
                accessorKey: 'status',
                header: ({ column }) => {
                    return h(TableSortButton, {
                        column,
                        label: t('common.status'),
                    });
                },
                cell: ({ row }) =>
                    row.original.status !== undefined
                        ? h(
                              'span',
                              { class: `font-medium ${resultStatusToTextColor(row.original.status)}` },
                              h(
                                  'span',
                                  { class: 'font-semibold' },
                                  t(`enums.qualifications.ResultStatus.${ResultStatus[row.original.status]}`),
                              ),
                          )
                        : null,
            },
            {
                accessorKey: 'score',
                header: t('common.score'),
                cell: ({ row }) => (row.original.score !== undefined ? $n(row.original.score) : null),
            },
            {
                accessorKey: 'summary',
                header: t('common.summary'),
                cell: ({ row }) =>
                    row.original.summary
                        ? h('p', { class: 'line-clamp-1 whitespace-normal hover:line-clamp-3' }, row.original.summary)
                        : null,
                meta: {
                    class: {
                        td: 'max-w-48 w-full min-w-0',
                    },
                },
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
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) =>
                    row.original.autoGraded
                        ? h(UBadge, { color: 'info', icon: 'i-mdi-robot', label: t('components.qualifications.auto_graded') })
                        : row.original.creator
                          ? h(CitizenInfoPopover, { user: row.original.creator })
                          : null,
            },
        ] as TableColumn<QualificationResult>[],
);

async function onRefresh(): Promise<void> {
    emit('refresh');
    return refresh();
}

const examViewResultModal = overlay.create(ExamViewResultModal);
const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <div class="flex h-full min-h-0 flex-col overflow-x-hidden">
        <UDashboardToolbar>
            <div class="flex w-full items-center justify-between gap-2">
                <h2 class="text-lg font-semibold text-highlighted">{{ $t('common.result', 2) }}</h2>

                <div class="flex items-center gap-2">
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

                    <slot name="actions" />
                </div>
            </div>
        </UDashboardToolbar>

        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
            :error="error"
            :retry="refresh"
        />

        <template v-else>
            <UTable
                v-model:sorting="query.sorting.columns"
                class="my-0 min-h-0 flex-1 overflow-auto"
                :columns="columns"
                :data="data?.results"
                :loading="isRequestPending(status)"
                :empty="$t('common.not_found', [$t('common.result', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
            />

            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </div>
</template>
