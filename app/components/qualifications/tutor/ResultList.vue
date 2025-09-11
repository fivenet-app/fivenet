<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { checkQualificationAccess, resultStatusToTextColor } from '~/components/qualifications/helpers';
import ExamViewResultModal from '~/components/qualifications/tutor/ExamViewResultModal.vue';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import {
    type Qualification,
    QualificationExamMode,
    type QualificationResult,
    ResultStatus,
} from '~~/gen/ts/resources/qualifications/qualifications';
import type {
    DeleteQualificationResultResponse,
    ListQualificationsResultsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualification: Qualification;
        status?: ResultStatus[];
        examMode?: QualificationExamMode;
    }>(),
    {
        qualificationId: undefined,
        status: () => [],
        examMode: QualificationExamMode.UNSPECIFIED,
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
    `qualifications-results-${JSON.stringify(query)}-${query.page}-${props.qualification.id}`,
    () => listQualificationResults(props.qualification.id, props.status),
    {
        watch: [query],
    },
);

defineExpose({
    refresh,
});

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

async function listQualificationResults(
    qualificationId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = qualificationsQualificationsClient.listQualificationsResults({
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

async function deleteQualificationResult(resultId: number): Promise<DeleteQualificationResultResponse> {
    try {
        const call = qualificationsQualificationsClient.deleteQualificationResult({
            resultId,
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
                        row.original.status === ResultStatus.PENDING &&
                            h(UTooltip, { text: t('common.grade') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-star',
                                    color: 'warning',
                                    onClick: () => {
                                        examViewResultModal.open({
                                            qualificationId: row.original.qualificationId,
                                            userId: row.original.userId,
                                            resultId: row.original.id,
                                            examMode: props.examMode,
                                            onRefresh: onRefresh,
                                        });
                                    },
                                }),
                            ),
                        props.examMode > QualificationExamMode.DISABLED &&
                            h(UTooltip, { text: t('common.show') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-star',
                                    color: 'warning',
                                    onClick: () => {
                                        examViewResultModal.open({
                                            qualificationId: row.original.qualificationId,
                                            userId: row.original.userId,
                                            resultId: row.original.id,
                                            examMode: props.examMode,
                                            viewOnly: true,
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
                                    class: 'flex-initial',
                                    variant: 'link',
                                    icon: 'i-mdi-delete',
                                    color: 'error',
                                    onClick: () => {
                                        confirmModal.open({
                                            confirm: async () => deleteQualificationResult(row.original.id),
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
                cell: ({ row }) => (row.original.score ? $n(row.original.score) : null),
            },
            {
                accessorKey: 'summary',
                header: t('common.summary'),
                cell: ({ row }) => (row.original.summary ? h('p', { class: 'text-sm' }, row.original.summary) : null),
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
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => (row.original.creator ? h(CitizenInfoPopover, { user: row.original.creator }) : null),
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
    <div class="overflow-hidden">
        <div class="px-1 sm:px-2">
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :error="error"
                :retry="refresh"
            />

            <template v-else>
                <UTable
                    v-model:sorting="query.sorting.columns"
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
    </div>
</template>
