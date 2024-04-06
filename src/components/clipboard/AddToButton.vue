<script lang="ts" setup>
withDefaults(
    defineProps<{
        callback: Function;
        title?: string;
    }>(),
    {
        title: '',
    },
);
</script>

<template>
    <UButton :ui="{ rounded: 'rounded-full' }" icon="i-mdi-clipboard-plus" color="gray" :title="title" @click="callback()" />
</template>
