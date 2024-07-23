<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type {
    DeleteQualificationResultResponse,
    ListQualificationsResultsResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import { resultStatusToTextColor } from '../helpers';
import ExamViewResultModal from './ExamViewResultModal.vue';

const props = withDefaults(
    defineProps<{
        qualificationId?: string;
        status?: ResultStatus[];
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
} = useLazyAsyncData(`qualifications-results-${page.value}-${props.qualificationId}`, () =>
    listQualificationsResults(props.qualificationId, props.status),
);

async function listQualificationsResults(
    qualificationId?: string,
    status?: ResultStatus[],
): Promise<ListQualificationsResultsResponse> {
    try {
        const call = getGRPCQualificationsClient().listQualificationsResults({
            pagination: {
                offset: offset.value,
            },
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

watch(offset, async () => refresh());

async function deleteQualificationResult(resultId: string): Promise<DeleteQualificationResultResponse> {
    try {
        const call = getGRPCQualificationsClient().deleteQualificationResult({
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
        key: 'citizen',
        label: t('common.citizen'),
    },
    {
        key: 'status',
        label: t('common.status'),
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
    },
    {
        key: 'creator',
        label: t('common.creator'),
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
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.qualifications', 2)])"
                :retry="refresh"
            />

            <template v-else>
                <UTable
                    :loading="loading"
                    :columns="columns"
                    :rows="data?.results"
                    :empty-state="{ icon: 'i-mdi-sigma', label: $t('common.not_found', [$t('common.result', 2)]) }"
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
                        <template v-if="result.score">{{ result.score }}</template>
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
                            <UButton
                                v-if="result.status === ResultStatus.PENDING"
                                variant="link"
                                icon="i-mdi-star"
                                color="amber"
                                @click="
                                    modal.open(ExamViewResultModal, {
                                        qualificationId: result.qualificationId,
                                        userId: result.userId,
                                        resultId: result.id,
                                        onRefresh: onRefresh,
                                    })
                                "
                            />

                            <UButton
                                v-if="can('QualificationsService.DeleteQualificationResult').value"
                                class="flex-initial"
                                variant="link"
                                icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => deleteQualificationResult(result.id),
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
