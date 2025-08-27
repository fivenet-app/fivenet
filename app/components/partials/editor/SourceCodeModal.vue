<script lang="ts" setup>
const props = defineProps<{
    content: string;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:content', val: string): void;
}>();

const content = useVModel(props, 'content', emit);
</script>

<template>
    <UModal :title="$t('common.source_code')" fullscreen :ui="{ body: 'flex flex-col flex-1' }">
        <template #body>
            <UTextarea v-model="content" class="mx-auto size-full max-w-(--breakpoint-xl) flex-1" autoresize :row="0" />
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
