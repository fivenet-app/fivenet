<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';
import type { ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import type { QualificationExamSettings } from '~~/gen/ts/resources/qualifications/qualifications';
import ExamEditorQuestion from './ExamEditorQuestion.vue';

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
    questions: z.custom<ExamQuestion>().array().max(50),
});
</script>

<template>
    <div class="mt-2 flex flex-col gap-2 px-2">
        <UForm :schema="schema" :state="settings">
            <h2 class="text- text-gray-900 dark:text-white">
                {{ $t('common.settings') }}
            </h2>

            <UFormGroup name="settings.time" :label="$t('common.duration')">
                <UInput
                    type="number"
                    :placeholder="$t('common.duration')"
                    :value="settings.time ? fromDuration(settings.time) : 600"
                    @update:model-value="settings.time = toDuration($event)"
                >
                    <template #trailing>
                        <span class="text-xs text-gray-500 dark:text-gray-400">s</span>
                    </template>
                </UInput>
            </UFormGroup>

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
                        :question="question"
                        @delete="questions.questions.splice(idx, 1)"
                    />
                </VueDraggable>

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
        </UForm>
    </div>
</template>
