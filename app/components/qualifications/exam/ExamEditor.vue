<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import type { File } from '~~/gen/ts/resources/file/file';
import type { ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import { AutoGradeMode, QualificationExamMode } from '~~/gen/ts/resources/qualifications/qualifications';
import ExamEditorQuestion from './ExamEditorQuestion.vue';

const props = defineProps<{
    qualificationId: number;
    disabled?: boolean;
}>();

defineEmits<{
    (e: 'fileUploaded', file: File): void;
}>();

const examModes = ref<{ mode: QualificationExamMode; selected?: boolean }[]>([
    { mode: QualificationExamMode.DISABLED },
    { mode: QualificationExamMode.REQUEST_NEEDED },
    { mode: QualificationExamMode.ENABLED },
]);

const examMode = defineModel<QualificationExamMode>('examMode', { required: true });

const examSettings = defineModel<ExamSettingsSchema>('settings', { required: true });

const exam = defineModel<ExamQuestions>('exam', { required: true });

const modes = ref<{ mode: AutoGradeMode; selected?: boolean }[]>([
    { mode: AutoGradeMode.STRICT, selected: true },
    { mode: AutoGradeMode.PARTIAL_CREDIT },
]);

const { moveUp, moveDown } = useListReorder(toRef(exam.value.questions));
</script>

<script lang="ts">
export const examSettings = z.object({
    time: zodDurationSchema,
    autoGrade: z.coerce.boolean().default(false),
    autoGradeMode: z.nativeEnum(AutoGradeMode).default(AutoGradeMode.STRICT),
    minimumPoints: z.coerce.number().min(0).default(0),
});

export type ExamSettingsSchema = z.output<typeof examSettings>;
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0' }">
        <template #body>
            <UPageCard :title="`${$t('common.exam')} ${$t('common.setting')}`">
                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="examMode"
                    :label="$t('components.qualifications.exam_mode')"
                    :ui="{ container: '' }"
                >
                    <ClientOnly>
                        <USelectMenu v-model="examMode" :items="examModes" value-key="mode" class="w-full">
                            <template #default>
                                {{ $t(`enums.qualifications.QualificationExamMode.${QualificationExamMode[examMode]}`) }}
                            </template>

                            <template #item="{ item }">
                                {{ $t(`enums.qualifications.QualificationExamMode.${QualificationExamMode[item.mode]}`) }}
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.type', 2)]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="examSettings.time"
                    :label="$t('components.qualifications.exam_editor.exam_duration')"
                    :ui="{ container: '' }"
                >
                    <UInputNumber
                        v-model="examSettings.time"
                        :min="1"
                        :step="1"
                        :placeholder="$t('common.duration')"
                        class="w-full"
                    >
                        <template #trailing>
                            <span class="text-xs text-muted">s</span>
                        </template>
                    </UInputNumber>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="examSettings.autoGrade"
                    :label="$t('components.qualifications.exam_editor.auto_grade.title')"
                    :description="$t('components.qualifications.exam_editor.auto_grade.description')"
                    :ui="{ container: '' }"
                >
                    <USwitch
                        v-model="examSettings.autoGrade"
                        :placeholder="$t('components.qualifications.exam_editor.auto_grade')"
                    />
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="mode"
                    :label="$t('components.qualifications.exam_editor.auto_grade_mode.title')"
                    :description="$t('components.qualifications.exam_editor.auto_grade_mode.description')"
                    :ui="{ container: '' }"
                >
                    <ClientOnly>
                        <USelectMenu
                            v-model="examSettings.autoGradeMode"
                            :items="modes"
                            value-key="mode"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            class="w-full"
                        >
                            <template #default>
                                {{ $t(`enums.qualifications.AutoGradeMode.${AutoGradeMode[settings.autoGradeMode ?? 0]}`) }}
                            </template>

                            <template #item="{ item }">
                                {{ $t(`enums.qualifications.AutoGradeMode.${AutoGradeMode[item.mode ?? 0]}`) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField
                    class="grid grid-cols-2 items-center gap-2"
                    name="examSettings.miniumPoints"
                    :label="$t('components.qualifications.exam_editor.minimum_points')"
                    :ui="{ container: '' }"
                >
                    <UInputNumber
                        v-model="examSettings.minimumPoints"
                        :min="0"
                        :max="999999"
                        :step="1"
                        :placeholder="$t('components.qualifications.exam_editor.minimum_points')"
                        class="w-full"
                    />
                </UFormField>
            </UPageCard>

            <div class="my-4 gap-4">
                <UPageCard :ui="{ title: 'inline-flex flex-1', body: 'flex w-full' }">
                    <template #title>
                        <h2>{{ $t('common.question', 2) }}</h2>
                        <div class="flex-1" />
                        <UBadge :label="`${$t('common.count')}: ${exam?.questions.length}`" />
                    </template>
                </UPageCard>

                <UContainer>
                    <VueDraggable
                        v-model="exam.questions"
                        class="flex flex-col gap-4 divide-y divide-default"
                        :disabled="disabled"
                        handle=".handle"
                    >
                        <ExamEditorQuestion
                            v-for="(question, idx) in exam?.questions"
                            :key="idx"
                            v-model="exam.questions[idx]"
                            :qualification-id="props.qualificationId"
                            :question="question"
                            :index="idx"
                            :disabled="disabled"
                            @delete="exam.questions.splice(idx, 1)"
                            @file-uploaded="(file) => $emit('fileUploaded', file)"
                            @move-up="moveUp(idx)"
                            @move-down="moveDown(idx)"
                        />
                    </VueDraggable>

                    <UButton
                        icon="i-mdi-plus"
                        @click="
                            exam.questions.push({
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
                                order: exam.questions.length + 1,
                            })
                        "
                    />
                </UContainer>
            </div>
        </template>
    </UDashboardPanel>
</template>
