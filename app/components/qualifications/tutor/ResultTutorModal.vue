<script lang="ts" setup>
import ResultTutorForm from './ResultTutorForm.vue';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';

defineProps<{
    qualificationId: number;
    qualification?: QualificationShort;
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
        :qualification="qualification"
        :user-id="userId"
        :result-id="resultId"
        :title="$t('components.qualifications.result_modal.add_title')"
        @refresh="$emit('refresh')"
        @close="$emit('close', false)"
    />
</template>
