<script lang="ts" setup>
import ExamViewQuestions from '~/components/qualifications/exam/ExamViewQuestions.vue';
import QualificationResultTutorForm from '~/components/qualifications/tutor/QualificationResultTutorForm.vue';
import type { GetUserExamResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualificationId: string;
        userId: number;
        resultId?: string;
        viewOnly?: boolean;
    }>(),
    {
        resultId: undefined,
        viewOnly: false,
    },
);

defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const { data } = useLazyAsyncData(`qualification-${props.qualificationId}-result-examinfo-${props.userId}`, () =>
    getUserExam(),
);

async function getUserExam(): Promise<GetUserExamResponse> {
    const call = getGRPCQualificationsClient().getUserExam({
        qualificationId: props.qualificationId,
        userId: props.userId,
    });
    const { response } = await call;

    totalQuestions.value = response.exam?.questions.filter((q) => q.data?.data.oneofKind !== 'separator').length ?? 0;

    return response;
}

const totalQuestions = ref(0);
const correctCount = ref(0);
const score = computed(() => (100 / totalQuestions.value) * correctCount.value);

function updateCount(add: boolean): void {
    if (add) {
        correctCount.value++;
    } else {
        correctCount.value--;
    }
}
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <QualificationResultTutorForm
            :qualification-id="qualificationId"
            :user-id="userId"
            :result-id="resultId"
            :score="score"
            :view-only="viewOnly"
            @refresh="$emit('refresh')"
            @close="isOpen = false"
        >
            <template v-if="data" #default>
                <ExamViewQuestions
                    v-if="data?.exam && data?.examUser && data?.responses"
                    :qualification-id="qualificationId"
                    :exam="data.exam"
                    :exam-user="data.examUser"
                    :responses="data.responses"
                >
                    <template #question-after="{ question }">
                        <UCheckbox
                            v-if="question.question.data?.data.oneofKind !== 'separator'"
                            :label="$t('components.qualifications.correct_question')"
                            @update:model-value="updateCount($event)"
                        />
                    </template>
                </ExamViewQuestions>

                <div v-if="!viewOnly" class="flex flex-1 justify-end p-2">
                    <p class="text-sm">
                        <span class="font-semibold">{{ $t('components.qualifications.correct_question') }}</span
                        >: {{ correctCount }} / {{ totalQuestions }} {{ $t('common.question', 2) }}
                    </p>
                </div>

                <UDivider class="mb-4 mt-2" />
            </template>
        </QualificationResultTutorForm>
    </UModal>
</template>
