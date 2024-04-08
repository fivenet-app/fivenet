<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import PenaltyCalculator from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

defineProps<{
    open: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <GenericModal :open="open" :title="$t('components.penaltycalculator.title')" @close="$emit('close')">
        <PenaltyCalculator :load-laws="open" />
    </GenericModal>
</template>
