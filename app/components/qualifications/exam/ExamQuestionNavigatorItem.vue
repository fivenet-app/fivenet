<script lang="ts" setup>
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';

const props = defineProps<{
    question: ExamQuestion;
    answerable: boolean;
    number?: number;
    answered: boolean;
    flagged: boolean;
}>();

const emit = defineEmits<{
    (e: 'navigate', questionId: number): void;
}>();

const questionIcon = computed(() => {
    if (!props.answerable) {
        return props.question.data?.data.oneofKind === 'image' ? 'i-mdi-image-outline' : 'i-mdi-information-outline';
    }

    return props.answered ? 'i-mdi-check-circle' : 'i-mdi-circle-outline';
});

const questionTypeIcon = computed(() => {
    switch (props.question.data?.data.oneofKind) {
        case 'yesno':
            return 'i-mdi-toggle-switch-outline';
        case 'freeText':
            return 'i-mdi-text-box-outline';
        case 'singleChoice':
            return 'i-mdi-radiobox-marked';
        case 'multipleChoice':
            return 'i-mdi-checkbox-multiple-marked-outline';
        default:
            return undefined;
    }
});

const questionTypeLabel = computed(() => {
    switch (props.question.data?.data.oneofKind) {
        case 'yesno':
            return 'components.qualifications.exam_view.question_types.yesno';
        case 'freeText':
            return 'components.qualifications.exam_view.question_types.free_text';
        case 'singleChoice':
            return 'components.qualifications.exam_view.question_types.single_choice';
        case 'multipleChoice':
            return 'components.qualifications.exam_view.question_types.multiple_choice';
        default:
            return undefined;
    }
});

function navigate(): void {
    emit('navigate', props.question.id);
}
</script>

<template>
    <UButton
        type="button"
        class="w-full justify-start"
        :data-question-id="question.id"
        :variant="answerable ? 'soft' : 'ghost'"
        :color="answerable && answered ? 'success' : 'neutral'"
        :icon="questionIcon"
        @click="navigate"
    >
        <span class="truncate text-left">
            <template v-if="number">{{ number }}. </template>
            {{ question.title || $t('common.untitled') }}
        </span>
        <UTooltip v-if="questionTypeIcon && questionTypeLabel" :text="$t(questionTypeLabel)">
            <UIcon :name="questionTypeIcon" class="shrink-0 text-muted" aria-hidden="true" />
        </UTooltip>
        <UIcon v-if="flagged" name="i-mdi-flag" class="ml-auto shrink-0 text-warning" />
    </UButton>
</template>
