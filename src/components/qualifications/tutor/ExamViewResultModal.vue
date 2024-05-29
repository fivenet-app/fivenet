<script lang="ts" setup>
import type { GetUserExamResponse } from '~~/gen/ts/services/qualifications/qualifications';
import QualificationResultTutorForm from './QualificationResultTutorForm.vue';
import ExamViewQuestions from '../exam/ExamViewQuestions.vue';

const props = defineProps<{
    qualificationId: string;
    userId: number;
}>();

defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const { data } = useLazyAsyncData('', () => getUserExam());

async function getUserExam(): Promise<GetUserExamResponse> {
    try {
        const call = getGRPCQualificationsClient().getUserExam({
            qualificationId: props.qualificationId,
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        throw e;
    }
}

// TODO
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <QualificationResultTutorForm
            :qualification-id="qualificationId"
            :user-id="userId"
            @refresh="$emit('refresh')"
            @close="isOpen = false"
        >
            <template #default>
                <ExamViewQuestions
                    v-if="data?.exam && data?.examUser && data?.responses"
                    :qualification-id="qualificationId"
                    :exam="data.exam"
                    :exam-user="data.examUser"
                    :responses="data.responses"
                />

                <UDivider class="mb-4 mt-2" />
            </template>
        </QualificationResultTutorForm>
    </UModal>
</template>
