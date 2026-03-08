<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        requireExplicitSave?: boolean;
        onSave?: () => void;
        onCancel?: () => void;
    }>(),
    {
        requireExplicitSave: false,
        onSave: undefined,
        onCancel: undefined,
    },
);

const emits = defineEmits<{
    close: [boolean];
}>();

function closeDrawer(): void {
    if (props.requireExplicitSave) props.onCancel?.();
    emits('close', false);
}

function saveDrawer(): void {
    props.onSave?.();
    emits('close', false);
}
</script>

<template>
    <UDrawer
        :title="$t('components.penaltycalculator.title')"
        :overlay="false"
        :close="{ onClick: () => closeDrawer() }"
        side="bottom"
        handle-only
        :ui="{ title: 'flex gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('components.penaltycalculator.title') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="closeDrawer" />
        </template>

        <template #body>
            <div class="flex justify-center">
                <LazyQuickbuttonsPenaltycalculatorPenaltyCalculator class="w-full max-w-[80%] min-w-1/2" />
            </div>
        </template>

        <template v-if="props.requireExplicitSave" #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.cancel')" @click="closeDrawer" />
                <UButton class="flex-1" block :label="$t('common.save')" @click="saveDrawer" />
            </UFieldGroup>
        </template>
    </UDrawer>
</template>
