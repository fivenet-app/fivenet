<script lang="ts" setup>
import QualificationResultTutorForm from './QualificationResultTutorForm.vue';

defineProps<{
    qualificationId: number;
    userId?: number;
    resultId?: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();
</script>

<template>
    <UModal>
        <template #content>
            <QualificationResultTutorForm
                :qualification-id="qualificationId"
                :user-id="userId"
                :result-id="resultId"
                @refresh="$emit('refresh')"
                @close="$emit('close', false)"
            />
        </template>
    </UModal>
</template>
