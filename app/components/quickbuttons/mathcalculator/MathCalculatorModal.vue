<script lang="ts" setup>
import GenericModal from '~/components/partials/elements/GenericModal.vue';
import MathCalculator from '~/components/quickbuttons/mathcalculator/MathCalculator.vue';
import { useSettingsStore } from '~/stores/settings';

const { isOpen } = useModal();

const settingsStore = useSettingsStore();
const { calculatorPosition: position } = storeToRefs(settingsStore);

const containerPosition = computed(() => {
    switch (position.value) {
        case 'top':
            return 'sm:items-start justify-end';
        case 'bottom':
            return 'sm:items-end justify-end';
        default:
            return 'sm:items-center justify-end';
    }
});
</script>

<template>
    <GenericModal
        :open="isOpen"
        :title="$t('components.mathcalculator.title')"
        :ui="{
            container: containerPosition,
            width: 'w-full sm:max-w-md',
        }"
        :overlay="false"
        @close="isOpen = false"
    >
        <div class="flex gap-2">
            <MathCalculator class="flex-1" />

            <UButtonGroup class="my-auto flex-initial" orientation="vertical">
                <UButton icon="i-mdi-arrow-up-bold" @click="position = 'top'" />
                <UButton color="black" icon="i-mdi-format-vertical-align-center" @click="position = 'middle'" />
                <UButton icon="i-mdi-arrow-down-bold" @click="position = 'bottom'" />
            </UButtonGroup>
        </div>
    </GenericModal>
</template>
