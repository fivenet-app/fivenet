<script lang="ts" setup>
import type { WatchStopHandle } from 'vue';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';
import QuestionMutipleChoice from './QuestionMutipleChoice.vue';
import QuestionSingleChoice from './QuestionSingleChoice.vue';

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const props = defineProps<{
    qualificationId: number;
    index: number;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'delete'): void;
    (e: 'fileUploaded', file: File): void;
    (e: 'move-down'): void;
    (e: 'move-up'): void;
}>();

const question = defineModel<ExamQuestion>();

const appConfig = useAppConfig();
const settingsStore = useSettingsStore();

const { nuiEnabled } = storeToRefs(settingsStore);

const schema = z.object({
    id: z.coerce.number(),
    title: z.coerce.string().min(0).max(512),
    description: z.coerce.string().max(1024).optional(),
    data: z.object({
        data: z.union([
            z.object({
                oneofKind: z.literal(undefined),
            }),
            z.object({
                oneofKind: z.literal('separator'),
                separator: z.object({}),
            }),
            z.object({
                oneofKind: z.literal('image'),
                image: z.object({
                    alt: z.coerce.string().max(128).optional(),
                    image: z.custom<File>().optional(),
                }),
            }),
            z.object({
                oneofKind: z.literal('yesno'),
                yesno: z.object({}),
            }),
            z.object({
                oneofKind: z.literal('freeText'),
                freeText: z.object({
                    minLength: z.coerce.number().nonnegative(),
                    maxLength: z.coerce.number().nonnegative(),
                }),
            }),
            z.object({
                oneofKind: z.literal('singleChoice'),
                singleChoice: z.object({
                    choices: z.coerce.string().max(255).array().max(1).default(['']),
                }),
            }),
            z.object({
                oneofKind: z.literal('multipleChoice'),
                multipleChoice: z.object({
                    choices: z.coerce.string().max(255).array().max(10).default([]),
                    limit: z.coerce.number().positive().min(0).max(10).default(10).optional(),
                }),
            }),
        ]),
    }),
    answer: z
        .object({
            answerKey: z.coerce.string().max(1024),
        })
        .optional(),
    points: z.coerce.number().min(0).max(99999),
});

function handleQuestionChange(): void {
    if (question.value === undefined) {
        question.value = {
            id: 0,
            qualificationId: props.qualificationId,
            title: '',
            answer: {
                answerKey: '',
                answer: {
                    oneofKind: undefined,
                },
            },
            order: props.index + 1,
        };
    } else {
        if (question.value.order <= 0) {
            question.value.order = props.index + 1; // Ensure order is always positive and starts from 1
        }
    }

    if (question.value.answer?.answer.oneofKind === undefined) {
        switch (question.value.data?.data.oneofKind) {
            case 'yesno':
                question.value.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: 'yesno',
                        yesno: {
                            value: false,
                        },
                    },
                };
                break;

            case 'freeText':
                question.value.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: 'freeText',
                        freeText: {
                            text: '',
                        },
                    },
                };
                break;

            case 'singleChoice':
                question.value.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: 'singleChoice',
                        singleChoice: {
                            choice: '',
                        },
                    },
                };
                break;

            case 'multipleChoice':
                question.value.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: 'multipleChoice',
                        multipleChoice: {
                            choices: [],
                        },
                    },
                };
                break;

            case 'separator':
            default:
                question.value.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: undefined,
                    },
                };
                break;
        }
    }
}

watch(question, () => handleQuestionChange());
handleQuestionChange();

// Make sure to set the order of the question based on the index prop
watch(
    () => props.index,
    () => {
        if (!question.value) return;

        question.value.order = props.index + 1;
    },
);

const { resizeAndUpload } = useFileUploader(
    (opts) => qualificationsQualificationsClient.uploadFile(opts),
    'qualifications-exam-questions',
    props.qualificationId,
);

async function handleImage(file: globalThis.File | null | undefined): Promise<void> {
    if (!file || question.value!.data!.data.oneofKind !== 'image') return;

    const resp = await resizeAndUpload(file);
    if (question.value?.data?.data.oneofKind === 'image') {
        question.value.data.data.image.image = resp.file;
    }

    question.value!.data!.data.image.image = resp.file!;
    emit('fileUploaded', resp.file!);
}

const questionTypes = ['separator', 'image', 'yesno', 'freeText', 'singleChoice', 'multipleChoice'];

function changeQuestionType(qt: string): void {
    if (question.value === undefined) {
        return;
    }

    switch (qt) {
        case 'image':
            question.value.data = {
                data: {
                    oneofKind: 'image',
                    image: {},
                },
            };
            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: undefined,
                },
            };
            break;

        case 'yesno':
            question.value.data = {
                data: {
                    oneofKind: 'yesno',
                    yesno: {},
                },
            };

            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: 'yesno',
                    yesno: {
                        value: false,
                    },
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

            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: 'freeText',
                    freeText: {
                        text: '',
                    },
                },
            };
            break;

        case 'singleChoice': {
            const choices = [];
            if (
                question.value.data?.data.oneofKind === 'multipleChoice' &&
                question.value.data.data.multipleChoice.choices.length > 0
            ) {
                choices.push(...question.value.data.data.multipleChoice.choices);
            } else {
                choices.push(''); // Start with an empty choice
            }

            question.value.data = {
                data: {
                    oneofKind: 'singleChoice',
                    singleChoice: {
                        choices: choices,
                    },
                },
            };
            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: 'singleChoice',
                    singleChoice: {
                        choice: '', // Placeholder for an undefined choice
                    },
                },
            };
            break;
        }

        case 'multipleChoice': {
            const choices = [];
            if (
                question.value.data?.data.oneofKind === 'singleChoice' &&
                question.value.data.data.singleChoice.choices.length > 0
            ) {
                choices.push(...question.value.data.data.singleChoice.choices);
            } else {
                choices.push(''); // Start with an empty choice
            }

            question.value.data = {
                data: {
                    oneofKind: 'multipleChoice',
                    multipleChoice: {
                        choices: choices,
                        limit: 3,
                    },
                },
            };
            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: 'multipleChoice',
                    multipleChoice: {
                        choices: [],
                    },
                },
            };
            break;
        }

        case 'separator':
        default:
            question.value.data = {
                data: {
                    oneofKind: 'separator',
                    separator: {},
                },
            };
            question.value.answer = {
                answerKey: '',
                answer: {
                    oneofKind: undefined,
                },
            };
            break;
    }
}

let multipleChoiceWatcher: WatchStopHandle | null = null;
let singleChoiceWatcher: WatchStopHandle | null = null;

watch(
    () => question.value?.data?.data.oneofKind,
    (newKind, oldKind) => {
        // Stop previous watchers if they exist
        if (multipleChoiceWatcher) {
            multipleChoiceWatcher();
            multipleChoiceWatcher = null;
        }
        if (singleChoiceWatcher) {
            singleChoiceWatcher();
            singleChoiceWatcher = null;
        }

        // Activate watchers based on the new oneofKind
        if (newKind === 'multipleChoice') {
            multipleChoiceWatcher = watch(
                () =>
                    question.value?.data?.data.oneofKind === 'multipleChoice'
                        ? question.value.data.data.multipleChoice?.choices
                        : undefined,
                (newChoices) => {
                    if (
                        question.value?.answer?.answer?.oneofKind === 'multipleChoice' &&
                        question.value?.answer?.answer?.multipleChoice
                    ) {
                        // Filter answer choices to ensure they are valid values
                        question.value.answer.answer.multipleChoice.choices =
                            question.value.answer.answer.multipleChoice.choices.filter((value) => newChoices?.includes(value));
                    }
                },
                { immediate: true, deep: true },
            );
        } else if (newKind === 'singleChoice') {
            singleChoiceWatcher = watch(
                () =>
                    question.value?.data?.data.oneofKind === 'singleChoice'
                        ? question.value.data.data.singleChoice?.choices
                        : undefined,
                (newChoices) => {
                    if (
                        question.value?.answer?.answer?.oneofKind === 'singleChoice' &&
                        question.value?.answer?.answer?.singleChoice
                    ) {
                        // Reset singleChoice answer if it becomes invalid
                        if (!newChoices?.includes(question.value.answer.answer.singleChoice.choice)) {
                            question.value.answer.answer.singleChoice.choice = ''; // Reset to a placeholder
                        }
                    }
                },
                { immediate: true, deep: true },
            );
        }

        // Reset answer values if the oneofKind changes to a different type
        if (oldKind !== newKind) {
            if (oldKind === 'multipleChoice' && question.value?.answer?.answer.oneofKind === 'multipleChoice') {
                question.value.answer.answer.multipleChoice = { choices: [] };
            } else if (oldKind === 'singleChoice' && question.value?.answer?.answer.oneofKind === 'singleChoice') {
                question.value.answer.answer.singleChoice = { choice: '' };
            }
        }
    },
);
</script>

<template>
    <UForm v-if="question" class="flex items-center gap-2" :schema="schema" :state="question">
        <div class="inline-flex items-center gap-1">
            <UTooltip :text="$t('common.draggable')">
                <UIcon class="handle size-7 cursor-move" name="i-mdi-drag-horizontal" />
            </UTooltip>

            <UButtonGroup>
                <UButton size="xs" variant="link" icon="i-mdi-arrow-up" @click="$emit('move-up')" />
                <UButton size="xs" variant="link" icon="i-mdi-arrow-down" @click="$emit('move-down')" />
            </UButtonGroup>
        </div>

        <UFormField name="data.data.oneofKind">
            <ClientOnly>
                <USelectMenu
                    :model-value="question.data!.data.oneofKind"
                    class="w-40 max-w-40"
                    :items="questionTypes"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :disabled="disabled"
                    @update:model-value="changeQuestionType($event)"
                >
                    <template #default>
                        {{ $t(`components.qualifications.exam_editor.question_types.${question.data!.data.oneofKind}`) }}
                    </template>

                    <template #item-label="{ item }">
                        {{ $t(`components.qualifications.exam_editor.question_types.${item}`) }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                </USelectMenu>
            </ClientOnly>
        </UFormField>

        <div class="flex flex-1 flex-col gap-2 p-4">
            <div class="flex flex-1 flex-col gap-2">
                <UFormField name="title" :label="$t('common.title')" required>
                    <UInput v-model="question.title" type="text" :placeholder="$t('common.title')" size="xl" class="w-full" />
                </UFormField>

                <UFormField class="flex-1" name="description" :label="$t('common.description')">
                    <UTextarea
                        v-model="question.description"
                        type="text"
                        :rows="3"
                        resize
                        :placeholder="$t('common.description')"
                        :disabled="disabled"
                        class="w-full"
                    />
                </UFormField>
            </div>
            <div class="flex-1">
                <template v-if="question.data!.data.oneofKind === 'separator'">
                    <USeparator class="mt-2 mb-2 text-xl">
                        <template v-if="question.title !== ''" #default>
                            <h4 class="text-xl">{{ question.title }}</h4>
                        </template>
                    </USeparator>

                    <p class="mb-2">{{ question.description }}</p>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'image'">
                    <div class="flex flex-col gap-2">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <template v-else>
                            <UFileUpload
                                :accept="appConfig.fileUpload.types.images.join(',')"
                                :disabled="disabled"
                                class="w-full"
                                :placeholder="$t('common.image')"
                                :label="$t('common.file_upload_label')"
                                :description="$t('common.allowed_file_types')"
                                @update:model-value="($event) => handleImage($event)"
                            />
                        </template>

                        <div v-if="question.data?.data.image.image" class="flex flex-1 items-center justify-center">
                            <GenericImg
                                class="min-h-12 min-w-12"
                                :enable-popup="true"
                                :rounded="false"
                                :src="question.data?.data.image.image.filePath"
                                :alt="question.data?.data.image.alt"
                            />
                        </div>
                    </div>
                </template>

                <template
                    v-else-if="question.data!.data.oneofKind === 'yesno' && question.answer!.answer.oneofKind === 'yesno'"
                >
                    <div class="flex flex-col gap-2">
                        <UButtonGroup>
                            <UButton
                                :model-value="question.answer!.answer.yesno.value"
                                color="green"
                                :label="$t('common.yes')"
                                :variant="question.answer!.answer.yesno.value ? 'solid' : 'outline'"
                                :disabled="disabled"
                                @click="question.answer!.answer.yesno.value = true"
                            />
                            <UButton
                                :model-value="question.answer!.answer.yesno.value"
                                color="error"
                                :label="$t('common.no')"
                                :variant="!question.answer!.answer.yesno.value ? 'solid' : 'outline'"
                                :disabled="disabled"
                                @click="question.answer!.answer.yesno.value = false"
                            />
                        </UButtonGroup>
                    </div>
                </template>

                <template
                    v-else-if="question.data!.data.oneofKind === 'freeText' && question.answer!.answer.oneofKind === 'freeText'"
                >
                    <div class="flex flex-col gap-2">
                        <div class="flex gap-2">
                            <UFormField class="flex-1" name="data.data.freeText.minLength" :label="$t('common.min')">
                                <UInputNumber
                                    v-model="question.data!.data.freeText.minLength"
                                    :step="10"
                                    :min="0"
                                    :max="Number.MAX_SAFE_INTEGER"
                                    :disabled="disabled"
                                />
                            </UFormField>

                            <UFormField class="flex-1" name="data.data.freeText.maxLength" :label="$t('common.max')">
                                <UInputNumber
                                    v-model="question.data!.data.freeText.maxLength"
                                    :step="10"
                                    :min="0"
                                    :max="Number.MAX_SAFE_INTEGER"
                                    :disabled="disabled"
                                />
                            </UFormField>
                        </div>

                        <UTextarea
                            v-model="question.answer!.answer.freeText.text"
                            :rows="5"
                            resize
                            :disabled="disabled"
                            class="w-full"
                        />
                    </div>
                </template>

                <template
                    v-else-if="
                        question.data!.data.oneofKind === 'singleChoice' && question.answer!.answer.oneofKind === 'singleChoice'
                    "
                >
                    <QuestionSingleChoice v-model="question" :disabled="disabled" />
                </template>

                <template
                    v-else-if="
                        question.data!.data.oneofKind === 'multipleChoice' &&
                        question.answer!.answer.oneofKind === 'multipleChoice'
                    "
                >
                    <QuestionMutipleChoice v-model="question" :disabled="disabled" />
                </template>

                <div
                    v-if="question.data!.data.oneofKind !== 'separator' && question.data!.data.oneofKind !== 'image'"
                    class="mt-2 flex flex-row gap-2"
                >
                    <UFormField class="flex-1" name="answer.answerKey" :label="$t('common.answer_key')">
                        <UTextarea
                            v-model="question.answer!.answerKey"
                            :placeholder="$t('common.answer_key')"
                            :rows="2"
                            resize
                            :disabled="disabled"
                            class="w-full"
                        />
                    </UFormField>

                    <UFormField class="max-w-24" name="points" :label="$t('common.points', 2)">
                        <UInputNumber
                            v-model="question.points"
                            name="points"
                            :min="0"
                            :placeholder="$t('common.points', 2)"
                            :disabled="disabled"
                        />
                    </UFormField>
                </div>
            </div>
        </div>

        <UTooltip :text="$t('components.qualifications.remove_question')">
            <UButton class="mt-1 flex-initial self-start" icon="i-mdi-close" color="error" @click="$emit('delete')" />
        </UTooltip>
    </UForm>
</template>
