<script lang="ts" setup>
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import type { ExamResponse } from '~~/gen/ts/resources/qualifications/exam/exam';
import ExamViewQuestionHeader from './ExamViewQuestionHeader.vue';

const props = withDefaults(
    defineProps<{
        disabled?: boolean;
        yesNoAnswered?: boolean;
        flagged?: boolean;
    }>(),
    {
        disabled: false,
    },
);

const emit = defineEmits<{
    (e: 'yesNoAnswered'): void;
    (e: 'toggleFlag'): void;
}>();

const modelValue = defineModel<ExamResponse | undefined>({
    required: true,
});

const yesNoValue = computed(() => {
    if (modelValue.value?.response?.response.oneofKind !== 'yesno') {
        return undefined;
    }

    return modelValue.value.response.response.yesno.value;
});

// In the live exam, yesNoAnswered distinguishes an unanswered question from
// an explicitly selected "No". Read-only result views do not have that local
// tracker, so the persisted response itself is the source of truth there.
const isYesNoAnswered = computed(() =>
    props.yesNoAnswered === undefined ? modelValue.value?.response?.response.oneofKind === 'yesno' : props.yesNoAnswered,
);

function setYesNoValue(value: boolean): void {
    if (modelValue.value?.response?.response.oneofKind !== 'yesno') {
        return;
    }

    modelValue.value.response.response.yesno.value = value;
    emit('yesNoAnswered');
}

function enforceMultipleChoiceLimit(choices: string[]): void {
    if (
        modelValue.value?.question?.data?.data.oneofKind !== 'multipleChoice' ||
        modelValue.value?.response?.response.oneofKind !== 'multipleChoice'
    ) {
        return;
    }

    const limit = modelValue.value.question.data.data.multipleChoice.limit ?? 0;
    if (limit > 0 && choices.length > limit) {
        modelValue.value.response.response.multipleChoice.choices = choices.slice(0, limit);
    }
}
</script>

<template>
    <UCard v-if="modelValue?.question" :ui="{ header: 'p-4 sm:p-4', body: 'p-4 sm:p-4', footer: 'p-4 sm:p-4' }">
        <template #header>
            <template v-if="modelValue.question.data?.data.oneofKind === 'separator'">
                <USeparator class="text-xl">
                    <template v-if="modelValue.question.title !== ''" #default>
                        <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue.question.id}`">
                            {{ modelValue.question.title }}
                        </h4>
                    </template>
                </USeparator>
            </template>

            <ExamViewQuestionHeader
                v-else
                :question="modelValue.question"
                :flagged="flagged"
                :disabled="disabled"
                @toggle-flag="emit('toggleFlag')"
            />
        </template>

        <div class="flex flex-1 flex-row gap-2">
            <div v-if="modelValue.question.data!.data.oneofKind === 'separator'" class="flex flex-col gap-2">
                <p v-if="modelValue.question.description" class="text-muted">{{ modelValue.question.description }}</p>
            </div>

            <div v-else-if="modelValue.question.data?.data.oneofKind === 'image'" class="flex flex-col gap-2">
                <GenericImg
                    class="min-h-12 min-w-12"
                    enable-popup
                    :rounded="false"
                    :src="modelValue.question.data.data.image?.image?.filePath"
                    :alt="modelValue.question.data.data.image?.alt ?? $t('common.image')"
                    src-fallback
                />
            </div>

            <div
                v-else-if="
                    modelValue?.question.data!.data.oneofKind === 'yesno' &&
                    modelValue?.response?.response.oneofKind === 'yesno'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <UFieldGroup>
                    <UButton
                        class="w-20"
                        :variant="isYesNoAnswered && yesNoValue ? 'solid' : 'outline'"
                        color="success"
                        :label="$t('common.yes')"
                        block
                        :disabled="disabled"
                        @click="setYesNoValue(true)"
                    />
                    <UButton
                        class="w-20"
                        :variant="isYesNoAnswered && yesNoValue === false ? 'solid' : 'outline'"
                        color="error"
                        :label="$t('common.no')"
                        block
                        :disabled="disabled"
                        @click="setYesNoValue(false)"
                    />
                </UFieldGroup>
            </div>

            <div
                v-else-if="
                    modelValue?.question.data!.data.oneofKind === 'freeText' &&
                    modelValue?.response?.response.oneofKind === 'freeText'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <div
                    v-if="
                        modelValue?.question.data!.data.freeText.minLength > 0 ||
                        modelValue?.question.data!.data.freeText.maxLength > 0
                    "
                    class="flex flex-1 flex-col gap-2"
                >
                    <div>
                        <UBadge
                            v-if="modelValue?.question.data!.data.freeText.minLength > 0"
                            :label="`${$t('common.min')}: ${modelValue?.question.data!.data.freeText.minLength} ${$t('common.chars', modelValue?.question.data!.data.freeText.minLength)}`"
                        />
                        <UBadge
                            v-if="modelValue?.question.data!.data.freeText.maxLength > 0"
                            :label="`${$t('common.max')}: ${modelValue?.question.data!.data.freeText.maxLength} ${$t('common.chars', modelValue?.question.data!.data.freeText.maxLength)}`"
                        />
                    </div>
                </div>

                <UTextarea
                    v-model="modelValue.response.response.freeText.text"
                    :rows="5"
                    :maxlength="modelValue?.question.data!.data.freeText.maxLength || undefined"
                    :disabled="disabled"
                />
            </div>

            <div
                v-else-if="
                    modelValue?.question.data!.data.oneofKind === 'singleChoice' &&
                    modelValue?.response?.response.oneofKind === 'singleChoice'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <UFormField
                    class="flex-1"
                    name="data.data.singleChoice.choices"
                    :label="$t('common.option', 2)"
                    :ui="{ label: 'font-semibold' }"
                >
                    <URadioGroup
                        v-model="modelValue.response.response.singleChoice.choice"
                        :name="modelValue?.question.data!.data.singleChoice.choices.join(':')"
                        :items="modelValue?.question.data!.data.singleChoice?.choices"
                        :disabled="disabled"
                    />
                </UFormField>
            </div>

            <div
                v-else-if="
                    modelValue?.question.data?.data.oneofKind === 'multipleChoice' &&
                    modelValue?.response?.response.oneofKind === 'multipleChoice'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <div
                    v-if="
                        modelValue?.question.data!.data.multipleChoice.limit &&
                        modelValue?.question.data!.data.multipleChoice.limit > 0
                    "
                >
                    <UBadge
                        :label="`${$t('common.max')}: ${modelValue?.question.data!.data.multipleChoice.limit} ${$t('common.option', modelValue?.question.data!.data.multipleChoice.limit)}`"
                    />
                </div>

                <UFormField class="flex-1" :label="$t('common.option', 2)" :ui="{ label: 'font-semibold' }">
                    <div class="flex flex-1 flex-col gap-2">
                        <UCheckboxGroup
                            v-model="modelValue.response.response.multipleChoice.choices"
                            name="data.data.multipleChoice.choices"
                            :disabled="disabled"
                            :items="modelValue.question.data.data.multipleChoice.choices"
                            @update:model-value="enforceMultipleChoiceLimit"
                        />
                    </div>
                </UFormField>
            </div>

            <slot name="question-after" :question="modelValue?.question" :disabled="disabled" />
        </div>

        <template v-if="$slots['question-below']" #footer>
            <slot name="question-below" :question="modelValue.question" :disabled="disabled" />
        </template>
    </UCard>
</template>
