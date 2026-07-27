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

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <div class="p-2">
        <UAlert
            :title="message"
            variant="subtle"
            block
            :icon="paused ? 'i-mdi-pause' : 'i-mdi-refresh'"
            :ui="{ icon: paused ? '' : 'animate-spin' }"
            v-bind="$attrs"
        />
    </div>
</template>
