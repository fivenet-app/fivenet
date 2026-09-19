<script lang="ts" setup>
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';
import ExamQuestionFlagButton from './ExamQuestionFlagButton.vue';

const props = withDefaults(
    defineProps<{
        question: ExamQuestion;
        flagged?: boolean;
        disabled?: boolean;
        showFlag?: boolean;
    }>(),
    {
        flagged: false,
        disabled: false,
        showFlag: true,
    },
);

const emit = defineEmits<{
    (e: 'toggleFlag'): void;
}>();

const canFlag = computed(() => props.showFlag && !['separator', 'image'].includes(props.question.data?.data.oneofKind ?? ''));
</script>

<template>
    <div class="flex flex-1 flex-col gap-1">
        <div class="flex items-start justify-between gap-2">
            <h4 class="text-xl" :title="`${$t('common.id')}: ${question.id}`">
                {{ question.title }}
            </h4>
            <div class="flex shrink-0 items-center gap-2">
                <p v-if="question.points">{{ $t('common.point', question.points) }}</p>

                <UTooltip
                    v-if="canFlag"
                    :text="
                        flagged
                            ? $t('components.qualifications.exam_view.unflag')
                            : $t('components.qualifications.exam_view.flag')
                    "
                >
                    <ExamQuestionFlagButton :flagged="flagged" :disabled="disabled" @toggle="emit('toggleFlag')" />
                </UTooltip>
            </div>
        </div>

        <p v-if="question.description" class="text-muted">{{ question.description }}</p>
    </div>
</template>
