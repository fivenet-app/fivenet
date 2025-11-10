<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';

defineProps<{
    disabled?: boolean;
    index: number;
}>();

const question = defineModel<ExamQuestion>({ required: true });

// Only access singleChoice if oneofKind is 'singleChoice'
const singleChoiceChoices = computed<string[]>(() =>
    question.value.data?.data.oneofKind === 'multipleChoice' ? question.value.data.data.multipleChoice.choices : [],
);
const { moveUp, moveDown } = useListReorder(singleChoiceChoices);
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
            :name="`exam.questions.${index}.data.data.multipleChoice.choices`"
            class="flex-1"
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
                        <UTooltip :text="$t('common.draggable')">
                            <UIcon class="handle-choice size-6 cursor-move" name="i-mdi-drag-horizontal" />
                        </UTooltip>

                        <UFieldGroup>
                            <UButton size="xs" variant="link" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                            <UButton size="xs" variant="link" icon="i-mdi-arrow-down" @click="moveDown(idx)" />
                        </UFieldGroup>
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
                :name="`exam.questions.${index}.answer.multipleChoice.choices`"
                :label="$t('common.answer')"
                class="mt-2"
            >
                <USelect
                    v-model="question.answer!.answer.multipleChoice.choices"
                    multiple
                    :items="question.data!.data.multipleChoice?.choices"
                    class="w-full"
                />
            </UFormField>
        </UFormField>
    </div>
</template>
