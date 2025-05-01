<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';

const props = defineProps<{
    modelValue?: ExamQuestion;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ExamQuestion): void;
    (e: 'delete'): void;
}>();

const question = useVModel(props, 'modelValue', emit);

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const schema = z.object({
    id: z.number(),
    title: z.string().min(3).max(512),
    description: z.string().max(1024).optional(),
    data: z.object({
        data: z.union([
            z.object({
                oneofKind: z.literal('separator'),
                separator: z.object({}),
            }),
            z.object({
                oneofKind: z.literal('image'),
                image: z.object({
                    alt: z.string().max(128).optional(),
                    image: z.object({
                        url: z.string().optional(),
                        data: z.any(),
                    }),
                }),
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
    answer: z
        .object({
            answerKey: z.string().max(1024),
        })
        .optional(),
    points: z.coerce.number().min(0).max(99999),
});

watch(question, () => {
    if (question.value === undefined) {
        question.value = {
            id: 0,
            qualificationId: 0,
            title: '',
            answer: {
                answerKey: '',
            },
        };
    } else {
        if (question.value.data?.data.oneofKind === 'image') {
            imageUrl.value = question.value.data?.data.image.image?.url;
        }
    }
});

const imageUrl = ref<string | undefined>();
if (question.value?.data?.data.oneofKind === 'image') {
    imageUrl.value = question.value.data?.data.image.image?.url;
}

async function handleImage(files: FileList): Promise<void> {
    if (question.value!.data!.data.oneofKind !== 'image') {
        return;
    }

    if (!files || files.length === 0 || !files[0]) {
        return;
    }

    question.value!.data!.data.image.image = { data: new Uint8Array(await files[0].arrayBuffer()) };

    imageUrl.value = URL.createObjectURL(files[0]);
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
            break;

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
    <UForm v-if="question" :schema="schema" :state="question" class="flex items-center gap-2">
        <UIcon name="i-mdi-drag-horizontal" class="size-7" />

        <UFormGroup name="data.data.oneofKind">
            <ClientOnly>
                <USelectMenu
                    v-model="question.data!.data.oneofKind"
                    :options="questionTypes"
                    class="w-40 max-w-40"
                    @update:model-value="changeQuestionType($event)"
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
            </ClientOnly>
        </UFormGroup>

        <div class="flex flex-1 flex-col gap-2 p-4">
            <div class="flex flex-1 flex-col gap-2">
                <UFormGroup name="title" :label="$t('common.title')" required>
                    <UInput v-model="question.title" type="text" :placeholder="$t('common.title')" size="xl" />
                </UFormGroup>

                <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                    <UTextarea v-model="question.description" type="text" :rows="3" :placeholder="$t('common.description')" />
                </UFormGroup>
            </div>
            <div class="flex-1">
                <template v-if="question.data!.data.oneofKind === 'separator'">
                    <UDivider class="mb-2 mt-2 text-xl">
                        <h4 class="text-xl">{{ question.title }}</h4>
                    </UDivider>

                    <p class="mb-2">{{ question.description }}</p>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'image'">
                    <div class="flex flex-col gap-2">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <template v-else>
                            <UInput
                                type="file"
                                :accept="appConfig.fileUpload.types.images.join(',')"
                                :placeholder="$t('common.image')"
                                @change="handleImage($event)"
                            />
                        </template>

                        <img v-if="imageUrl" :src="imageUrl" />
                    </div>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'yesno'">
                    <div class="flex flex-col gap-2">
                        <UButtonGroup>
                            <UButton color="green" :label="$t('common.yes')" disabled />
                            <UButton color="error" :label="$t('common.no')" disabled />
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
                                />
                            </UFormGroup>

                            <UFormGroup name="data.data.freeText.maxLength" :label="$t('common.max')" class="flex-1">
                                <UInput
                                    v-model="question.data!.data.freeText.maxLength"
                                    type="number"
                                    :min="0"
                                    :max="Number.MAX_SAFE_INTEGER"
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
                                    :key="idx"
                                    class="inline-flex items-center gap-2"
                                >
                                    <UIcon name="i-mdi-drag-horizontal" class="size-6" />
                                    <URadio disabled />
                                    <UFormGroup :name="`data.data.singleChoices.choices.${idx}`">
                                        <UInput
                                            v-model="question.data!.data.singleChoice.choices[idx]"
                                            type="text"
                                            class="w-full"
                                        />
                                    </UFormGroup>

                                    <UTooltip :text="$t('components.qualifications.remove_option')">
                                        <UButton
                                            icon="i-mdi-close"
                                            :ui="{ rounded: 'rounded-full' }"
                                            class="flex-initial"
                                            @click="question.data!.data.singleChoice.choices.splice(idx, 1)"
                                        />
                                    </UTooltip>
                                </div>
                            </VueDraggable>

                            <UTooltip :text="$t('components.qualifications.add_option')">
                                <UButton
                                    icon="i-mdi-plus"
                                    :ui="{ rounded: 'rounded-full' }"
                                    :class="question.data!.data.singleChoice.choices.length ? 'mt-2' : ''"
                                    @click="question.data!.data.singleChoice.choices.push('')"
                                />
                            </UTooltip>
                        </UFormGroup>
                    </div>
                </template>

                <template v-else-if="question.data!.data.oneofKind === 'multipleChoice'">
                    <div class="flex flex-col gap-2">
                        <UFormGroup name="data.data.multipleChoice.limit" :label="$t('common.max')">
                            <UInput
                                v-model="question.data!.data.multipleChoice.limit"
                                type="number"
                                :min="1"
                                :max="question.data!.data.multipleChoice.choices.length"
                            />
                        </UFormGroup>

                        <UFormGroup :label="$t('common.option', 2)" required class="flex-1">
                            <VueDraggable v-model="question.data!.data.multipleChoice.choices" class="flex flex-col gap-2">
                                <div
                                    v-for="(_, idx) in question.data!.data.multipleChoice?.choices"
                                    :key="idx"
                                    class="inline-flex items-center gap-2"
                                >
                                    <UIcon name="i-mdi-drag-horizontal" class="size-6" />
                                    <UCheckbox disabled />
                                    <UInput
                                        v-model="question.data!.data.multipleChoice.choices[idx]"
                                        type="text"
                                        block
                                        class="w-full"
                                    />

                                    <UTooltip :text="$t('components.qualifications.remove_option')">
                                        <UButton
                                            icon="i-mdi-close"
                                            :ui="{ rounded: 'rounded-full' }"
                                            class="flex-initial"
                                            @click="question.data!.data.multipleChoice.choices.splice(idx, 1)"
                                        />
                                    </UTooltip>
                                </div>
                            </VueDraggable>

                            <UTooltip :text="$t('components.qualifications.add_option')">
                                <UButton
                                    icon="i-mdi-plus"
                                    :ui="{ rounded: 'rounded-full' }"
                                    :class="question.data!.data.multipleChoice.choices.length ? 'mt-2' : ''"
                                    @click="question.data!.data.multipleChoice.choices.push('')"
                                />
                            </UTooltip>
                        </UFormGroup>
                    </div>
                </template>

                <div
                    v-if="question.data!.data.oneofKind !== 'separator' && question.data!.data.oneofKind !== 'image'"
                    class="mt-2 flex flex-row gap-2"
                >
                    <UFormGroup name="answer.answerKey" :label="$t('common.answer_key')" class="flex-1">
                        <UTextarea v-model="question.answer!.answerKey" :placeholder="$t('common.answer_key')" />
                    </UFormGroup>

                    <UFormGroup name="points" :label="$t('common.points', 2)" class="max-w-24">
                        <UInput
                            v-model="question.points"
                            type="number"
                            name="points"
                            :min="0"
                            :placeholder="$t('common.points', 2)"
                        />
                    </UFormGroup>
                </div>
            </div>
        </div>

        <UTooltip :text="$t('components.qualifications.remove_question')">
            <UButton
                icon="i-mdi-close"
                :ui="{ rounded: 'rounded-full' }"
                class="mt-1 flex-initial self-start"
                @click="$emit('delete')"
            />
        </UTooltip>
    </UForm>
</template>
