<script lang="ts" setup>
import type { ExamResponse } from '~~/gen/ts/resources/qualifications/exam';

const props = withDefaults(
    defineProps<{
        modelValue: ExamResponse | undefined;
        disabled?: boolean;
    }>(),
    {
        disabled: false,
    },
);

const emits = defineEmits<{
    (e: 'update:modelValue', value: ExamResponse | undefined): void;
}>();

const response = useVModel(props, 'modelValue', emits);
</script>

<template>
    <div v-if="modelValue?.question" class="flex flex-1 justify-between gap-2 py-4">
        <div v-if="modelValue?.question.data!.data.oneofKind === 'separator'">
            <UDivider class="mb-2 mt-2 text-xl">
                <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                    {{ modelValue?.question.title }}
                </h4>
            </UDivider>

            <p>{{ modelValue?.question.description }}</p>
        </div>

        <div
            v-else-if="
                modelValue?.question.data!.data.oneofKind === 'yesno' && response?.response?.response.oneofKind === 'yesno'
            "
            class="flex flex-1 flex-col gap-2"
        >
            <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">{{ modelValue?.question.title }}</h4>

            <div class="flex flex-1 justify-between gap-2">
                <p>{{ modelValue?.question.description }}</p>
                <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
            </div>

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
            v-else-if="
                modelValue?.question.data!.data.oneofKind === 'freeText' &&
                response?.response?.response.oneofKind === 'freeText'
            "
            class="flex flex-1 flex-col gap-2"
        >
            <div class="flex flex-1 flex-col gap-2">
                <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                    {{ modelValue?.question.title }}
                </h4>

                <div class="inline-flex gap-2">
                    <p>{{ modelValue?.question.description }}</p>
                    <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                </div>

                <div>
                    <UBadge v-if="modelValue?.question.data!.data.freeText.minLength > 0">
                        {{ $t('common.min') }}: {{ modelValue?.question.data!.data.freeText.minLength }}
                        {{ $t('common.chars', modelValue?.question.data!.data.freeText.minLength) }}
                    </UBadge>
                    <UBadge v-if="modelValue?.question.data!.data.freeText.maxLength > 0">
                        {{ $t('common.max') }}: {{ modelValue?.question.data!.data.freeText.maxLength }}
                        {{ $t('common.chars', modelValue?.question.data!.data.freeText.maxLength) }}
                    </UBadge>
                </div>
            </div>

            <UTextarea v-model="response.response.response.freeText.text" :rows="5" :disabled="disabled" />
        </div>

        <div
            v-else-if="
                modelValue?.question.data!.data.oneofKind === 'singleChoice' &&
                response?.response?.response.oneofKind === 'singleChoice'
            "
            class="flex flex-1 flex-col gap-2"
        >
            <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">{{ modelValue?.question.title }}</h4>

            <div class="inline-flex gap-2">
                <p>{{ modelValue?.question.description }}</p>
                <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
            </div>

            <UFormGroup name="data.data.singleChoices.choices" :label="$t('common.option', 2)" required class="flex-1">
                <URadioGroup
                    v-model="response.response.response.singleChoice.choice"
                    :name="modelValue?.question.data!.data.singleChoice.choices.join(':')"
                    :options="modelValue?.question.data!.data.singleChoice?.choices"
                    :disabled="disabled"
                />
            </UFormGroup>
        </div>

        <div
            v-else-if="
                modelValue?.question.data?.data.oneofKind === 'multipleChoice' &&
                response?.response?.response.oneofKind === 'multipleChoice'
            "
            class="flex flex-1 flex-col gap-2"
        >
            <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">{{ modelValue?.question.title }}</h4>

            <div class="inline-flex gap-2">
                <p>{{ modelValue?.question.description }}</p>
                <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
            </div>

            <div>
                <UBadge
                    v-if="
                        modelValue?.question.data!.data.multipleChoice.limit &&
                        modelValue?.question.data!.data.multipleChoice.limit > 0
                    "
                >
                    {{ $t('common.max') }}: {{ modelValue?.question.data!.data.multipleChoice.limit }}
                    {{ $t('common.option', modelValue?.question.data!.data.multipleChoice.limit) }}
                </UBadge>
            </div>

            <UFormGroup :label="$t('common.option', 2)" required class="flex-1"> </UFormGroup>
            <div class="flex flex-1 flex-col gap-2">
                <UCheckbox
                    v-for="choice in modelValue?.question.data.data.multipleChoice.choices"
                    :key="choice"
                    v-model="response.response.response.multipleChoice.choices"
                    name="data.data.multipleChoice.choices"
                    :label="choice"
                    :disabled="disabled"
                    :value="choice"
                />
            </div>
        </div>

        <slot name="question-after" :question="modelValue?.question"></slot>
    </div>
</template>
