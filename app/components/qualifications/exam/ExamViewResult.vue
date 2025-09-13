<script lang="ts" setup>
import type { ExamResponses } from '~~/gen/ts/resources/qualifications/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import ExamViewQuestion from './ExamViewQuestion.vue';

defineProps<{
    qualificationId: number;
    qualification?: QualificationShort;
    responses?: ExamResponses;
}>();
</script>

<template>
    <UCard>
        <UContainer>
            <div class="flex flex-col gap-4">
                <ExamViewQuestion
                    v-for="(question, idx) in responses?.responses"
                    :key="question.questionId"
                    v-model="responses!.responses[idx]"
                    disabled
                >
                    <template #question-after="{ disabled }">
                        <slot name="question-after" :question="{ question }" :disabled="disabled" />
                    </template>

                    <template #question-below="{ disabled }">
                        <slot name="question-below" :question="{ question }" :disabled="disabled" />
                    </template>
                </ExamViewQuestion>
            </div>
        </UContainer>
    </UCard>
</template>
