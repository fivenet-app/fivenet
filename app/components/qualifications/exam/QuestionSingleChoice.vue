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
    question.value.data?.data.oneofKind === 'singleChoice' ? question.value.data.data.singleChoice.choices : [],
);
const { moveUp, moveDown } = useListReorder(singleChoiceChoices);
</script>

<template>
    <div
        v-if="question.data!.data.oneofKind === 'singleChoice' && question.answer!.answer.oneofKind === 'singleChoice'"
        class="flex flex-col gap-2"
    >
        <UFormField
            class="flex-1"
            :name="`exam.questions.${index}.data.data.singleChoice.choices`"
            :label="$t('common.option', 2)"
        >
            <VueDraggable
                v-model="question.data!.data.singleChoice.choices"
                class="flex w-full flex-col gap-2"
                :disabled="disabled"
                handle=".handle-choice"
            >
                <div
                    v-for="(_, idx) in question.data!.data.singleChoice?.choices"
                    :key="idx"
                    class="flex flex-1 items-center gap-2"
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

                    <UFormField :name="`exam.questions.${index}.data.data.singleChoice.choices.${idx}`" class="w-full">
                        <UInput
                            v-model="question.data!.data.singleChoice.choices[idx]"
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
                            @click="question.data!.data.singleChoice.choices.splice(idx, 1)"
                        />
                    </UTooltip>
                </div>
            </VueDraggable>

            <UTooltip :text="$t('components.qualifications.add_option')">
                <UButton
                    :class="question.data!.data.singleChoice.choices.length ? 'mt-2' : ''"
                    icon="i-mdi-plus"
                    :disabled="disabled"
                    @click="question.data!.data.singleChoice.choices.push('')"
                />
            </UTooltip>

            <UFormField :name="`exam.questions.${index}.answer.singleChoice.choice`" :label="$t('common.answer')" class="mt-2">
                <USelect
                    v-model="question.answer!.answer.singleChoice.choice"
                    :items="question.data!.data.singleChoice?.choices"
                    class="w-full"
                />
            </UFormField>
        </UFormField>
    </div>
</template>
