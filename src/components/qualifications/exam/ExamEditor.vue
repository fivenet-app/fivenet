<script lang="ts" setup>
import { z } from 'zod';
import { ExamQuestionData, type ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import type { QualificationExamSettings } from '~~/gen/ts/resources/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        settings: QualificationExamSettings;
        questions: ExamQuestions;
        qualificationId?: string;
    }>(),
    {
        qualificationId: '0',
    },
);

const emits = defineEmits<{
    (e: 'update:settings', value: ExamQuestions): void;
    (e: 'update:questions', value: ExamQuestions): void;
}>();

const { settings, questions } = useVModels(props, emits);

const schema = z.object({
    settings: z.object({
        time: zodDurationSchema,
    }),
    questions: z
        .object({
            id: z.string(),
            title: z.string().min(3).max(512),
            description: z.string().max(1024).optional(),
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
                        minLength: z.number().positive(),
                        maxLength: z.number().positive(),
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
        })
        .array()
        .max(50),
});

type Schema = z.output<typeof schema>;

const questionTypes = ['separator', 'yesno', 'freeText', 'multipleChoice'];

function changeQuestionType(idx: number, qt: string): void {
    switch (qt) {
        case 'yesno':
            questions.value.questions[idx].data = {
                data: {
                    oneofKind: 'yesno',
                    yesno: {},
                },
            };
            break;

        case 'freeText':
            questions.value.questions[idx].data = {
                data: {
                    oneofKind: 'freeText',
                    freeText: {
                        minLength: 0,
                        maxLength: 0,
                    },
                },
            };
            break;

        case 'multipleChoice':
            questions.value.questions[idx].data = {
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
            questions.value.questions[idx].data = {
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
    <div class="mt-2 flex flex-col gap-2 px-2">
        <UForm :schema="schema" :state="settings">
            <div>
                <h2 class="text- text-gray-900 dark:text-white">
                    {{ $t('common.settings') }}
                </h2>

                <UFormGroup name="settings.time">
                    <UInput
                        type="text"
                        :placeholder="$t('common.duration')"
                        :value="settings.time ? fromDuration(settings.time) : '600s'"
                        @update:modelValue="settings.time = toDuration($event)"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>

                <h3>{{ $t('common.question', 2) }}</h3>

                <UContainer>
                    <div class="flex flex-col gap-2">
                        <div v-for="(question, idx) in questions?.questions" class="flex items-center gap-2">
                            <UFormGroup :label="$t('common.type')" required>
                                <USelectMenu
                                    v-model="questions!.questions[idx].data!.data.oneofKind"
                                    :options="questionTypes"
                                    @update:model-value="changeQuestionType(idx, $event)"
                                >
                                    <template #label>
                                        {{ questions!.questions[idx].data!.data.oneofKind }}
                                    </template>
                                    <template #option="{ option }">
                                        {{ option }}
                                    </template>
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                </USelectMenu>
                            </UFormGroup>

                            <div class="flex flex-1 flex-col gap-2">
                                <div class="flex flex-1 flex-row gap-2">
                                    <UFormGroup :label="$t('common.title')" required>
                                        <UInput v-model="question.title" type="text" :placeholder="$t('common.title')" />
                                    </UFormGroup>

                                    <UFormGroup
                                        v-if="question.data!.data.oneofKind !== 'separator'"
                                        :label="$t('common.description')"
                                    >
                                        <UInput
                                            v-model="question.description"
                                            type="text"
                                            :placeholder="$t('common.description')"
                                        />
                                    </UFormGroup>
                                </div>
                                <div class="flex-1">
                                    <template v-if="question.data!.data.oneofKind === 'separator'">
                                        <UDivider class="mb-4 mt-2" :label="question.title" />
                                    </template>

                                    <template v-else-if="question.data!.data.oneofKind === 'yesno'">
                                        <div class="flex flex-col gap-2">
                                            <h3 v-if="question.title" class="text-2xl">{{ question.title }}</h3>
                                            <p v-if="question.description">{{ question.description }}</p>

                                            <UButtonGroup>
                                                <UButton color="green" :label="$t('common.yes')" />
                                                <UButton color="red" :label="$t('common.no')" />
                                            </UButtonGroup>
                                        </div>
                                    </template>

                                    <template v-else-if="question.data!.data.oneofKind === 'multipleChoice'">
                                        <div class="flex flex-col gap-2">
                                            <h3 v-if="question.title" class="text-2xl">{{ question.title }}</h3>
                                            <p v-if="question.description">{{ question.description }}</p>

                                            <UFormGroup label="Limit">
                                                <UInput type="number" min="0" max="10" />
                                            </UFormGroup>

                                            <div class=""></div>

                                            <div class="flex flex-col gap-2">
                                                <div
                                                    v-for="(_, idx) in question.data!.data.multipleChoice?.choices"
                                                    class="inline-flex items-center gap-2"
                                                >
                                                    <UCheckbox disabled />
                                                    <UInput
                                                        v-model="question.data!.data.multipleChoice.choices[idx]"
                                                        type="text"
                                                        block
                                                        class="w-full"
                                                    />

                                                    <UButton
                                                        icon="i-mdi-close"
                                                        :ui="{ rounded: 'rounded-full' }"
                                                        class="flex-initial"
                                                        @click="
                                                            question.data!.data.multipleChoice.choices =
                                                                question.data!.data.multipleChoice.choices.splice(idx, 1)
                                                        "
                                                    />
                                                </div>
                                            </div>

                                            <UButton
                                                icon="i-mdi-plus"
                                                :ui="{ rounded: 'rounded-full' }"
                                                @click="question.data!.data.multipleChoice.choices.push('')"
                                            />
                                        </div>
                                    </template>

                                    <template v-else-if="question.data!.data.oneofKind === 'freeText'">
                                        <div class="flex flex-col gap-2">
                                            <h3 v-if="question.title" class="text-2xl">{{ question.title }}</h3>
                                            <p v-if="question.description">{{ question.description }}</p>

                                            <UTextarea :rows="5" />
                                        </div>
                                    </template>
                                </div>
                            </div>

                            <UButton
                                icon="i-mdi-close"
                                :ui="{ rounded: 'rounded-full' }"
                                class="flex-initial"
                                @click="questions.questions = questions.questions.splice(idx, 1)"
                            />
                        </div>
                    </div>

                    <UButton
                        icon="i-mdi-plus"
                        :ui="{ rounded: 'rounded-full' }"
                        @click="
                            questions.questions.push({
                                id: '0',
                                qualificationId: props.qualificationId,
                                title: '',
                                description: '',
                                data: {
                                    data: {
                                        oneofKind: 'separator',
                                        separator: {},
                                    },
                                },
                            })
                        "
                    />
                </UContainer>

                <!-- TODO -->
            </div>
        </UForm>
    </div>
</template>
