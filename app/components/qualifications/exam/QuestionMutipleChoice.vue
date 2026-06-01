<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';

defineProps<{
    disabled?: boolean;
    index: number;
}>();

const question = defineModel<ExamQuestion>({ required: true });

// Only access multipleChoice if oneofKind is 'multipleChoice'
const multipleChoiceChoices = computed<string[]>(() =>
    question.value.data?.data.oneofKind === 'multipleChoice' ? question.value.data.data.multipleChoice.choices : [],
);
const validMultipleChoiceChoices = computed<string[]>(() =>
    multipleChoiceChoices.value.filter((choice) => choice.trim().length > 0),
);
const { moveUp, moveDown } = useListReorder(multipleChoiceChoices);

watch(
    validMultipleChoiceChoices,
    (choices) => {
        if (question.value.answer?.answer.oneofKind === 'multipleChoice') {
            const selectedChoices = question.value.answer.answer.multipleChoice.choices;
            question.value.answer.answer.multipleChoice.choices = selectedChoices.filter((choice) => choices.includes(choice));
        }
    },
    { immediate: true },
);
</script>

<template>
    <div
        v-if="question.data!.data.oneofKind === 'multipleChoice' && question.answer!.answer.oneofKind === 'multipleChoice'"
        class="flex flex-col gap-2"
    >
        <UFormField :name="`exam.questions.${index}.data.data.multipleChoice.limit`" :label="$t('common.max')">
            <UInputNumber
                v-model="question.data!.data.multipleChoice.limit"
                :min="1"
                :max="question.data!.data.multipleChoice.choices.length"
                :disabled="disabled"
            />
        </UFormField>

        <UFormField
            class="flex-1"
            :name="`exam.questions.${index}.data.data.multipleChoice.choices`"
            :label="$t('common.option', 2)"
            required
        >
            <VueDraggable
                v-model="question.data!.data.multipleChoice.choices"
                class="flex flex-col gap-2"
                :disabled="disabled"
                handle=".handle-choice"
            >
                <div
                    v-for="(_, idx) in question.data!.data.multipleChoice?.choices"
                    :key="idx"
                    class="inline-flex items-center gap-2"
                >
                    <div class="inline-flex items-center gap-1">
                        <DraggableHandle handle-class="handle-choice" />

                        <ReorderButtons :idx="idx" :move-up="moveUp" :move-down="moveDown" />
                    </div>

                    <UFormField>
                        <UCheckboxGroup
                            v-model="question.answer!.answer.multipleChoice.choices"
                            :value="question.data!.data.multipleChoice.choices[idx]"
                            :disabled="disabled"
                        />
                    </UFormField>

                    <UFormField :name="`exam.questions.${index}.data.data.multipleChoice.choices.${idx}`">
                        <UInput
                            v-model="question.data!.data.multipleChoice.choices[idx]"
                            class="w-full"
                            type="text"
                            block
                            :disabled="disabled"
                        />
                    </UFormField>

                    <UTooltip :text="$t('components.qualifications.remove_option')">
                        <UButton
                            class="flex-initial"
                            icon="i-mdi-close"
                            :disabled="disabled"
                            @click="question.data!.data.multipleChoice.choices.splice(idx, 1)"
                        />
                    </UTooltip>
                </div>
            </VueDraggable>

            <UTooltip :text="$t('components.qualifications.add_option')">
                <UButton
                    :class="question.data!.data.multipleChoice.choices.length ? 'mt-2' : ''"
                    icon="i-mdi-plus"
                    :disabled="disabled"
                    @click="question.data!.data.multipleChoice.choices.push('')"
                />
            </UTooltip>

            <UFormField
                class="mt-2"
                :name="`exam.questions.${index}.answer.multipleChoice.choices`"
                :label="$t('common.answer')"
            >
                <USelect
                    v-model="question.answer!.answer.multipleChoice.choices"
                    class="w-full"
                    multiple
                    :items="validMultipleChoiceChoices"
                    :disabled="disabled || validMultipleChoiceChoices.length === 0"
                />
            </UFormField>
        </UFormField>
    </div>
</template>
