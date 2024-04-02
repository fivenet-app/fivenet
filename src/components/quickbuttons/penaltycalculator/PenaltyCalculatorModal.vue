<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import PenaltyCalculator from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

const open = ref(true);
</script>

<template>
    <GenericModal :open="open" :unmount="false" :title="$t('components.penaltycalculator.title')" @close="open = false">
        <PenaltyCalculator :load-laws="open" />
    </GenericModal>
</template>
