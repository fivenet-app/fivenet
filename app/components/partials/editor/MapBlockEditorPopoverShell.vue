<script lang="ts" setup>
const open = defineModel<boolean>('open', {
    required: true,
});

const props = withDefaults(
    defineProps<{
        title: string;
        contentWidthClass?: string;
    }>(),
    {
        contentWidthClass: 'w-[28rem]',
    },
);
</script>

<template>
    <UPopover v-model:open="open">
        <slot name="trigger" :open="open" />

        <template #content>
            <div class="flex max-w-[calc(100vw-2rem)] flex-col gap-3 p-4" :class="props.contentWidthClass">
                <h3 class="block font-medium">{{ props.title }}</h3>

                <slot />
            </div>
        </template>
    </UPopover>
</template>
