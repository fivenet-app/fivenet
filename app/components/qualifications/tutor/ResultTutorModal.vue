<script lang="ts" setup>
import ResultTutorForm from './ResultTutorForm.vue';

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
    <ResultTutorForm
        :qualification-id="qualificationId"
        :user-id="userId"
        :result-id="resultId"
        @refresh="$emit('refresh')"
        @close="$emit('close', false)"
    />
</template>
