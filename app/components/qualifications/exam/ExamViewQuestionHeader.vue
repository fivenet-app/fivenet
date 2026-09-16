<script lang="ts" setup>
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';

defineProps<{
    question: ExamQuestion;
}>();
</script>

<template>
    <div class="flex flex-1 flex-col gap-1">
        <div class="flex items-start justify-between gap-2">
            <h4 class="text-xl" :title="`${$t('common.id')}: ${question.id}`">
                {{ question.title }}
            </h4>
            <p v-if="question.points">{{ $t('common.point', question.points) }}</p>
        </div>
        <p class="text-muted">{{ question.description }}</p>
    </div>
</template>
