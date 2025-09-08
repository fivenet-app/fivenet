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
    <div class="m-2">
        <UAlert
            :title="message"
            variant="subtle"
            block
            :icon="paused ? 'i-mdi-pause' : 'i-mdi-refresh'"
            :ui="{ icon: paused ? '' : 'animate-spin' }"
        />
    </div>
</template>
