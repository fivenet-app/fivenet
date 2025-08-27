<script lang="ts" setup>
const emit = defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :title="$t('components.mathcalculator.title')"
        :overlay="false"
        :close="{ onClick: () => emit('close', false) }"
        side="bottom"
        handle-only
    >
        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsMathcalculatorMathCalculator class="w-40 min-w-1/4" />
            </div>
        </template>
    </UDrawer>
</template>
