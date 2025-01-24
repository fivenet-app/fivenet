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
                    :disabled="true"
                >
                    <template #question-after>
                        <slot name="question-after" :question="{ question }" />
                    </template>
                </ExamViewQuestion>
            </div>
        </UContainer>
    </UCard>
</template>
