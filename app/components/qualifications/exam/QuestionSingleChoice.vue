<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';

defineProps<{
    disabled?: boolean;
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
        <UFormGroup class="flex-1" name="data.data.singleChoices.choices" :label="$t('common.option', 2)" required>
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

                        <UButtonGroup>
                            <UButton size="xs" variant="link" :padded="false" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                            <UButton size="xs" variant="link" :padded="false" icon="i-mdi-arrow-down" @click="moveDown(idx)" />
                        </UButtonGroup>
                    </div>

                    <URadio
                        v-model="question.answer!.answer.singleChoice.choice"
                        :value="question.data!.data.singleChoice.choices[idx]"
                        :disabled="disabled"
                    />
                    <UFormGroup :name="`data.data.singleChoices.choices.${idx}`" class="w-full">
                        <UInput
                            v-model="question.data!.data.singleChoice.choices[idx]"
                            class="w-full"
                            type="text"
                            block
                            :disabled="disabled"
                        />
                    </UFormGroup>

                    <UTooltip :text="$t('components.qualifications.remove_option')">
                        <UButton
                            class="flex-initial"
                            icon="i-mdi-close"
                            :ui="{ rounded: 'rounded-full' }"
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
                    :ui="{ rounded: 'rounded-full' }"
                    :disabled="disabled"
                    @click="question.data!.data.singleChoice.choices.push('')"
                />
            </UTooltip>
        </UFormGroup>
    </div>
</template>
