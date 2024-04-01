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
    <UButton
        :title="title"
        icon="i-mdi-clipboard-plus"
        :ui="{ rounded: 'rounded-full' }"
        class="fixed bottom-16 right-6 z-10 inline-flex size-12 items-center justify-center shadow-float"
        @click="callback()"
    >
    </UButton>
</template>
