<script lang="ts" setup>
import { CalculatorIcon } from 'mdi-vue3';
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import Calculator from '~/components/penaltycalculator/PenaltyCalculator.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-neutral shadow-float hover:bg-primary-400"
        @click="open = true"
    >
        <CalculatorIcon class="w-10 h-auto" />
    </button>

    <GenericModal :open="open" @close="open = false">
        <Calculator />
    </GenericModal>
</template>
