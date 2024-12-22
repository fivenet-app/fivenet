<script lang="ts" setup>
const props = defineProps<{
    content: string;
}>();

const emit = defineEmits<{
    (e: 'update:content', val: string): void;
}>();

const { isOpen } = useModal();

const content = useVModel(props, 'content', emit);
</script>

<template>
    <UModal fullscreen>
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.source_code') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div class="mx-auto flex w-full max-w-screen-xl flex-1 flex-col">
                <UTextarea v-model="content" autoresize class="flex flex-1 flex-col" :ui="{ base: 'flex-1' }" />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton block class="flex-1" color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
