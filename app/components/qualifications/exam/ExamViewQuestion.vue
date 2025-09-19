<script lang="ts" setup>
import GenericImg from '~/components/partials/elements/GenericImg.vue';
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

const emit = defineEmits<{
    (e: 'update:modelValue', value: ExamResponse | undefined): void;
}>();

const response = useVModel(props, 'modelValue', emit);
</script>

<template>
    <div v-if="modelValue?.question" class="flex flex-1 flex-col justify-between gap-2 py-4">
        <div class="flex flex-1 flex-row gap-2">
            <div v-if="modelValue?.question.data!.data.oneofKind === 'separator'">
                <USeparator class="mt-2 mb-2 text-xl">
                    <template v-if="modelValue?.question.title !== ''" #default>
                        <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                            {{ modelValue?.question.title }}
                        </h4>
                    </template>
                </USeparator>

                <p>{{ modelValue?.question.description }}</p>
            </div>

            <div v-else-if="modelValue?.question!.data?.data.oneofKind === 'image'">
                <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                    {{ modelValue?.question.title }}
                </h4>

                <div class="flex flex-1 justify-between gap-2">
                    <p>{{ modelValue?.question.description }}</p>
                    <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                </div>

                <GenericImg
                    class="min-h-12 min-w-12"
                    enable-popup
                    :rounded="false"
                    :src="modelValue?.question!.data?.data.image?.image?.filePath"
                    :alt="modelValue?.question!.data?.data.image?.alt ?? $t('common.image')"
                />
            </div>

            <div
                v-else-if="
                    modelValue?.question.data!.data.oneofKind === 'yesno' && response?.response?.response.oneofKind === 'yesno'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <div class="flex flex-1 flex-row gap-2">
                    <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                        {{ modelValue?.question.title }}
                    </h4>

                    <div class="flex flex-1 justify-between gap-2">
                        <p>{{ modelValue?.question.description }}</p>
                        <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                    </div>
                </div>

                <UButtonGroup>
                    <UButton
                        class="w-20"
                        :variant="
                            response.response?.response.oneofKind === 'yesno' && response.response?.response.yesno.value
                                ? 'solid'
                                : 'outline'
                        "
                        color="green"
                        :label="$t('common.yes')"
                        block
                        :disabled="disabled"
                        @click="response.response.response.yesno.value = true"
                    />
                    <UButton
                        class="w-20"
                        :variant="
                            response.response?.response.oneofKind === 'yesno' && !response.response?.response.yesno.value
                                ? 'solid'
                                : 'outline'
                        "
                        color="error"
                        :label="$t('common.no')"
                        block
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
                    <div class="flex flex-1 flex-row gap-2">
                        <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                            {{ modelValue?.question.title }}
                        </h4>

                        <div class="flex flex-1 justify-between gap-2">
                            <p>{{ modelValue?.question.description }}</p>
                            <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                        </div>
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
                <div class="flex flex-1 flex-row gap-2">
                    <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                        {{ modelValue?.question.title }}
                    </h4>

                    <div class="flex flex-1 justify-between gap-2">
                        <p>{{ modelValue?.question.description }}</p>
                        <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                    </div>
                </div>

                <UFormField class="flex-1" name="data.data.singleChoice.choices" :label="$t('common.option', 2)" required>
                    <URadioGroup
                        v-model="response.response.response.singleChoice.choice"
                        :name="modelValue?.question.data!.data.singleChoice.choices.join(':')"
                        :items="modelValue?.question.data!.data.singleChoice?.choices"
                        :disabled="disabled"
                    />
                </UFormField>
            </div>

            <div
                v-else-if="
                    modelValue?.question.data?.data.oneofKind === 'multipleChoice' &&
                    response?.response?.response.oneofKind === 'multipleChoice'
                "
                class="flex flex-1 flex-col gap-2"
            >
                <div class="flex flex-1 flex-row gap-2">
                    <h4 class="text-xl" :title="`${$t('common.id')}: ${modelValue?.question.id}`">
                        {{ modelValue?.question.title }}
                    </h4>

                    <div class="flex flex-1 justify-between gap-2">
                        <p>{{ modelValue?.question.description }}</p>
                        <p v-if="modelValue?.question.points">{{ $t('common.point', modelValue?.question.points) }}</p>
                    </div>
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

                <UFormField class="flex-1" :label="$t('common.option', 2)" required>
                    <div class="flex flex-1 flex-col gap-2">
                        <UCheckboxGroup
                            v-model="response.response.response.multipleChoice.choices"
                            name="data.data.multipleChoice.choices"
                            :disabled="disabled"
                        />
                    </div>
                </UFormField>
            </div>

            <slot name="question-after" :question="modelValue?.question" :disabled="disabled" />
        </div>

        <slot name="question-below" :question="modelValue?.question" :disabled="disabled" />
    </div>
</template>
