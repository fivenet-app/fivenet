<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { checkQualificationAccess, resultStatusToTextColor } from '~/components/qualifications/helpers';
import ExamViewResultModal from '~/components/qualifications/tutor/ExamViewResultModal.vue';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import { type Qualification, QualificationExamMode, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
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

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const page = useRouteQuery('page', '1', { transform: Number });

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const { data, status, refresh, error } = useLazyAsyncData(
    `qualifications-results-${sort.value.column}:${sort.value.direction}-${page.value}-${props.qualification.id}`,
    () => listQualificationsResults(props.qualification.id, props.status),
    {
        watch: [sort],
    },
);

async function listQualificationsResults(
    qualificationId?: number,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = $grpc.qualifications.qualifications.listQualificationsResults({
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

async function deleteQualificationResult(resultId: number): Promise<DeleteQualificationResultResponse> {
    try {
        const call = $grpc.qualifications.qualifications.deleteQualificationResult({
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
        key: 'status',
        label: t('common.status'),
        sortable: true,
    },
    {
        key: 'score',
        label: t('common.score'),
    },
    {
        key: 'summary',
        label: t('common.summary'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
        sortable: true,
    },
    {
        key: 'creator',
        label: t('common.creator'),
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
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :error="error"
                :retry="refresh"
            />

            <template v-else>
                <UTable
                    v-model:sort="sort"
                    :loading="isRequestPending(status)"
                    :columns="columns"
                    :rows="data?.results"
                    :empty-state="{ icon: 'i-mdi-sigma', label: $t('common.not_found', [$t('common.result', 2)]) }"
                    sort-mode="manual"
                >
                    <template #citizen-data="{ row: result }">
                        <CitizenInfoPopover :user="result.user" />
                    </template>

                    <template #status-data="{ row: result }">
                        <template v-if="result.status !== undefined">
                            <span class="font-medium" :class="resultStatusToTextColor(result.status)">
                                <span class="font-semibold">{{
                                    $t(`enums.qualifications.ResultStatus.${ResultStatus[result.status]}`)
                                }}</span>
                            </span>
                        </template>
                    </template>

                    <template #score-data="{ row: result }">
                        <template v-if="result.score">{{ $n(result.score) }}</template>
                    </template>

                    <template #summary-data="{ row: result }">
                        <p v-if="result.summary" class="text-sm">
                            {{ result.summary }}
                        </p>
                    </template>

                    <template #createdAt-data="{ row: result }">
                        <GenericTime :value="result.createdAt" />
                    </template>

                    <template #creator-data="{ row: result }">
                        <CitizenInfoPopover v-if="result.creator" :user="result.creator" />
                    </template>

                    <template #actions-data="{ row: result }">
                        <div :key="result.id">
                            <UTooltip v-if="result.status === ResultStatus.PENDING" :text="$t('common.grade')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-star"
                                    color="amber"
                                    @click="
                                        modal.open(ExamViewResultModal, {
                                            qualificationId: result.qualificationId,
                                            userId: result.userId,
                                            resultId: result.id,
                                            examMode: examMode,
                                            onRefresh: onRefresh,
                                        })
                                    "
                                />
                            </UTooltip>

                            <UTooltip v-if="examMode > QualificationExamMode.DISABLED" :text="$t('common.show')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-star"
                                    color="amber"
                                    @click="
                                        modal.open(ExamViewResultModal, {
                                            qualificationId: result.qualificationId,
                                            userId: result.userId,
                                            resultId: result.id,
                                            examMode: examMode,
                                            viewOnly: true,
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
                                    class="flex-initial"
                                    variant="link"
                                    icon="i-mdi-delete"
                                    color="error"
                                    @click="
                                        modal.open(ConfirmModal, {
                                            confirm: async () => deleteQualificationResult(result.id),
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
