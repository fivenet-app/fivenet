<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import MathCalculator from '~/components/quickbuttons/mathcalculator/MathCalculator.vue';

const { isOpen } = useModal();
</script>

<template>
    <GenericModal
        :open="isOpen"
        :title="$t('components.mathcalculator.title')"
        :ui="{ width: 'w-full sm:max-w-md' }"
        @close="isOpen = false"
    >
        <MathCalculator />
    </GenericModal>
</template>
