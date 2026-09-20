<script lang="ts" setup>
import type { BadgeProps, FormSubmitEvent } from '@nuxt/ui';
import { isPast } from 'date-fns';
import { z } from 'zod';
import {
    areExamChoicesAllowed,
    areExamChoicesUnique,
    examTextLength,
    isAnsweredExamSingleChoice,
    isExamChoiceLimitExceeded,
    isValidExamResponseKind,
} from '~/utils/qualificationExam';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import ScrollToTop from '~/components/partials/ScrollToTop.vue';
import { authKeys } from '~/composables/useAuth';
import { getQualificationsExamClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ExamQuestions, ExamResponse, ExamResponses, ExamUser } from '~~/gen/ts/resources/qualifications/exam/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { SubmitExamResponse } from '~~/gen/ts/services/qualifications/exam';
import ExamQuestionNavigator from './ExamQuestionNavigator.vue';
import ExamViewQuestion from './ExamViewQuestion.vue';

const props = defineProps<{
    qualificationId: number;
    exam: ExamQuestions;
    examUser: ExamUser;
    examResponses?: ExamResponses;
    qualification?: QualificationShort;
    responses?: ExamResponses;
}>();

const emits = defineEmits<{
    (e: 'submit', response: SubmitExamResponse): void;
    (e: 'cancel'): void;
    (e: 'expired'): void;
}>();

const notifications = useNotificationsStore();

const { accountId, activeChar } = useAuth();

const formatDuration = useDurationFormatter();

const qualificationsExamClient = await getQualificationsExamClient();
const { t } = useI18n();
const submitExamModal = useOverlay().create(ConfirmModal);

const schema = z
    .object({
        responses: z.custom<ExamResponse>().array().max(100).default([]),
    })
    .superRefine((values, ctx) => {
        values.responses.forEach((examResponse, index) => {
            const question = props.exam.questions.find((item) => item.id === examResponse.questionId);
            const data = question?.data?.data;
            const response = examResponse.response?.response;
            if (!data || !response) {
                ctx.addIssue({
                    code: 'custom',
                    path: ['responses', index],
                    message: 'zod.custom.qualification_exam.invalid_response',
                });
                return;
            }

            if (!isValidExamResponseKind(data.oneofKind, response.oneofKind)) {
                ctx.addIssue({
                    code: 'custom',
                    path: ['responses', index],
                    message: 'zod.custom.qualification_exam.invalid_response',
                });
                return;
            }

            if (data.oneofKind === 'freeText' && response.oneofKind === 'freeText') {
                const length = examTextLength(response.freeText.text);
                if (data.freeText.minLength > 0 && length < data.freeText.minLength) {
                    ctx.addIssue({
                        code: 'custom',
                        path: ['responses', index],
                        message: 'zod.custom.qualification_exam.answer_too_short',
                    });
                }
                if (data.freeText.maxLength > 0 && length > data.freeText.maxLength) {
                    ctx.addIssue({
                        code: 'custom',
                        path: ['responses', index],
                        message: 'zod.custom.qualification_exam.answer_too_long',
                    });
                }
            }

            if (data.oneofKind === 'singleChoice' && response.oneofKind === 'singleChoice') {
                if (!areExamChoicesAllowed([response.singleChoice.choice], data.singleChoice.choices)) {
                    ctx.addIssue({
                        code: 'custom',
                        path: ['responses', index],
                        message: 'zod.custom.qualification_exam.select_answer',
                    });
                }
            }

            if (data.oneofKind === 'multipleChoice' && response.oneofKind === 'multipleChoice') {
                const choices = response.multipleChoice.choices;
                if (!areExamChoicesUnique(choices) || !areExamChoicesAllowed(choices, data.multipleChoice.choices)) {
                    ctx.addIssue({
                        code: 'custom',
                        path: ['responses', index],
                        message: 'zod.custom.qualification_exam.select_valid_answers',
                    });
                }
                if (isExamChoiceLimitExceeded(choices, data.multipleChoice.limit)) {
                    ctx.addIssue({
                        code: 'custom',
                        path: ['responses', index],
                        message: 'zod.custom.qualification_exam.too_many_answers',
                    });
                }
            }
        });
    });

type Schema = z.output<typeof schema>;

const disabled = ref<boolean>(false);

const endsAtTime = toDate(props.examUser.endsAt).getTime();
const startsAtTime = toDate(props.examUser.startedAt).getTime();
const timeLowAtTime = endsAtTime - (endsAtTime - startsAtTime) * 0.15;
const now = useNow({ interval: 1_000 });

const remainingSeconds = computed(() => Math.max(0, Math.ceil((endsAtTime - now.value.getTime()) / 1_000)));
const remainingTimeLabel = computed(
    () => `${Math.floor(remainingSeconds.value / 60)}:${String(remainingSeconds.value % 60).padStart(2, '0')}`,
);
const timerAnnouncement = computed(() => {
    if (remainingSeconds.value <= 5) return t('components.qualifications.exam_view.timer.announcement_5');
    if (remainingSeconds.value <= 10) return t('components.qualifications.exam_view.timer.announcement_10');
    if (remainingSeconds.value <= 30) return t('components.qualifications.exam_view.timer.announcement_30');
    return t('components.qualifications.exam_view.timer.announcement_60');
});
const remainingTimeColor = computed<BadgeProps['color']>(() => {
    if (remainingSeconds.value <= 60) return 'error';
    if (now.value.getTime() >= timeLowAtTime) return 'warning';
    return 'success';
});

const responseStorageKey = `qualifications-exam-responses-${authKeys.character(accountId.value, activeChar.value?.userId)}-${props.qualificationId}-${toDate(props.examUser.startedAt).getTime()}`;
const storedResponses = useLocalStorage<ExamResponse[]>(responseStorageKey, []);
const storedYesNoAnswered = useLocalStorage<number[]>(`${responseStorageKey}-yesno`, []);
const storedFlaggedQuestions = useLocalStorage<number[]>(`${responseStorageKey}-flags`, []);
const restoredFromLocal = ref(false);
const localServerConflict = ref(false);

const state = useState<Schema>(responseStorageKey, () => ({
    responses: props.examResponses?.responses ?? storedResponses.value,
}));

const yesNoAnswered = reactive(new Set<number>(storedYesNoAnswered.value));
const flaggedQuestionIds = reactive(new Set<number>(storedFlaggedQuestions.value));

function toggleFlag(questionId: number): void {
    if (flaggedQuestionIds.has(questionId)) flaggedQuestionIds.delete(questionId);
    else flaggedQuestionIds.add(questionId);
}

function responsesFingerprint(responses: ExamResponse[]): string {
    return JSON.stringify(
        partialResponses(responses)
            .slice()
            .sort((left, right) => left.questionId - right.questionId),
    );
}

function mergeStoredResponses(serverResponses: ExamResponse[]): ExamResponse[] {
    const localByQuestionId = new Map(storedResponses.value.map((response) => [response.questionId, response]));
    const merged = serverResponses.map((response) => localByQuestionId.get(response.questionId) ?? response);

    for (const response of serverResponses) {
        localByQuestionId.delete(response.questionId);
    }

    return [...merged, ...localByQuestionId.values()];
}

const answerableQuestions = computed(() =>
    props.exam.questions.filter((question) => !['separator', 'image'].includes(question.data?.data.oneofKind ?? '')),
);

function isAnswered(question: (typeof props.exam.questions)[number]): boolean {
    const response = state.value.responses.find((item) => item.questionId === question.id)?.response?.response;

    switch (question.data?.data.oneofKind) {
        case 'yesno':
            return yesNoAnswered.has(question.id);
        case 'freeText':
            return response?.oneofKind === 'freeText' && response.freeText.text.trim().length > 0;
        case 'singleChoice':
            return response?.oneofKind === 'singleChoice'
                ? isAnsweredExamSingleChoice(response.singleChoice.choice, question.data.data.singleChoice.choices)
                : false;
        case 'multipleChoice':
            return response?.oneofKind === 'multipleChoice' && response.multipleChoice.choices.length > 0;
        default:
            return false;
    }
}

const answeredCount = computed(() => answerableQuestions.value.filter(isAnswered).length);
const answeredQuestionIds = computed(
    () => new Set(answerableQuestions.value.filter(isAnswered).map((question) => question.id)),
);
const answerProgress = computed(() =>
    answerableQuestions.value.length > 0 ? (answeredCount.value / answerableQuestions.value.length) * 100 : 0,
);

type AutoSaveState = 'idle' | 'saving' | 'saved' | 'error';
const autoSaveState = ref<AutoSaveState>('idle');
const lastSavedAt = ref<Date>();
let lastSavedResponses = '';

const autoSaveColor = computed<BadgeProps['color']>(() => {
    if (autoSaveState.value === 'error') return 'error';
    if (autoSaveState.value === 'saving') return 'warning';
    return 'success';
});

const autoSaveLabel = computed(() => {
    if (autoSaveState.value === 'saving') return t('components.qualifications.exam_view.auto_save.saving');
    if (autoSaveState.value === 'error') return t('components.qualifications.exam_view.auto_save.error');
    if (autoSaveState.value === 'saved') return t('components.qualifications.exam_view.auto_save.saved');
    return t('components.qualifications.exam_view.auto_save.ready');
});

function questionAnchorId(questionId: number): string {
    return `question-${questionId}`;
}

const activeQuestionId = ref<string>();

function navigateToQuestion(questionId: number): void {
    const anchorId = questionAnchorId(questionId);
    activeQuestionId.value = anchorId;
    window.history.replaceState(null, '', `#${anchorId}`);

    const question = document.getElementById(anchorId);
    question?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    nextTick(() => question?.focus({ preventScroll: true }));
}

function focusQuestionFromHash(): void {
    const hash = window.location.hash.slice(1);
    activeQuestionId.value = hash.startsWith('question-') ? hash : undefined;
    if (!activeQuestionId.value) return;

    nextTick(() => document.getElementById(hash)?.focus({ preventScroll: true }));
}

onMounted(() => {
    focusQuestionFromHash();
    window.addEventListener('hashchange', focusQuestionFromHash);
});

onBeforeUnmount(() => window.removeEventListener('hashchange', focusQuestionFromHash));

let partialSubmitInFlight: Promise<SubmitExamResponse> | undefined;

function partialResponses(responses: ExamResponse[]): ExamResponse[] {
    return responses.filter((examResponse) => {
        const question = props.exam.questions.find((item) => item.id === examResponse.questionId);
        const data = question?.data?.data;
        const response = examResponse.response?.response;

        // State is deliberately initialized for every question so controls can
        // bind with v-model. Do not serialize an untouched single-choice
        // placeholder as an exam answer, though: it is neither a response nor
        // valid input for the server to persist.
        if (data?.oneofKind === 'singleChoice' && response?.oneofKind === 'singleChoice') {
            return isAnsweredExamSingleChoice(response.singleChoice.choice, data.singleChoice.choices);
        }

        return true;
    });
}

async function submitExam(values: Schema, partial: boolean = false): Promise<SubmitExamResponse> {
    if (partial && partialSubmitInFlight) {
        await partialSubmitInFlight.catch(() => undefined);
    }

    const submit = (async () => {
        try {
            if (partial) autoSaveState.value = 'saving';

            const call = qualificationsExamClient.submitExam({
                qualificationId: props.qualificationId,
                responses: {
                    qualificationId: props.qualificationId,
                    userId: 0,
                    attemptId: '',
                    responses: partial ? partialResponses(values.responses) : values.responses,
                },
                partial: partial,
            });
            const { response } = await call;

            if (partial) {
                lastSavedResponses = JSON.stringify(partialResponses(values.responses));
                autoSaveState.value = 'saved';
                lastSavedAt.value = new Date();
                return response;
            }

            state.value.responses = [];
            storedResponses.value = [];
            storedYesNoAnswered.value = [];
            storedFlaggedQuestions.value = [];

            emits('submit', response);

            return response;
        } catch (e) {
            if (partial) autoSaveState.value = 'error';
            handleGRPCError(e as RpcError);
            throw e;
        }
    })();

    if (!partial) return submit;

    partialSubmitInFlight = submit;
    try {
        return await submit;
    } finally {
        if (partialSubmitInFlight === submit) partialSubmitInFlight = undefined;
    }
}

onBeforeMount(() => {
    props.exam.questions.forEach((q) => {
        // Question already in state? Skip it
        if (state.value.responses.find((r) => r.questionId === q.id)) {
            if (q.data?.data.oneofKind === 'yesno') yesNoAnswered.add(q.id);
            return;
        }

        switch (q.data?.data.oneofKind ?? 'separator') {
            case 'separator':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'separator',
                            separator: {},
                        },
                    },
                });
                break;

            case 'image':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'separator',
                            separator: {},
                        },
                    },
                });
                break;

            case 'yesno':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'yesno',
                            yesno: {
                                value: false,
                            },
                        },
                    },
                });
                break;

            case 'freeText':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'freeText',
                            freeText: {
                                text: '',
                            },
                        },
                    },
                });
                break;

            case 'singleChoice':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'singleChoice',
                            singleChoice: {
                                choice: '',
                            },
                        },
                    },
                });
                break;

            case 'multipleChoice':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
                    question: q,
                    response: {
                        response: {
                            oneofKind: 'multipleChoice',
                            multipleChoice: {
                                choices: [],
                            },
                        },
                    },
                });
                break;
        }
    });
});

// Auto-save every 30 seconds if there are changes
const { pause: pauseAutoSave } = useIntervalFn(() => {
    const responses = partialResponses(state.value.responses);
    if (JSON.stringify(responses) === lastSavedResponses) return;
    return submitExam(state.value, true);
}, 30_000);

if (!props.responses) {
    let timeLowNotificationSent = false;
    const { pause } = useIntervalFn(async () => {
        if (isPast(endsAtTime)) {
            pauseAutoSave();
            pause();

            let finalSaveSucceeded = false;
            try {
                // The server accepts partial submissions during the grace period.
                await submitExam(state.value, true);
                finalSaveSucceeded = true;
            } catch {
                // Retry once after a short delay so a transient failure does not
                // lose the user's final answers.
                await new Promise((resolve) => setTimeout(resolve, 4000));
                try {
                    await submitExam(state.value, true);
                    finalSaveSucceeded = true;
                } catch {
                    // The RPC error has already been surfaced by submitExam.
                }
            }

            disabled.value = true;
            emits('expired');

            if (finalSaveSucceeded) {
                notifications.add({
                    title: { key: 'notifications.qualifications.times_up.title', parameters: {} },
                    description: { key: 'notifications.qualifications.times_up.content', parameters: {} },
                    type: NotificationType.INFO,
                });
            }
        } else if (!timeLowNotificationSent && Date.now() >= timeLowAtTime) {
            notifications.add({
                title: { key: 'notifications.qualifications.time_low.title', parameters: {} },
                description: { key: 'notifications.qualifications.time_low.content', parameters: {} },
                type: NotificationType.INFO,
            });
            timeLowNotificationSent = true;
        }
    }, 1000);
}

function setResponses(): void {
    if (!props.responses) {
        disabled.value = false;
        if (props.examResponses && storedResponses.value.length > 0) {
            localServerConflict.value =
                responsesFingerprint(storedResponses.value) !== responsesFingerprint(props.examResponses.responses);
            state.value.responses = mergeStoredResponses(props.examResponses.responses);
            restoredFromLocal.value = true;
        }
        if (!props.examResponses && state.value.responses.length === 0 && storedResponses.value.length > 0) {
            state.value.responses = storedResponses.value;
            restoredFromLocal.value = true;
        } else if (storedResponses.value.length > 0 && !props.examResponses) {
            restoredFromLocal.value = true;
        }
        storedYesNoAnswered.value.forEach((questionId) => yesNoAnswered.add(questionId));
        return;
    }

    disabled.value = true;
    state.value.responses = mergeStoredResponses(props.responses.responses);
    storedYesNoAnswered.value.forEach((questionId) => yesNoAnswered.add(questionId));
}

setResponses();

watch(
    () => props.responses,
    () => setResponses(),
);

watch(
    [state, () => [...yesNoAnswered], () => [...flaggedQuestionIds]],
    () => {
        if (!disabled.value) {
            storedResponses.value = partialResponses(state.value.responses);
            storedYesNoAnswered.value = [...yesNoAnswered];
            storedFlaggedQuestions.value = [...flaggedQuestionIds];
        }
    },
    { deep: true },
);

const formRef = useTemplateRef('formRef');

const containerRef = useTemplateRef('containerRef');

function openSubmitConfirmation(): void {
    const unansweredCount = answerableQuestions.value.length - answeredCount.value;
    const incomplete = unansweredCount > 0;

    submitExamModal.open({
        title: t(
            incomplete
                ? 'components.qualifications.exam_view.submit.incomplete_title'
                : 'components.qualifications.exam_view.submit.title',
        ),
        description: t(
            incomplete
                ? 'components.qualifications.exam_view.submit.incomplete_description'
                : 'components.qualifications.exam_view.submit.description',
            {
                answered: answeredCount.value,
                total: answerableQuestions.value.length,
                unanswered: unansweredCount,
                flagged: flaggedQuestionIds.size,
            },
        ),
        confirm: () => formRef.value?.submit(),
    });
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await submitExam(event.data, false).then(() => {
        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    });
});
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'overflow-y-hidden gap-0 sm:gap-0 p-0 sm:p-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.id.exam.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UFieldGroup>
                        <UButton
                            color="error"
                            variant="outline"
                            :disabled="!canSubmit"
                            :label="$t('common.cancel')"
                            @click="$emit('cancel')"
                        />

                        <UButton
                            type="submit"
                            icon="i-mdi-content-save"
                            :disabled="!canSubmit"
                            :loading="isSubmitting"
                            :label="$t('common.submit')"
                            @click="openSubmitConfirmation"
                        />
                    </UFieldGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="!responses">
                <template v-if="qualification" #default>
                    <div class="flex flex-1 flex-row flex-wrap items-center justify-between gap-2">
                        <div class="flex-1">
                            <h1 class="px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                                <template v-if="qualification.abbreviation">{{ qualification.abbreviation }}: </template>
                                {{ !qualification.title ? $t('common.untitled') : qualification.title }}
                            </h1>

                            <p v-if="qualification.description" class="px-0.5 py-1 text-base font-bold break-words sm:pl-1">
                                {{ qualification.description }}
                            </p>
                        </div>

                        <div class="flex flex-col items-center gap-2">
                            <div class="inline-flex items-center gap-1">
                                <UIcon class="size-6" name="i-mdi-clock" />

                                <span>{{ $t('common.time_remaining') }}:</span>
                            </div>

                            <UBadge
                                :color="remainingTimeColor"
                                class="font-semibold tabular-nums"
                                :label="remainingTimeLabel"
                            />
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar v-if="!responses">
                <template #left>
                    <div class="flex justify-between gap-2">
                        <div class="flex gap-2">
                            <UBadge
                                v-if="qualification?.examSettings?.time"
                                class="inline-flex gap-1"
                                icon="i-mdi-clock"
                                :label="`${$t('common.duration')}: ${formatDuration(qualification.examSettings.time)}`"
                            />
                            <UBadge
                                class="inline-flex gap-1"
                                icon="i-mdi-question-mark"
                                :label="`${$t('common.count')}: ${exam.questions.length} ${$t('common.question', exam.questions.length)}`"
                            />
                            <UBadge
                                class="inline-flex gap-1"
                                icon="i-mdi-progress-check"
                                :label="`${$t('components.qualifications.exam_view.progress')}: ${answeredCount}/${answerableQuestions.length}`"
                            />
                            <UBadge
                                class="inline-flex gap-1"
                                :color="autoSaveColor"
                                :icon="
                                    autoSaveState === 'error'
                                        ? 'i-mdi-alert-circle-outline'
                                        : autoSaveState === 'saving'
                                          ? 'i-mdi-content-save-edit-outline'
                                          : 'i-mdi-content-save-check-outline'
                                "
                                :label="autoSaveLabel"
                                :title="
                                    lastSavedAt
                                        ? $t('components.qualifications.exam_view.auto_save.last_saved', {
                                              time: $d(lastSavedAt, 'long'),
                                          })
                                        : undefined
                                "
                            />
                        </div>
                    </div>
                </template>

                <template #right>
                    <div class="flex gap-1">
                        <UBadge v-if="props.examUser.startedAt" class="inline-flex gap-1">
                            <span class="font-semibold">{{ $t('common.begins_at') }}:</span>
                            <span>{{ $d(toDate(props.examUser.startedAt), 'long') }}</span>
                        </UBadge>

                        <UBadge v-if="props.examUser.endsAt" class="inline-flex gap-1">
                            <span class="font-semibold">{{ $t('common.ends_at') }}:</span>
                            <span>{{ $d(toDate(props.examUser.endsAt), 'long') }}</span>
                        </UBadge>
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="relative h-full min-h-0">
                <ExamQuestionNavigator
                    :questions="exam.questions"
                    :answered-ids="answeredQuestionIds"
                    :flagged-ids="flaggedQuestionIds"
                    @navigate="navigateToQuestion"
                />

                <div ref="containerRef" class="h-full overflow-y-auto p-4 sm:gap-6 sm:p-6 lg:pr-72">
                    <UForm ref="formRef" :schema="schema" :state="state" @submit="submit">
                        <UContainer>
                            <UCard class="sticky top-0 z-10 mb-4 rounded-lg bg-elevated shadow-sm" :ui="{ body: 'p-4 sm:p-4' }">
                                <UAlert
                                    v-if="remainingSeconds <= 60"
                                    class="mb-2"
                                    color="error"
                                    variant="subtle"
                                    icon="i-mdi-clock-alert-outline"
                                    :title="$t('components.qualifications.exam_view.timer.final_warning_title')"
                                    :description="
                                        $t('components.qualifications.exam_view.timer.final_warning_description', {
                                            time: remainingTimeLabel,
                                        })
                                    "
                                    aria-hidden="true"
                                />

                                <span v-if="remainingSeconds <= 60" class="sr-only" aria-live="polite">
                                    {{ timerAnnouncement }}
                                </span>

                                <UAlert
                                    v-if="restoredFromLocal"
                                    class="mb-2"
                                    color="info"
                                    variant="subtle"
                                    icon="i-mdi-history"
                                    :title="$t('components.qualifications.exam_view.restored.title')"
                                    :description="$t('components.qualifications.exam_view.restored.description')"
                                    :close="{ onClick: () => (restoredFromLocal = false) }"
                                />

                                <UAlert
                                    v-if="localServerConflict"
                                    class="mb-2"
                                    color="warning"
                                    variant="subtle"
                                    icon="i-mdi-alert-outline"
                                    :title="$t('components.qualifications.exam_view.local_conflict.title')"
                                    :description="$t('components.qualifications.exam_view.local_conflict.description')"
                                    :close="{ onClick: () => (localServerConflict = false) }"
                                />

                                <div class="mb-1 flex justify-between text-sm" aria-live="polite">
                                    <span>{{ $t('components.qualifications.exam_view.progress') }}</span>
                                    <span>{{ answeredCount }} / {{ answerableQuestions.length }}</span>
                                </div>
                                <UProgress :model-value="answerProgress" />
                            </UCard>

                            <div class="flex flex-col gap-4">
                                <ExamViewQuestion
                                    v-for="(question, idx) in exam.questions"
                                    :id="questionAnchorId(question.id)"
                                    :key="question.id"
                                    v-model="state.responses[idx]"
                                    class="scroll-mt-48"
                                    :class="
                                        activeQuestionId === questionAnchorId(question.id)
                                            ? 'ring-2 !ring-primary-500'
                                            : undefined
                                    "
                                    tabindex="-1"
                                    :disabled="disabled"
                                    :yes-no-answered="yesNoAnswered.has(question.id)"
                                    :flagged="flaggedQuestionIds.has(question.id)"
                                    @yes-no-answered="yesNoAnswered.add(question.id)"
                                    @toggle-flag="toggleFlag(question.id)"
                                >
                                    <template #question-after>
                                        <slot name="question-after" :question="{ question }" />
                                    </template>
                                </ExamViewQuestion>
                            </div>

                            <UCard v-if="!disabled" class="mt-4">
                                <UButton
                                    class="w-full"
                                    type="button"
                                    icon="i-mdi-content-save"
                                    block
                                    :disabled="!canSubmit"
                                    :loading="isSubmitting"
                                    :label="$t('common.submit')"
                                    @click="openSubmitConfirmation"
                                />
                            </UCard>

                            <ScrollToTop :element="containerRef" />
                        </UContainer>
                    </UForm>
                </div>
            </div>
        </template>
    </UDashboardPanel>
</template>
