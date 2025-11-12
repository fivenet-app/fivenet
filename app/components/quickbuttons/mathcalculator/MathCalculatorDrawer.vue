<script lang="ts" setup>
defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :title="$t('components.mathcalculator.title')"
        :overlay="false"
        :close="{ onClick: () => $emit('close', false) }"
        side="bottom"
        handle-only
        :ui="{ title: 'flex flex-row gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('components.mathcalculator.title') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsMathcalculatorMathCalculator class="w-48 min-w-3/4 sm:min-w-2/4" />
            </div>
        </template>
    </UDrawer>
</template>
