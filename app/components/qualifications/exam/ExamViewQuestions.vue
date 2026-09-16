<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
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
import ScrollToTop from '~/components/partials/ScrollToTop.vue';
import { authKeys } from '~/composables/useAuth';
import { getQualificationsExamClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ExamQuestions, ExamResponse, ExamResponses, ExamUser } from '~~/gen/ts/resources/qualifications/exam/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { SubmitExamResponse } from '~~/gen/ts/services/qualifications/exam';
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

const state = useState<Schema>(
    `qualifications-exam-responses-${authKeys.character(accountId.value, activeChar.value?.userId)}-${props.qualificationId}-${toDate(props.examUser.startedAt).getTime()}`,
    () => ({
        responses: props.examResponses?.responses ?? [],
    }),
);

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

            if (partial) return response;

            state.value.responses = [];

            emits('submit', response);

            return response;
        } catch (e) {
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
        if (state.value.responses.find((r) => r.questionId === q.id)) return;

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
const { pause: pauseAutoSave } = useIntervalFn(() => submitExam(state.value, true), 30_000);

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
                // The RPC error has already been surfaced by submitExam.
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
        return;
    }

    disabled.value = true;
    state.value.responses = props.responses.responses;
}

setResponses();

watch(
    () => props.responses,
    () => setResponses(),
);

const formRef = useTemplateRef('formRef');

const containerRef = useTemplateRef('containerRef');

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await submitExam(event.data, false)
        .then(() => {
            notifications.add({
                title: { key: 'notifications.action_successful.title', parameters: {} },
                description: { key: 'notifications.action_successful.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        })
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'overflow-y-hidden gap-0 sm:gap-0 p-0 sm:p-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.id.exam.title')">
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
                            :loading="!canSubmit"
                            :label="$t('common.submit')"
                            @click="formRef?.submit()"
                        />
                    </UFieldGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="!responses">
                <template v-if="qualification" #default>
                    <div class="flex flex-1 flex-row items-center justify-between gap-2">
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

                            <div class="font-semibold">
                                {{
                                    useLocaleTimeAgo(toDate(props.examUser.endsAt), {
                                        showSecond: true,
                                        updateInterval: 1_000,
                                    }).value
                                }}
                            </div>
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar v-if="!responses">
                <template #left>
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

                <template #right>
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
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div ref="containerRef" class="gap-4 overflow-y-auto p-4 sm:gap-6 sm:p-6">
                <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                    <UContainer>
                        <div class="flex flex-col gap-4">
                            <UCard v-for="(question, idx) in exam.questions" :key="question.id">
                                <ExamViewQuestion :key="question.id" v-model="state.responses[idx]" :disabled="disabled">
                                    <template #question-after>
                                        <slot name="question-after" :question="{ question }" />
                                    </template>
                                </ExamViewQuestion>
                            </UCard>
                        </div>

                        <UCard v-if="!disabled" class="mt-4">
                            <UButton
                                class="w-full"
                                type="submit"
                                icon="i-mdi-content-save"
                                block
                                :disabled="!canSubmit"
                                :loading="!canSubmit"
                                :label="$t('common.submit')"
                            />
                        </UCard>

                        <ScrollToTop :element="containerRef" />
                    </UContainer>
                </UForm>
            </div>
        </template>
    </UDashboardPanel>
</template>
