<script lang="ts" setup>
withDefaults(
    defineProps<{
        callback: () => void;
        title?: string;
    }>(),
    {
        title: '',
    },
);
</script>

<template>
    <UTooltip :text="title">
        <UButton :ui="{ rounded: 'rounded-full' }" icon="i-mdi-clipboard-plus" color="gray" @click="callback()" />
    </UTooltip>
</template>
