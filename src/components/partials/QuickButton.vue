<script lang="ts" setup>
import { CalculatorIcon, HelpCircleIcon } from 'mdi-vue3';
import { DefineComponent } from 'vue';
import ComponentModal from '~/components/partials/elements/ComponentModal.vue';
import Calculator from '~/components/penaltycalculator/Calculator.vue';

export type ButtonComponents = 'PenaltyCalculator';

const props = defineProps<{
    comp: ButtonComponents;
}>();

let icon: DefineComponent;
switch (props.comp) {
    case 'PenaltyCalculator':
        icon = CalculatorIcon;
        break;
    default:
        icon = HelpCircleIcon;
        break;
}

const open = ref(false);
</script>

<template>
    <button
        type="button"
        @click="open = true"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-2 right-20 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
    >
        <component :is="icon" class="w-10 h-auto" />
    </button>
    <ComponentModal :open="open" @close="open = false">
        <Calculator v-if="comp === 'PenaltyCalculator'" />
    </ComponentModal>
</template>
