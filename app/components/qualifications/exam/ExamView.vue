<script lang="ts" setup>
import { isPast } from 'date-fns';
import { emojiBlasts } from 'emoji-blast';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DraftBadge from '~/components/partials/DraftBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { getQualificationsExamClient } from '~~/gen/ts/clients';
import {
    QualificationExamMode,
    type ExamQuestions,
    type ExamResponses,
    type ExamUser,
} from '~~/gen/ts/resources/qualifications/exam/exam';
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
const takingExam = ref(false);

const submitted = computed(() => data.value?.examUser?.endedAt !== undefined);
const expired = computed(() => {
    const endsAt = data.value?.examUser?.endsAt;
    return examExpired.value || (!submitted.value && endsAt !== undefined && isPast(toDate(endsAt)));
});
const result = computed(() => data.value?.qualification?.result);
const submittedAt = computed(() => {
    const endedAt = data.value?.examUser?.endedAt;
    return endedAt ? toDate(endedAt) : undefined;
});
const submissionDuration = computed(() => {
    const startedAt = data.value?.examUser?.startedAt;
    const endedAt = data.value?.examUser?.endedAt;
    if (!startedAt || !endedAt) return undefined;

    const elapsedSeconds = Math.max(0, Math.floor((toDate(endedAt).getTime() - toDate(startedAt).getTime()) / 1_000));
    const elapsedMinutes = Math.round(elapsedSeconds / 60);

    return formatDuration({ seconds: elapsedMinutes * 60, nanos: 0 }, 'minute');
});
const qualificationTitle = computed(() => {
    const qualification = data.value?.qualification;
    if (!qualification) return t('pages.qualifications.id.exam.title');

    return [qualification.abbreviation, qualification.title].filter(Boolean).join(': ') || t('common.untitled');
});
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

async function startOrResumeExam(): Promise<void> {
    if (takingExam.value) return;

    takingExam.value = true;
    try {
        await takeExam(false);
    } finally {
        takingExam.value = false;
    }
}

async function handleExpired(): Promise<void> {
    examExpired.value = true;
    exam.value = undefined;
    examResponses.value = undefined;
    examUser.value = undefined;
    await refresh();
}

const cancelExamModal = overlay.create(ConfirmModal);
const startExamModal = overlay.create(ConfirmModal);

function openCancelExamConfirmation(): void {
    cancelExamModal.open({
        title: t('components.qualifications.exam_view.cancel.title'),
        description: t('components.qualifications.exam_view.cancel.description'),
        confirm: cancelExam,
    });
}

function openStartExamConfirmation(): void {
    const isResume = data.value?.examUser?.startedAt !== undefined;

    startExamModal.open({
        title: t(
            isResume
                ? 'components.qualifications.exam_view.resume_confirm.title'
                : 'components.qualifications.exam_view.start_confirm.title',
        ),
        description: t(
            isResume
                ? 'components.qualifications.exam_view.resume_confirm.description'
                : 'components.qualifications.exam_view.start_confirm.description',
        ),
        color: isResume ? 'primary' : 'success',
        icon: isResume ? 'i-mdi-play' : 'i-mdi-play-circle-outline',
        iconClass: isResume ? 'text-primary-500 dark:text-primary-400' : 'text-green-500 dark:text-green-400',
        confirm: startOrResumeExam,
    });
}
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
                <div class="flex flex-1 flex-row flex-wrap items-center justify-between gap-2">
                    <div class="flex-1">
                        <h1 class="px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                            {{ qualificationTitle }}
                        </h1>
                        <p v-if="data.qualification?.description" class="px-0.5 py-1 text-base font-bold break-words sm:pl-1">
                            {{ data.qualification.description }}
                        </p>
                    </div>

                    <div class="flex flex-wrap gap-2">
                        <OpenClosedBadge :closed="data.qualification?.closed" />
                        <DraftBadge v-if="data.qualification?.draft" />
                        <UBadge
                            v-if="data.qualification?.public"
                            class="inline-flex gap-1"
                            icon="i-mdi-earth"
                            color="neutral"
                            :label="$t('common.public')"
                        />
                        <UBadge
                            v-if="data.qualification?.examMode"
                            class="inline-flex gap-1"
                            icon="i-mdi-test-tube"
                            :label="`${$t('common.exam', 1)}: ${$t(
                                `enums.qualifications.QualificationExamMode.${QualificationExamMode[data.qualification.examMode]}`,
                            )}`"
                        />
                    </div>
                </div>
            </UDashboardToolbar>

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
            <UContainer>
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
                    <UCard :ui="{ body: 'p-4 sm:p-4' }">
                        <div v-if="submittedAt || submissionDuration" class="mb-4 flex flex-wrap gap-2">
                            <UBadge v-if="submittedAt" icon="i-mdi-calendar-check" class="inline-flex gap-1">
                                <span class="font-semibold">{{ $t('common.ended_at') }}:</span>
                                <span>{{ $d(submittedAt, 'long') }}</span>
                            </UBadge>

                            <UBadge v-if="submissionDuration" icon="i-mdi-timer-outline" class="inline-flex gap-1">
                                <span class="font-semibold">{{ $t('common.duration') }}:</span>
                                <span>{{ submissionDuration }}</span>
                            </UBadge>
                        </div>

                        <template v-if="result?.status === ResultStatus.PENDING">
                            <UAlert
                                :title="$t('components.qualifications.exam_view.result_pending.title')"
                                color="warning"
                                variant="subtle"
                                icon="i-mdi-clock-outline"
                                :description="$t('components.qualifications.exam_view.result_pending.description')"
                                :ui="{ icon: 'size-8', title: 'text-xl' }"
                            />

                            <PartialsBackButton class="mt-4" block :to="`/qualifications/${qualificationId}`" />
                        </template>
                        <template v-else-if="result?.status === ResultStatus.SUCCESSFUL">
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
                            :loading="takingExam"
                            :disabled="takingExam"
                            :label="
                                $t(
                                    data?.examUser?.startedAt
                                        ? 'components.qualifications.exam_view.resume'
                                        : 'components.qualifications.exam_view.start',
                                )
                            "
                            @click="openStartExamConfirmation"
                        />
                    </UCard>
                </template>
            </UContainer>
        </template>
    </UDashboardPanel>
</template>
