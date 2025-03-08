<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import QualificationResultTutorForm from '~/components/qualifications/tutor/QualificationResultTutorForm.vue';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/qualifications';
import type { GetUserExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import ExamViewResult from '../exam/ExamViewResult.vue';

const props = withDefaults(
    defineProps<{
        qualificationId: number;
        userId: number;
        resultId?: number;
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

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`qualification-${props.qualificationId}-result-examinfo-${props.userId}`, () => getUserExam());

async function getUserExam(): Promise<GetUserExamResponse> {
    const call = $grpc.qualifications.qualifications.getUserExam({
        qualificationId: props.qualificationId,
        userId: props.userId,
    });
    const { response } = await call;

    totalQuestions.value =
        response.responses?.responses?.filter((q) => q.question?.data?.data.oneofKind !== 'separator').length ?? 0;

    totalPoints.value = 0;
    response.responses?.responses
        .filter((q) => q.question?.data?.data.oneofKind !== 'separator')
        .forEach((q) => (totalPoints.value += q.question?.points ?? 0));

    // Make sure the grading responses list is "valid" for the questions/responses
    if (!response.grading) {
        response.grading = {
            responses: [],
        };
    }
    response.responses?.responses.forEach((q) => {
        // Check if there is already a grading response
        if (response.grading?.responses.find((r) => r.questionId === q.questionId)) {
            return;
        }

        response.grading?.responses.push({
            questionId: q.questionId,
            checked: false,
            points: q.question?.points ?? 0,
        });
    });

    return response;
}

function getGradingIndex(id: number): number {
    return data.value!.grading!.responses.findIndex((a) => a.questionId === id);
}

const totalPoints = ref(0);
const pointCount = computed(() => data.value?.grading?.responses.map((a) => a.points).reduce((sum, a) => sum + a, 0));

const totalQuestions = ref(0);
const correctCount = ref(0);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <QualificationResultTutorForm
            :qualification-id="qualificationId"
            :user-id="userId"
            :result-id="resultId"
            :score="pointCount"
            :view-only="viewOnly"
            :grading="data?.grading"
            @refresh="$emit('refresh')"
            @close="isOpen = false"
        >
            <template v-if="examMode >= QualificationExamMode.REQUEST_NEEDED" #default>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.exam')])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.exam')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!data" :type="$t('common.exam')" icon="i-mdi-sigma" />

                <ExamViewResult v-else-if="data.responses" :qualification-id="qualificationId" :responses="data.responses">
                    <template #question-after="{ question }">
                        <div
                            v-if="
                                question.question.question?.data?.data.oneofKind !== 'separator' &&
                                getGradingIndex(question.question.questionId) > -1
                            "
                            class="flex flex-col gap-2"
                        >
                            <div class="flex flex-col gap-2 md:flex-row">
                                <UFormGroup :label="$t('common.corrected')">
                                    <div class="flex flex-col md:items-center">
                                        <UCheckbox
                                            v-model="
                                                data.grading!.responses[getGradingIndex(question.question.questionId)]!.checked
                                            "
                                        />
                                    </div>
                                </UFormGroup>

                                <UFormGroup :label="$t('common.points', 2)">
                                    <UInput
                                        v-model="data.grading!.responses[getGradingIndex(question.question.questionId)]!.points"
                                        type="number"
                                        :step="0.5"
                                        :min="0"
                                        :max="question.question.question?.points"
                                        class="max-w-24"
                                    />
                                </UFormGroup>
                            </div>

                            <div v-if="question.question.question?.answer?.answerKey" class="inline-flex flex-col gap-2">
                                <p class="text-sm font-semibold">{{ $t('common.answer_key') }}:</p>
                                <p class="text-sm">{{ question.question.question?.answer?.answerKey }}</p>
                            </div>
                        </div>
                    </template>
                </ExamViewResult>

                <div class="flex flex-1 justify-end gap-2 p-2">
                    <p class="text-sm">
                        <span class="font-semibold">{{ $t('common.corrected') }}</span
                        >: {{ correctCount }} / {{ totalQuestions }} {{ $t('common.question', 2) }}
                    </p>
                    <p class="text-sm">
                        <span class="font-semibold">{{ $t('common.points', 2) }}</span
                        >: {{ pointCount }} / {{ totalPoints }} {{ $t('common.points', 2) }}
                    </p>
                </div>

                <UDivider v-if="!viewOnly" class="mb-4 mt-2" />
            </template>
        </QualificationResultTutorForm>
    </UModal>
</template>
