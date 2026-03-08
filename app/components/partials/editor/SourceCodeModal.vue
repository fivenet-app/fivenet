<script lang="ts" setup>
defineProps<{
    disabled?: boolean;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const content = defineModel<string>('content', { required: true });
</script>

<template>
    <UModal :title="$t('common.source_code')" fullscreen :ui="{ body: 'flex flex-col flex-1 overflow-y-hidden' }">
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
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
