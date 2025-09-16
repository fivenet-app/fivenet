<script lang="ts" setup>
defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :title="$t('components.penaltycalculator.title')"
        :overlay="false"
        :close="{ onClick: () => $emit('close', false) }"
        side="bottom"
        handle-only
        :ui="{ title: 'flex gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('components.penaltycalculator.title') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsPenaltycalculatorPenaltyCalculator class="w-full max-w-[80%] min-w-1/2" />
            </div>
        </template>
    </UDrawer>
</template>
