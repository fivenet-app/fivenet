<script lang="ts" setup>
import QualificationResultTutorForm from './QualificationResultTutorForm.vue';

defineProps<{
    qualificationId: string;
    userId: number;
}>();

defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <QualificationResultTutorForm
            :qualification-id="qualificationId"
            :user-id="userId"
            @refresh="$emit('refresh')"
            @close="isOpen = false"
        />
    </UModal>
</template>
