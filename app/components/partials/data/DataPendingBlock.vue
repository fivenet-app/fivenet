<script lang="ts" setup>
import type { AlertProps } from '@nuxt/ui';

withDefaults(
    defineProps<{
        message: NonNullable<AlertProps['title']>;
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
