<script lang="ts" setup>
defineProps<{
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const content = defineModel<string>({
    required: true,
});

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(content);

async function closeModal(): Promise<void> {
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="$t('components.fabric_editor.import.title')"
        fullscreen
        :ui="{ body: 'flex flex-col flex-1 overflow-y-hidden' }"
        :close="false"
        :dismissible="!hasUnsavedChanges"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.fabric_editor.import.title') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UTextarea
                v-model="content"
                class="mx-auto h-full w-full max-w-(--breakpoint-xl) overflow-y-hidden"
                :disabled="disabled"
                :autoresize="false"
                :row="0"
                :ui="{ base: '!resize-none h-full w-full' }"
            />
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="closeModal" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
