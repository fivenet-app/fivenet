<script lang="ts" setup>
import { CalculatorIcon } from 'mdi-vue3';
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import Calculator from '~/components/penaltycalculator/PenaltyCalculator.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-neutral shadow-float hover:bg-primary-400"
        @click="open = true"
    >
        <CalculatorIcon class="h-auto w-10" />
    </button>

    <ComponentModal :open="open" @close="open = false">
        <Calculator />
    </ComponentModal>
</template>
