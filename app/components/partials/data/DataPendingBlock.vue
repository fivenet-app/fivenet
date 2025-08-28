<script lang="ts" setup>
withDefaults(
    defineProps<{
        message: string;
        paused?: boolean;
    }>(),
    {
        paused: false,
    },
);
</script>

<template>
    <UAlert
        class="my-2"
        :title="message"
        variant="subtle"
        block
        :icon="paused ? 'i-mdi-pause' : 'i-mdi-refresh'"
        :ui="{ icon: paused ? '' : 'animate-spin' }"
    />
</template>
