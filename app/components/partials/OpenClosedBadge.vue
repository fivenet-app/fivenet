<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';

interface Props extends /* @vue-ignore */ BadgeProps {
    closed: boolean | undefined;
}

defineProps<Props>();
</script>

<template>
    <UBadge v-if="closed" color="error" icon="i-mdi-lock" :label="$t('common.close', 2)" v-bind="$attrs" />
    <UBadge v-else color="green" icon="i-mdi-lock-open-variant" :label="$t('common.open', 2)" v-bind="$attrs" />
</template>
