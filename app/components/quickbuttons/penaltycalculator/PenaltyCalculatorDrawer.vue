<script lang="ts" setup>
const emit = defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :overlay="false"
        :close="{ onClick: () => emit('close', false) }"
        side="bottom"
        handle-only
        :ui="{ title: 'inline-flex w-full items-center gap-2 justify-between' }"
    >
        <template #title>
            <span>{{ $t('components.penaltycalculator.title') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="emit('close', false)" />
        </template>

        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsPenaltycalculatorPenaltyCalculator class="w-full max-w-[80%] min-w-1/2" />
            </div>
        </template>
    </UDrawer>
</template>
