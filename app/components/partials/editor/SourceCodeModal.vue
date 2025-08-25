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
    <UModal fullscreen>
        <template #title>
            <h3 class="text-2xl leading-6 font-semibold">
                {{ $t('common.source_code') }}
            </h3>
        </template>

        <template #body>
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 flex-col">
                <UTextarea v-model="content" class="flex flex-1 flex-col" autoresize :ui="{ base: 'flex-1' }" />
            </div>
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
