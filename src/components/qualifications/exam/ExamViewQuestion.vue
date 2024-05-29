<script lang="ts" setup>
import type { ExamQuestion, ExamResponse } from '~~/gen/ts/resources/qualifications/exam';

const props = withDefaults(
    defineProps<{
        question: ExamQuestion;
        modelValue: ExamResponse;
        disabled?: boolean;
    }>(),
    {
        disabled: false,
    },
);

const emits = defineEmits<{
    (e: 'update:model-value', value: ExamResponse): void;
}>();

const response = useVModel(props, 'modelValue', emits);

function updateMultipleChoices(value: string, remove?: boolean): void {
    if (response.value.response?.response.oneofKind !== 'multipleChoice') {
        return;
    }

    if (!remove) {
        if (!response.value.response?.response.multipleChoice.choices.includes(value)) {
            return;
        }

        response.value.response?.response.multipleChoice.choices.push(value);
    } else {
        response.value.response.response.multipleChoice.choices =
            response.value.response?.response.multipleChoice.choices.filter((c) => c !== value);
    }
}
</script>

<template>
    <div class="py-4">
        <div v-if="question.data!.data.oneofKind === 'separator'">
            <UDivider class="mb-2 mt-2 text-xl">
                <h4 class="text-xl">{{ question.title }}</h4>
            </UDivider>

            <p>{{ question.description }}</p>
        </div>

        <div
            v-else-if="question.data!.data.oneofKind === 'yesno' && response?.response?.response.oneofKind === 'yesno'"
            class="flex flex-col gap-2"
        >
            <h4 class="text-xl">{{ question.title }}</h4>
            <p>{{ question.description }}</p>

            <UButtonGroup>
                <UButton
                    :variant="
                        response.response?.response.oneofKind === 'yesno' && response.response?.response.yesno.value
                            ? 'solid'
                            : 'outline'
                    "
                    color="green"
                    :label="$t('common.yes')"
                    block
                    class="w-20"
                    :disabled="disabled"
                    @click="response.response.response.yesno.value = true"
                />
                <UButton
                    :variant="
                        response.response?.response.oneofKind === 'yesno' && !response.response?.response.yesno.value
                            ? 'solid'
                            : 'outline'
                    "
                    color="red"
                    :label="$t('common.no')"
                    block
                    class="w-20"
                    :disabled="disabled"
                    @click="response.response.response.yesno.value = false"
                />
            </UButtonGroup>
        </div>

        <div
            v-else-if="question.data!.data.oneofKind === 'freeText' && response.response?.response.oneofKind === 'freeText'"
            class="flex flex-col gap-2"
        >
            <div class="flex flex-col gap-2">
                <h4 class="text-xl">{{ question.title }}</h4>
                <p>{{ question.description }}</p>

                <div>
                    <UBadge v-if="question.data!.data.freeText.minLength > 0">
                        {{ $t('common.min') }}: {{ question.data!.data.freeText.minLength }}
                        {{ $t('common.chars', question.data!.data.freeText.minLength) }}
                    </UBadge>
                    <UBadge v-if="question.data!.data.freeText.maxLength > 0">
                        {{ $t('common.max') }}: {{ question.data!.data.freeText.maxLength }}
                        {{ $t('common.chars', question.data!.data.freeText.maxLength) }}
                    </UBadge>
                </div>
            </div>

            <UTextarea v-model="response.response.response.freeText.text" :rows="5" :disabled="disabled" />
        </div>

        <div
            v-else-if="
                question.data!.data.oneofKind === 'singleChoice' && response.response?.response.oneofKind === 'singleChoice'
            "
            class="flex flex-col gap-2"
        >
            <h4 class="text-xl">{{ question.title }}</h4>
            <p>{{ question.description }}</p>

            <UFormGroup name="data.data.singleChoices.choices" :label="$t('common.option', 2)" required class="flex-1">
                <URadioGroup
                    v-model="response.response.response.singleChoice.choice"
                    :options="question.data!.data.singleChoice?.choices"
                    :disabled="disabled"
                />
            </UFormGroup>
        </div>

        <div
            v-else-if="
                question.data!.data.oneofKind === 'multipleChoice' && response.response?.response.oneofKind === 'multipleChoice'
            "
            class="flex flex-col gap-2"
        >
            <h4 class="text-xl">{{ question.title }}</h4>
            <p>{{ question.description }}</p>

            <div>
                <UBadge v-if="question.data!.data.multipleChoice.limit && question.data!.data.multipleChoice.limit > 0">
                    {{ $t('common.max') }}: {{ question.data!.data.multipleChoice.limit }}
                    {{ $t('common.option', question.data!.data.multipleChoice.limit) }}
                </UBadge>
            </div>

            <UFormGroup :label="$t('common.option', 2)" required class="flex-1">
                <div class="flex flex-col gap-2">
                    <UCheckbox
                        v-for="choice in question.data!.data.multipleChoice?.choices"
                        :label="choice"
                        :disabled="disabled"
                        @update:model-value="updateMultipleChoices(choice, $event)"
                    />
                </div>
            </UFormGroup>
        </div>
    </div>
</template>
