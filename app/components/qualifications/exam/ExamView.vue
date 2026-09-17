<script lang="ts" setup>
import { isPast } from 'date-fns';
import { emojiBlasts } from 'emoji-blast';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { getQualificationsExamClient } from '~~/gen/ts/clients';
import type { ExamQuestions, ExamResponses, ExamUser } from '~~/gen/ts/resources/qualifications/exam/exam';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { GetExamInfoResponse, TakeExamResponse } from '~~/gen/ts/services/qualifications/exam';
import ExamViewQuestions from './ExamViewQuestions.vue';

const props = defineProps<{
    qualificationId: number;
}>();

const overlay = useOverlay();

const { t } = useI18n();

const formatDuration = useDurationFormatter();

const qualificationsExamClient = await getQualificationsExamClient();

const { data, status, refresh, error } = useAuthedLazyAsyncData(
    'userState',
    `qualification-${props.qualificationId}-examinfo`,
    ({ signal }) => getExamInfo(props.qualificationId, signal),
);

async function getExamInfo(qualificationId: number, signal: AbortSignal): Promise<GetExamInfoResponse> {
    try {
        const call = qualificationsExamClient.getExamInfo(
            {
                qualificationId: qualificationId,
            },
            { abort: signal },
        );
        const { response } = await call;

        examUser.value = response.examUser;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function takeExam(cancel = false): Promise<TakeExamResponse> {
    try {
        const call = qualificationsExamClient.takeExam({
            qualificationId: props.qualificationId,
            cancel: cancel,
        });
        const { response } = await call;

        examResponses.value = response.responses;
        exam.value = response.exam;
        examUser.value = response.examUser;

        if (response.timesUp) {
            examExpired.value = true;
            exam.value = undefined;
            examResponses.value = undefined;
            examUser.value = undefined;
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const exam = ref<ExamQuestions | undefined>();
const examUser = ref<ExamUser | undefined>();
const examResponses = ref<ExamResponses | undefined>();
const examExpired = ref(false);

const submitted = computed(() => data.value?.examUser?.endedAt !== undefined);
const expired = computed(() => {
    const endsAt = data.value?.examUser?.endsAt;
    return examExpired.value || (!submitted.value && endsAt !== undefined && isPast(toDate(endsAt)));
});
const result = computed(() => data.value?.qualification?.result);
const celebratedResultId = ref<number>();

watch(
    result,
    (value) => {
        if (import.meta.server || value?.status !== ResultStatus.SUCCESSFUL || value.id === celebratedResultId.value) return;

        celebratedResultId.value = value.id;
        const { cancel } = emojiBlasts({
            emojis: ['🎉', '🥳', '🎊', '🏆', '✨', '🙌', '💯'],
        });
        useTimeoutFn(cancel, 5000);
    },
    { immediate: true },
);

async function cancelExam(): Promise<void> {
    await takeExam(true);
    exam.value = undefined;
    examResponses.value = undefined;
    examUser.value = undefined;
    await refresh();
}

async function handleExpired(): Promise<void> {
    examExpired.value = true;
    exam.value = undefined;
    examResponses.value = undefined;
    examUser.value = undefined;
    await refresh();
}

const cancelExamModal = overlay.create(ConfirmModal);

function openCancelExamConfirmation(): void {
    cancelExamModal.open({
        title: t('components.qualifications.exam_view.cancel.title'),
        description: t('components.qualifications.exam_view.cancel.description'),
        confirm: cancelExam,
    });
}

watch(data, async () => {
    if (
        !examExpired.value &&
        data.value?.examUser?.endsAt !== undefined &&
        data.value?.examUser?.endedAt === undefined &&
        !isPast(toDate(data.value.examUser.endsAt))
    ) {
        await takeExam(false);
    }
});
</script>

<template>
    <ExamViewQuestions
        v-if="data && !examExpired && exam && examUser && examUser?.endsAt && !examUser?.endedAt"
        :qualification-id="qualificationId"
        :exam="exam"
        :exam-user="examUser"
        :exam-responses="examResponses"
        :qualification="data.qualification"
        @cancel="openCancelExamConfirmation"
        @submit="
            () => {
                examUser = undefined;
                refresh();
            }
        "
        @expired="handleExpired"
    />

    <UDashboardPanel v-else :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.id.exam.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton :fallback-to="`/qualifications/${qualificationId}`" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="data">
                <template #left>
                    <div class="flex gap-2">
                        <UBadge
                            v-if="data?.qualification?.examSettings?.time"
                            class="inline-flex gap-1"
                            icon="i-mdi-clock"
                            :label="`${$t('common.duration')}: ${formatDuration(data.qualification.examSettings.time)}`"
                        />

                        <UBadge
                            class="inline-flex gap-1"
                            icon="i-mdi-question-mark"
                            :label="`${$t('common.count')}: ${data?.questionCount} ${$t('common.question', data?.questionCount ?? 1)}`"
                        />
                    </div>
                </template>

                <template #right>
                    <div class="flex gap-2">
                        <UBadge
                            v-if="data.examUser?.startedAt"
                            :label="`${$t('common.begins_at')} ${$d(toDate(data.examUser?.startedAt), 'long')}`"
                        />

                        <UBadge
                            v-if="data?.examUser?.endsAt"
                            :label="`${$t('common.ends_at')} ${$d(toDate(data?.examUser?.endsAt), 'long')}`"
                        />
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.exam', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.exam', 1)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data"
                icon="i-mdi-account-school"
                :message="$t('common.not_found', [$t('common.qualification', 1)])"
            />

            <template v-else>
                <UCard>
                    <template v-if="result?.status === ResultStatus.SUCCESSFUL">
                        <UAlert
                            color="success"
                            variant="subtle"
                            icon="i-mdi-party-popper"
                            :ui="{ icon: 'size-8', title: 'text-xl' }"
                        >
                            <template #title>
                                {{ $t('components.qualifications.exam_view.result_successful.title') }}
                            </template>
                            <template #description>
                                <div class="space-y-1">
                                    <p v-if="result.autoGraded">
                                        {{ $t('components.qualifications.exam_view.result_successful.description') }}
                                    </p>
                                    <p v-else>
                                        {{ $t('components.qualifications.exam_view.result_successful.manual_description') }}
                                    </p>
                                    <p class="font-semibold">
                                        {{ $t('common.score') }}: {{ $t('common.point', result.score ?? 0) }} 🏆
                                    </p>
                                    <p v-if="result.summary">{{ result.summary }}</p>
                                </div>
                            </template>
                        </UAlert>

                        <PartialsBackButton class="mt-4" block :to="`/qualifications/${qualificationId}`" />
                    </template>
                    <template v-else-if="result?.status === ResultStatus.FAILED">
                        <UAlert
                            :title="$t('components.qualifications.exam_view.result_failed.title')"
                            color="error"
                            variant="subtle"
                            icon="i-mdi-close-circle-outline"
                            :ui="{ icon: 'size-8', title: 'text-xl' }"
                        >
                            <template #description>
                                <div class="space-y-1">
                                    <p v-if="result.autoGraded">
                                        {{ $t('components.qualifications.exam_view.result_failed.description') }}
                                    </p>
                                    <p v-else>
                                        {{ $t('components.qualifications.exam_view.result_failed.manual_description') }}
                                    </p>
                                    <p class="font-semibold">
                                        {{ $t('common.score') }}:
                                        {{ $t('common.point', result.score ?? 0) }}
                                    </p>
                                    <p v-if="result.summary">{{ result.summary }}</p>
                                </div>
                            </template>
                        </UAlert>

                        <PartialsBackButton class="mt-4" block :to="`/qualifications/${qualificationId}`" />
                    </template>
                    <template v-else-if="submitted">
                        <UAlert
                            :title="$t('components.qualifications.exam_view.submitted.title')"
                            color="success"
                            variant="subtle"
                            icon="i-mdi-check-circle-outline"
                            :description="$t('components.qualifications.exam_view.submitted.description')"
                            :ui="{ icon: 'size-8', title: 'text-xl' }"
                        />

                        <PartialsBackButton class="mt-4" block :to="`/qualifications/${qualificationId}`" />
                    </template>
                    <template v-else-if="expired">
                        <UAlert
                            :title="$t('components.qualifications.exam_view.expired.title')"
                            color="warning"
                            variant="subtle"
                            icon="i-mdi-clock-alert-outline"
                            :description="$t('components.qualifications.exam_view.expired.description')"
                            :ui="{ icon: 'size-8', title: 'text-xl' }"
                        />

                        <PartialsBackButton class="mt-4" block :to="`/qualifications/${qualificationId}`" />
                    </template>

                    <UButton
                        v-else-if="!data?.examUser?.endedAt"
                        class="w-full"
                        size="xl"
                        color="neutral"
                        icon="i-mdi-play"
                        block
                        :label="$t('components.qualifications.take_test')"
                        @click="
                            () => {
                                takeExam(false);
                            }
                        "
                    />
                </UCard>
            </template>
        </template>
    </UDashboardPanel>
</template>
