<script lang="ts" setup>
import { CalculatorIcon } from 'mdi-vue3';
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import Calculator from '~/components/penaltycalculator/Calculator.vue';

const open = ref(false);
</script>

<template>
    <button
        type="button"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-10 bottom-2 right-20 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
        @click="open = true"
    >
        <CalculatorIcon class="w-10 h-auto" />
    </button>
    <ComponentModal :open="open" @close="open = false">
        <Calculator />
    </ComponentModal>
</template>
