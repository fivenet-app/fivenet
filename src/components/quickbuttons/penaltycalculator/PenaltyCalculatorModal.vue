<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import PenaltyCalculator from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

const { isOpen } = useModal();
</script>

<template>
    <GenericModal :open="isOpen" :title="$t('components.penaltycalculator.title')" @close="isOpen = false">
        <PenaltyCalculator />
    </GenericModal>
</template>
