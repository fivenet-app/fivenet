<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { UForm } from '#components';
import type { ExamQuestions, ExamResponse, ExamResponses } from '~~/gen/ts/resources/qualifications/exam';
import type { SubmitExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ExamUser } from '~~/gen/ts/resources/qualifications/exam';
import ExamViewQuestion from './ExamViewQuestion.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    qualificationId: string;
    exam: ExamQuestions;
    examUser: ExamUser;
    qualification?: QualificationShort;
    responses?: ExamResponses;
}>();

const notifications = useNotificatorStore();

const schema = z.object({
    responses: z.custom<ExamResponse>().array().max(50),
});

type Schema = z.output<typeof schema>;

const disabled = ref(false);

const endsAtTime = toDate(props.examUser.endsAt).getTime();
const now = new Date().getTime();

const state = useState<Schema>('qualifications-exam-responses', () => ({
    responses: [],
}));

async function submitExam(values: Schema): Promise<SubmitExamResponse> {
    try {
        const call = getGRPCQualificationsClient().submitExam({
            qualificationId: props.qualificationId,
            responses: {
                qualificationId: props.qualificationId,
                userId: 0,
                responses: values.responses,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
        switch (q.data?.data.oneofKind ?? 'separator') {
            case 'separator':
                state.value.responses.push({
                    questionId: q.id,
                    userId: 0,
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

const form = ref<InstanceType<typeof UForm> | null>(null);

if (!props.responses) {
    useIntervalFn(async () => {
        if (endsAtTime - now < 0) {
            await form.value?.submit();

            // TODO send notification that the time is over

            await navigateTo({
                name: 'qualifications-id',
                params: { id: props.qualificationId },
            });
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
    <UDashboardToolbar>
        <template v-if="qualification" #default>
            <div class="mb-4 flex flex-1 justify-between gap-2">
                <div class="">
                    <h1 class="break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </h1>

                    <p v-if="qualification.description" class="break-words px-0.5 py-1 text-base font-bold sm:pl-1">
                        {{ qualification.description }}
                    </p>
                </div>

                <div class="inline-flex flex-col items-center gap-2">
                    <UIcon name="i-mdi-clock" class="size-10" />
                    <span class="font-semibold">
                        {{
                            useLocaleTimeAgo(toDate(props.examUser.endsAt), {
                                showSecond: true,
                                updateInterval: 1_000,
                            }).value
                        }}
                    </span>
                    <UBadge v-if="props.examUser.startedAt">
                        {{ $t('common.begins_at') }}
                        {{ $d(toDate(props.examUser.startedAt), 'long') }}
                    </UBadge>
                    <UBadge v-if="props.examUser.endsAt">
                        {{ $t('common.ends_at') }}
                        {{ $d(toDate(props.examUser.endsAt), 'long') }}
                    </UBadge>
                </div>
            </div>
        </template>
    </UDashboardToolbar>

    <UForm ref="form" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UCard>
            <UContainer>
                <div class="flex flex-col gap-4">
                    <ExamViewQuestion
                        v-for="(question, idx) in exam.questions"
                        v-model="state.responses[idx]"
                        :key="question.id"
                        :question="question"
                        :disabled="disabled"
                    />
                </div>
            </UContainer>

            <template v-if="!disabled" #footer>
                <UButton
                    type="submit"
                    icon="i-mdi-save-content"
                    block
                    class="w-full"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                >
                    {{ $t('common.submit') }}
                </UButton>
            </template>
        </UCard>
    </UForm>
</template>
