<script lang="ts" setup>
import type { Form, FormSubmitEvent } from '#ui/types';
import { differenceInMinutes, isPast } from 'date-fns';
import { z } from 'zod';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ExamQuestions, ExamResponse, ExamResponses, ExamUser } from '~~/gen/ts/resources/qualifications/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { SubmitExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import ExamViewQuestion from './ExamViewQuestion.vue';

const props = defineProps<{
    qualificationId: number;
    exam: ExamQuestions;
    examUser: ExamUser;
    qualification?: QualificationShort;
    responses?: ExamResponses;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const schema = z.object({
    responses: z.custom<ExamResponse>().array().max(50).default([]),
});

type Schema = z.output<typeof schema>;

const disabled = ref(false);

const endsAtTime = toDate(props.examUser.endsAt).getTime();

const state = useState<Schema>('qualifications-exam-responses', () => ({
    responses: [],
}));

async function submitExam(values: Schema): Promise<SubmitExamResponse> {
    try {
        const call = $grpc.qualifications.qualifications.submitExam({
            qualificationId: props.qualificationId,
            responses: {
                qualificationId: props.qualificationId,
                userId: 0,
                responses: values.responses,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        state.value.responses = [];
        await navigateTo({
            name: 'qualifications-id',
            params: { id: props.qualificationId },
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

onBeforeMount(() =>
    props.exam.questions.forEach((q) => {
        // Question already in state? Skip it
        if (state.value.responses.find((r) => r.questionId === q.id)) {
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
    }),
);

const form = ref<Form<Schema> | null>(null);

if (!props.responses) {
    let timeLowNotificationSent = false;
    useIntervalFn(async () => {
        const minutesLeft = differenceInMinutes(endsAtTime, new Date());
        if (isPast(endsAtTime)) {
            await form.value?.submit();

            notifications.add({
                title: { key: 'notifications.qualifications.times_up.title', parameters: {} },
                description: { key: 'notifications.qualifications.times_up.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });

            await navigateTo({
                name: 'qualifications-id',
                params: { id: props.qualificationId },
            });
        } else if (!timeLowNotificationSent && minutesLeft <= 4) {
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await submitExam(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardToolbar v-if="!responses">
        <template v-if="qualification" #default>
            <div class="mb-2 flex flex-1 flex-col justify-between gap-1">
                <div class="flex flex-1 flex-row justify-between gap-2">
                    <div>
                        <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                            <template v-if="qualification.abbreviation">{{ qualification.abbreviation }}: </template>
                            {{ !qualification.title ? $t('common.untitled') : qualification.title }}
                        </h1>

                        <p v-if="qualification.description" class="break-words px-0.5 py-1 text-base font-bold sm:pl-1">
                            {{ qualification.description }}
                        </p>
                    </div>

                    <div class="inline-flex flex-col items-end gap-2">
                        <UIcon class="size-8" name="i-mdi-clock" />

                        <span class="font-semibold">
                            {{
                                useLocaleTimeAgo(toDate(props.examUser.endsAt), {
                                    showSecond: true,
                                    updateInterval: 1_000,
                                }).value
                            }}
                        </span>
                    </div>
                </div>

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
            </div>
        </template>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <UForm ref="form" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ rounded: '' }">
                <UContainer>
                    <div class="flex flex-col gap-4">
                        <ExamViewQuestion
                            v-for="(question, idx) in exam.questions"
                            :key="question.id"
                            v-model="state.responses[idx]"
                            :disabled="disabled"
                        >
                            <template #question-after>
                                <slot name="question-after" :question="{ question }" />
                            </template>
                        </ExamViewQuestion>
                    </div>
                </UContainer>

                <template v-if="!disabled" #footer>
                    <UButton
                        class="w-full"
                        type="submit"
                        icon="i-mdi-save-content"
                        block
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                    >
                        {{ $t('common.submit') }}
                    </UButton>
                </template>
            </UCard>
        </UForm>
    </UDashboardPanelContent>
</template>
