<script lang="ts" setup>
import { CalculatorIcon } from 'mdi-vue3';
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import PenaltyCalculator from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex size-12 items-center justify-center rounded-full bg-accent-500 text-neutral shadow-float hover:bg-accent-400"
        @click="open = true"
    >
        <CalculatorIcon class="h-auto w-10" aria-hidden="true" />
    </button>

    <GenericModal :open="open" :unmount="false" @close="open = false">
        <PenaltyCalculator :load-laws="open" />
    </GenericModal>
</template>
