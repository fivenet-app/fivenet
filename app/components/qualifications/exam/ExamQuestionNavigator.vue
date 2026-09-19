<script lang="ts" setup>
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';
import { createExamNavigatorEntries } from './examNavigator';
import ExamQuestionNavigatorList from './ExamQuestionNavigatorList.vue';
import ExamQuestionNavigatorSummary from './ExamQuestionNavigatorSummary.vue';

const props = defineProps<{
    questions: ExamQuestion[];
    answeredIds: ReadonlySet<number>;
    flaggedIds: ReadonlySet<number>;
}>();

const emit = defineEmits<{
    (e: 'navigate', questionId: number): void;
}>();

const open = ref(false);

const entries = computed(() => createExamNavigatorEntries(props.questions));

const unansweredQuestions = computed(() =>
    entries.value.filter((entry) => entry.answerable && !props.answeredIds.has(entry.question.id)),
);
const flaggedQuestions = computed(() => entries.value.filter((entry) => props.flaggedIds.has(entry.question.id)));
const answeredCount = computed(
    () => entries.value.filter((entry) => entry.answerable && props.answeredIds.has(entry.question.id)).length,
);

function navigateToFirstUnanswered(): void {
    const first = unansweredQuestions.value[0];
    if (first) navigate(first.question.id);
}

function navigateToFirstFlagged(): void {
    const first = flaggedQuestions.value[0];
    if (first) navigate(first.question.id);
}

function navigate(questionId: number): void {
    emit('navigate', questionId);
    open.value = false;
}

function navigateFromKeyboard(questionId: number): void {
    emit('navigate', questionId);
}
</script>

<template>
    <div class="absolute top-2 right-4 z-20 hidden max-h-[calc(100%-1rem)] w-72 lg:block">
        <UButton
            v-if="unansweredQuestions.length > 0"
            class="mb-2 w-full justify-between"
            color="warning"
            variant="soft"
            trailing-icon="i-mdi-arrow-right"
            :label="$t('components.qualifications.exam_view.review_unanswered', { count: unansweredQuestions.length })"
            @click="navigateToFirstUnanswered"
        />
        <UButton
            v-if="flaggedQuestions.length > 0"
            class="mb-2 w-full justify-between"
            color="warning"
            variant="soft"
            trailing-icon="i-mdi-flag"
            :label="$t('components.qualifications.exam_view.review_flagged', { count: flaggedQuestions.length })"
            @click="navigateToFirstFlagged"
        />

        <div class="flex flex-col gap-2 rounded-lg border border-default bg-default/95 p-2 shadow-sm backdrop-blur">
            <div class="px-2 text-sm font-semibold">
                {{ $t('components.qualifications.exam_view.question_navigator') }}
            </div>

            <ExamQuestionNavigatorSummary
                class="px-1"
                :answered="answeredCount"
                :unanswered="unansweredQuestions.length"
                :flagged="flaggedQuestions.length"
            />

            <USeparator />

            <div class="max-h-[calc(100vh-13rem)] space-y-1 overflow-y-auto">
                <div class="flex items-center gap-3 px-2 pb-1 text-xs text-muted">
                    <span class="inline-flex items-center gap-1">
                        <UIcon name="i-mdi-check-circle" class="text-success" />
                        {{ $t('components.qualifications.exam_view.answered') }}
                    </span>
                    <span class="inline-flex items-center gap-1">
                        <UIcon name="i-mdi-circle-outline" />
                        {{ $t('components.qualifications.exam_view.unanswered') }}
                    </span>
                </div>

                <ExamQuestionNavigatorList
                    :questions="questions"
                    :answered-ids="props.answeredIds"
                    :flagged-ids="props.flaggedIds"
                    @keyboard-navigate="navigateFromKeyboard"
                    @navigate="navigate"
                />
            </div>
        </div>
    </div>

    <div class="absolute right-4 bottom-4 z-30 lg:hidden">
        <UTooltip :text="$t('components.qualifications.exam_view.question_navigator')">
            <UButton
                color="neutral"
                size="lg"
                icon="i-mdi-format-list-numbered"
                :aria-label="$t('components.qualifications.exam_view.question_navigator')"
                @click="open = true"
            />
        </UTooltip>
    </div>

    <USlideover v-model:open="open" :title="$t('components.qualifications.exam_view.question_navigator')">
        <template #body>
            <div class="flex flex-col gap-2">
                <p class="text-sm text-muted">
                    {{ $t('components.qualifications.exam_view.question_navigator_description') }}
                </p>
                <ExamQuestionNavigatorSummary
                    :answered="answeredCount"
                    :unanswered="unansweredQuestions.length"
                    :flagged="flaggedQuestions.length"
                />

                <UButton
                    v-if="unansweredQuestions.length > 0"
                    class="justify-between"
                    color="warning"
                    variant="soft"
                    trailing-icon="i-mdi-arrow-right"
                    :label="$t('components.qualifications.exam_view.review_unanswered', { count: unansweredQuestions.length })"
                    @click="navigateToFirstUnanswered"
                />
                <UButton
                    v-if="flaggedQuestions.length > 0"
                    class="justify-between"
                    color="warning"
                    variant="soft"
                    trailing-icon="i-mdi-flag"
                    :label="$t('components.qualifications.exam_view.review_flagged', { count: flaggedQuestions.length })"
                    @click="navigateToFirstFlagged"
                />

                <div class="mb-1 flex items-center gap-3 px-1 text-xs text-muted">
                    <span class="inline-flex items-center gap-1">
                        <UIcon name="i-mdi-check-circle" class="text-success" />
                        {{ $t('components.qualifications.exam_view.answered') }}
                    </span>
                    <span class="inline-flex items-center gap-1">
                        <UIcon name="i-mdi-circle-outline" />
                        {{ $t('components.qualifications.exam_view.unanswered') }}
                    </span>
                </div>

                <ExamQuestionNavigatorList
                    :questions="questions"
                    :answered-ids="props.answeredIds"
                    :flagged-ids="props.flaggedIds"
                    @keyboard-navigate="navigateFromKeyboard"
                    @navigate="navigate"
                />
            </div>
        </template>
    </USlideover>
</template>
