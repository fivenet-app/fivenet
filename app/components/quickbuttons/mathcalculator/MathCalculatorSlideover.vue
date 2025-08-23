<script lang="ts" setup>
const emit = defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <USlideover
        :title="$t('components.mathcalculator.title')"
        :overlay="false"
        :close="{ onClick: () => emit('close', false) }"
        side="bottom"
    >
        <template #body>
            <LazyQuickbuttonsMathcalculatorMathCalculator />
        </template>
    </USlideover>
</template>
