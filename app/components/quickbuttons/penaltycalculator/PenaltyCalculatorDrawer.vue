<script lang="ts" setup>
const emit = defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :title="$t('components.penaltycalculator.title')"
        :overlay="false"
        :close="{ onClick: () => emit('close', false) }"
        side="bottom"
        handle-only
    >
        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsPenaltycalculatorPenaltyCalculator class="w-full max-w-[80%] min-w-1/2" />
            </div>
        </template>
    </UDrawer>
</template>
