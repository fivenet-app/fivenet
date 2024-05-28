<script lang="ts" setup>
import { z } from 'zod';
import { VueDraggable } from 'vue-draggable-plus';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';

const props = defineProps<{
    modelValue: ExamQuestion;
}>();

const emits = defineEmits<{
    (e: 'update:modelValue', value: ExamQuestion): void;
    (e: 'delete'): void;
}>();

const question = useVModel(props, 'modelValue', emits);

const schema = z.object({
    id: z.string(),
    title: z.string().min(3).max(512),
    description: z.string().max(1024).optional(),
    data: z.object({
        data: z.union([
            z.object({
                oneofKind: z.literal('separator'),
                separator: z.object({}),
            }),
            z.object({
                oneofKind: z.literal('yesno'),
                yesno: z.object({}),
            }),
            z.object({
                oneofKind: z.literal('freeText'),
                freeText: z.object({
                    minLength: z.number().nonnegative(),
                    maxLength: z.number().nonnegative(),
                }),
            }),
            z.object({
                oneofKind: z.literal('multipleChoice'),
                multipleChoice: z.object({
                    multi: z.boolean(),
                    limit: z.number().positive().optional(),
                    choices: z.string().max(255).array().max(10),
                }),
            }),
        ]),
    }),
});

const questionTypes = ['separator', 'yesno', 'freeText', 'singleChoice', 'multipleChoice'];

function changeQuestionType(qt: string): void {
    switch (qt) {
        case 'yesno':
            question.value.data = {
                data: {
                    oneofKind: 'yesno',
                    yesno: {},
                },
            };
            break;

        case 'freeText':
            question.value.data = {
                data: {
                    oneofKind: 'freeText',
                    freeText: {
                        minLength: 0,
                        maxLength: 0,
                    },
                },
            };
            break;

        case 'singleChoice':
            question.value.data = {
                data: {
                    oneofKind: 'singleChoice',
                    singleChoice: {
                        choices: [''],
                    },
                },
            };
            break;

        case 'multipleChoice':
            question.value.data = {
                data: {
                    oneofKind: 'multipleChoice',
                    multipleChoice: {
                        choices: [''],
                        limit: 3,
                    },
                },
            };
            break;

        case 'separator':
        default:
            question.value.data = {
                data: {
                    oneofKind: 'separator',
                    separator: {},
                },
            };
            break;
    }
}
</script>

<template>
    <UForm :schema="schema" :state="question" class="flex items-center gap-2">
        <UIcon name="i-mdi-drag-horizontal" class="size-7" />

        <UFormGroup name="data.data.oneofKind">
            <USelectMenu
                v-model="question.data!.data.oneofKind"
                :options="questionTypes"
                class="w-40 max-w-40"
                @update:model-value="changeQuestionType($event)"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            >
                <template #label>
                    <span class="truncate">
                        {{ $t(`components.qualifications.exam_editor.question_types.${question.data!.data.oneofKind}`) }}
                    </span>
                </template>
                <template #option="{ option }">
                    <span class="truncate">
                        {{ $t(`components.qualifications.exam_editor.question_types.${option}`) }}
                    </span>
                </template>
                <template #option-empty="{ query: search }">
                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                </template>
                <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
            </USelectMenu>
        </UFormGroup>

        <div class="flex flex-1 flex-col gap-2 p-4">
            <div class="flex flex-1 flex-col gap-2">
                <UFormGroup name="title" :label="$t('common.title')" required>
                    <UInput
                        v-model="question.title"
                        type="text"
                        :placeholder="$t('common.title')"
                        size="xl"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>

                <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                    <UTextarea
                        v-model="question.description"
                        type="text"
                        :rows="3"
                        :placeholder="$t('common.description')"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>
            </div>
            <div class="flex-1">
                <template v-if="question.data!.data.oneofKind === 'separator'">
                    <UDivider class="mb-2 mt-2 text-xl">
                        <h4 class="text-xl">{{ question.title }}</h4>
                    </UDivider>

                    <p class="mb-2">{{ question.description }}</p>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'yesno'">
                    <div class="flex flex-col gap-2">
                        <UButtonGroup>
                            <UButton color="green" :label="$t('common.yes')" />
                            <UButton color="red" :label="$t('common.no')" />
                        </UButtonGroup>
                    </div>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'freeText'">
                    <div class="flex flex-col gap-2">
                        <div class="flex gap-2">
                            <UFormGroup name="data.data.freeText.minLength" :label="$t('common.min')" class="flex-1">
                                <UInput
                                    v-model="question.data!.data.freeText.minLength"
                                    type="number"
                                    :min="0"
                                    :max="Number.MAX_SAFE_INTEGER"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup name="data.data.freeText.minLength" :label="$t('common.max')" class="flex-1">
                                <UInput
                                    v-model="question.data!.data.freeText.minLength"
                                    type="number"
                                    :min="0"
                                    :max="Number.MAX_SAFE_INTEGER"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>
                        </div>

                        <UTextarea disabled :rows="5" />
                    </div>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'singleChoice'">
                    <div class="flex flex-col gap-2">
                        <UFormGroup
                            name="data.data.singleChoices.choices"
                            :label="$t('common.option', 2)"
                            required
                            class="flex-1"
                        >
                            <VueDraggable v-model="question.data!.data.singleChoice.choices" class="flex flex-col gap-2">
                                <div
                                    v-for="(_, idx) in question.data!.data.singleChoice?.choices"
                                    class="inline-flex items-center gap-2"
                                >
                                    <UIcon name="i-mdi-drag-horizontal" class="size-6" />
                                    <URadio disabled />
                                    <UFormGroup :name="`data.data.singleChoices.choices.${idx}`">
                                        <UInput
                                            v-model="question.data!.data.singleChoice.choices[idx]"
                                            type="text"
                                            class="w-full"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </UFormGroup>

                                    <UButton
                                        icon="i-mdi-close"
                                        :ui="{ rounded: 'rounded-full' }"
                                        class="flex-initial"
                                        @click="question.data!.data.singleChoice.choices.splice(idx, 1)"
                                    />
                                </div>
                            </VueDraggable>

                            <UButton
                                icon="i-mdi-plus"
                                :ui="{ rounded: 'rounded-full' }"
                                :class="question.data!.data.singleChoice.choices.length ? 'mt-2' : ''"
                                @click="question.data!.data.singleChoice.choices.push('')"
                            />
                        </UFormGroup>
                    </div>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'multipleChoice'">
                    <div class="flex flex-col gap-2">
                        <UFormGroup name="data.data.multipleChoice.limit" :label="$t('common.max')">
                            <UInput
                                v-model="question.data!.data.multipleChoice.limit"
                                type="number"
                                min="1"
                                :max="question.data!.data.multipleChoice.choices.length"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup :label="$t('common.option', 2)" required class="flex-1">
                            <VueDraggable v-model="question.data!.data.multipleChoice.choices" class="flex flex-col gap-2">
                                <div
                                    v-for="(_, idx) in question.data!.data.multipleChoice?.choices"
                                    class="inline-flex items-center gap-2"
                                >
                                    <UIcon name="i-mdi-drag-horizontal" class="size-6" />
                                    <UCheckbox disabled />
                                    <UInput
                                        v-model="question.data!.data.multipleChoice.choices[idx]"
                                        type="text"
                                        block
                                        class="w-full"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <UButton
                                        icon="i-mdi-close"
                                        :ui="{ rounded: 'rounded-full' }"
                                        class="flex-initial"
                                        @click="question.data!.data.multipleChoice.choices.splice(idx, 1)"
                                    />
                                </div>
                            </VueDraggable>

                            <UButton
                                icon="i-mdi-plus"
                                :ui="{ rounded: 'rounded-full' }"
                                :class="question.data!.data.multipleChoice.choices.length ? 'mt-2' : ''"
                                @click="question.data!.data.multipleChoice.choices.push('')"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </div>
        </div>

        <UButton
            icon="i-mdi-close"
            :ui="{ rounded: 'rounded-full' }"
            class="flex-initial self-start"
            @click="$emit('delete')"
        />
    </UForm>
</template>
