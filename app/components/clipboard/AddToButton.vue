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
        <UButton icon="i-mdi-clipboard-plus" color="neutral" @click="callback()" />
    </UTooltip>
</template>
