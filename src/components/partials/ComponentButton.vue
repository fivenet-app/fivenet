<script lang="ts" setup>
import { CalculatorIcon, QuestionMarkCircleIcon } from '@heroicons/vue/24/outline';
import { FunctionalComponent } from 'vue';
import Calculator from '~/components/penaltycalculator/Calculator.vue';
import ComponentModal from './ComponentModal.vue';

export type ButtonComponents = 'PenaltyCalculator';

const props = defineProps<{
    comp: ButtonComponents,
}>();

const open = ref(false);

let icon: FunctionalComponent = QuestionMarkCircleIcon;
switch (props.comp) {
    case 'PenaltyCalculator':
        icon = CalculatorIcon;
        break;
    default:
        break;
}
</script>

<template>
    <button type="button" @click="open = true"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-10 right-24 bg-primary-500 shadow-float text-neutral hover:bg-primary-400">
        <component :is="icon" class="w-10 h-auto" />
    </button>
    <ComponentModal :open="open" @close="open = false">
        <Calculator v-if="comp === 'PenaltyCalculator'" />
    </ComponentModal>
</template>
