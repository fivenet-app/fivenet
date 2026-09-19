<script lang="ts" setup>
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';
import { createExamNavigatorEntries, getExamNavigatorTargetIndex } from './examNavigator';
import ExamQuestionNavigatorItem from './ExamQuestionNavigatorItem.vue';

const props = defineProps<{
    questions: ExamQuestion[];
    answeredIds: ReadonlySet<number>;
    flaggedIds: ReadonlySet<number>;
}>();

const emit = defineEmits<{
    (e: 'navigate', questionId: number): void;
    (e: 'keyboardNavigate', questionId: number): void;
}>();

const entries = computed(() => createExamNavigatorEntries(props.questions));

function navigate(questionId: number): void {
    emit('navigate', questionId);
}

function handleKeydown(event: KeyboardEvent, entry: (typeof entries.value)[number]): void {
    if (!['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(event.key)) return;

    const navigableEntries = entries.value.filter((item) => item.kind !== 'separator');
    const currentIndex = navigableEntries.findIndex((item) => item.question.id === entry.question.id);
    if (currentIndex < 0) return;

    const targetIndex = getExamNavigatorTargetIndex(event.key, currentIndex, navigableEntries.length);
    const target = navigableEntries[targetIndex];
    if (!target) return;

    event.preventDefault();
    emit('keyboardNavigate', target.question.id);
    nextTick(() => {
        const items = [...document.querySelectorAll<HTMLElement>(`[data-question-id="${target.question.id}"]`)];
        items.find((item) => item.getClientRects().length > 0)?.focus();
    });
}
</script>

<template>
    <template v-for="entry in entries" :key="entry.question.id">
        <USeparator v-if="entry.kind === 'separator'" class="my-2">
            <span class="max-w-full truncate text-xs font-semibold">
                {{ entry.question.title || $t('components.qualifications.exam_view.section') }}
            </span>
        </USeparator>

        <ExamQuestionNavigatorItem
            v-else
            :question="entry.question"
            :answerable="entry.answerable"
            :number="entry.number"
            :answered="props.answeredIds.has(entry.question.id)"
            :flagged="props.flaggedIds.has(entry.question.id)"
            @keydown="handleKeydown($event, entry)"
            @navigate="navigate"
        />
    </template>
</template>
