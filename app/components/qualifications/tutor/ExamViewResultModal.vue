<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import ExamViewQuestions from '~/components/qualifications/exam/ExamViewQuestions.vue';
import QualificationResultTutorForm from '~/components/qualifications/tutor/QualificationResultTutorForm.vue';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/qualifications';
import type { GetUserExamResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        qualificationId: string;
        userId: number;
        resultId?: string;
        viewOnly?: boolean;
        examMode?: QualificationExamMode;
    }>(),
    {
        resultId: undefined,
        viewOnly: false,
        examMode: QualificationExamMode.DISABLED,
    },
);

defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}-result-examinfo-${props.userId}`, () => getUserExam());

async function getUserExam(): Promise<GetUserExamResponse> {
    const call = getGRPCQualificationsClient().getUserExam({
        qualificationId: props.qualificationId,
        userId: props.userId,
    });
    const { response } = await call;

    totalQuestions.value = response.exam?.questions.filter((q) => q.data?.data.oneofKind !== 'separator').length ?? 0;

    totalPoints.value = 0;
    response.exam?.questions
        .filter((q) => q.data?.data.oneofKind !== 'separator')
        .forEach((q) => (totalPoints.value += q.points ?? 0));

    return response;
}

const totalPoints = ref(0);
const pointCount = ref(0);

const totalQuestions = ref(0);
const correctCount = ref(0);

function updateCount(add: boolean, points?: number): void {
    if (add) {
        correctCount.value++;
        pointCount.value += points ?? 0;
    } else {
        correctCount.value--;
        pointCount.value -= points ?? 0;
    }
}
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <QualificationResultTutorForm
            :qualification-id="qualificationId"
            :user-id="userId"
            :result-id="resultId"
            :score="pointCount"
            :view-only="viewOnly"
            @refresh="$emit('refresh')"
            @close="isOpen = false"
        >
            <template v-if="examMode >= QualificationExamMode.REQUEST_NEEDED" #default>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.exam')])" />
                <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.exam')])" :retry="refresh" />
                <DataNoDataBlock v-else-if="!data" :type="$t('common.exam')" icon="i-mdi-sigma" />

                <ExamViewQuestions
                    v-else-if="data?.exam && data?.examUser && data?.responses"
                    :qualification-id="qualificationId"
                    :exam="data.exam"
                    :exam-user="data.examUser"
                    :responses="data.responses"
                >
                    <template #question-after="{ question }">
                        <div v-if="question.question.data?.data.oneofKind !== 'separator'" class="flex flex-col gap-2">
                            <UCheckbox
                                :label="$t('components.qualifications.correct_question')"
                                @update:model-value="updateCount($event, question.question.points)"
                            />

                            <div class="inline-flex flex-col gap-2">
                                <p class="text-sm font-semibold">{{ $t('common.answer_key') }}:</p>
                                <p class="text-sm">{{ question.question.answer?.answerKey ?? $t('common.na') }}</p>
                            </div>
                        </div>
                    </template>
                </ExamViewQuestions>

                <div v-if="!viewOnly" class="flex flex-1 justify-end gap-2 p-2">
                    <p class="text-sm">
                        <span class="font-semibold">{{ $t('components.qualifications.correct_question') }}</span
                        >: {{ correctCount }} / {{ totalQuestions }} {{ $t('common.question', 2) }}
                    </p>
                    <p class="text-sm">
                        <span class="font-semibold">{{ $t('common.points', 2) }}</span
                        >: {{ pointCount }} / {{ totalPoints }} {{ $t('common.points', 2) }}
                    </p>
                </div>

                <UDivider class="mb-4 mt-2" />
            </template>
        </QualificationResultTutorForm>
    </UModal>
</template>
