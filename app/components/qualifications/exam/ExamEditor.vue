<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import type { File } from '~~/gen/ts/resources/file/file';
import type { ExamQuestion, ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import { AutoGradeMode, type QualificationExamSettings } from '~~/gen/ts/resources/qualifications/qualifications';
import ExamEditorQuestion from './ExamEditorQuestion.vue';

const props = defineProps<{
    settings: QualificationExamSettings;
    questions: ExamQuestions;
    qualificationId: number;
}>();

const emit = defineEmits<{
    (e: 'update:settings', value: ExamQuestions): void;
    (e: 'update:questions', value: ExamQuestions): void;
    (e: 'fileUploaded', file: File): void;
}>();

const { settings, questions } = useVModels(props, emit);

const schema = z.object({
    settings: z.object({
        time: zodDurationSchema,
        autoGrade: z.coerce.boolean().default(false),
        autoGradeMode: z.nativeEnum(AutoGradeMode).default(AutoGradeMode.STRICT),
        minimumPoints: z.coerce.number().min(0).default(0),
    }),
    questions: z.custom<ExamQuestion>().array().max(100).default([]),
});

if (!settings.value.time) {
    settings.value.time = {
        seconds: 600,
        nanos: 0,
    };
}

const modes = ref<{ mode: AutoGradeMode; selected?: boolean }[]>([
    { mode: AutoGradeMode.STRICT, selected: true },
    { mode: AutoGradeMode.PARTIAL_CREDIT },
]);
</script>

<template>
    <div class="mt-2 flex flex-1 flex-col gap-2 px-2">
        <UForm :schema="schema" :state="settings">
            <h2 class="text-gray-900 dark:text-white">
                {{ $t('common.settings') }}
            </h2>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="settings.time"
                :label="$t('components.qualifications.exam_editor.exam_duration')"
                :ui="{ container: '' }"
            >
                <UInput v-model="settings.time!.seconds" type="number" :min="1" :step="1" :placeholder="$t('common.duration')">
                    <template #trailing>
                        <span class="text-xs text-gray-500 dark:text-gray-400">s</span>
                    </template>
                </UInput>
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="settings.autoGrade"
                :label="$t('components.qualifications.exam_editor.auto_grade.title')"
                :description="$t('components.qualifications.exam_editor.auto_grade.description')"
                :ui="{ container: '' }"
            >
                <UToggle v-model="settings.autoGrade" :placeholder="$t('components.qualifications.exam_editor.auto_grade')" />
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="mode"
                :label="$t('components.qualifications.exam_editor.auto_grade_mode.title')"
                :description="$t('components.qualifications.exam_editor.auto_grade_mode.description')"
                :ui="{ container: '' }"
            >
                <ClientOnly>
                    <USelectMenu
                        v-model="settings.autoGradeMode"
                        :options="modes"
                        value-attribute="mode"
                        :searchable-placeholder="$t('common.search_field')"
                    >
                        <template #label>
                            <span class="truncate">{{
                                $t(`enums.qualifications.AutoGradeMode.${AutoGradeMode[settings.autoGradeMode ?? 0]}`)
                            }}</span>
                        </template>

                        <template #option="{ option }">
                            <span class="truncate">{{
                                $t(`enums.qualifications.AutoGradeMode.${AutoGradeMode[option.mode ?? 0]}`)
                            }}</span>
                        </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="settings.miniumPoints"
                :label="$t('components.qualifications.exam_editor.minimum_points')"
                :ui="{ container: '' }"
            >
                <UInput
                    v-model="settings.minimumPoints"
                    type="number"
                    :min="0"
                    :max="999999"
                    :step="1"
                    :placeholder="$t('components.qualifications.exam_editor.minimum_points')"
                />
            </UFormGroup>

            <UDivider class="mt-2" />

            <h3>{{ $t('common.question', 2) }}</h3>

            <UContainer>
                <VueDraggable
                    v-model="questions.questions"
                    class="flex flex-col gap-4 divide-y divide-gray-100 dark:divide-gray-800"
                >
                    <ExamEditorQuestion
                        v-for="(question, idx) in questions?.questions"
                        :key="idx"
                        v-model="questions.questions[idx]"
                        :qualification-id="props.qualificationId"
                        :question="question"
                        @delete="questions.questions.splice(idx, 1)"
                        @file-uploaded="(file) => $emit('fileUploaded', file)"
                    />
                </VueDraggable>

                <UButton
                    icon="i-mdi-plus"
                    :ui="{ rounded: 'rounded-full' }"
                    @click="
                        questions.questions.push({
                            id: 0,
                            qualificationId: props.qualificationId,
                            title: '',
                            description: '',
                            data: {
                                data: {
                                    oneofKind: 'separator',
                                    separator: {},
                                },
                            },
                            answer: {
                                answerKey: '',
                                answer: {
                                    oneofKind: undefined,
                                },
                            },
                            points: 0,
                        })
                    "
                />
            </UContainer>
        </UForm>
    </div>
</template>
