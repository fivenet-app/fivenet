<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';

interface Props extends /* @vue-ignore */ BadgeProps {
    closed: boolean | undefined;
}

defineProps<Props>();
</script>

<template>
    <UBadge
        :color="closed ? 'error' : 'success'"
        :icon="closed ? 'i-mdi-lock' : 'i-mdi-lock-open-variant'"
        :label="closed ? $t('common.close', 2) : $t('common.open', 2)"
        v-bind="$attrs"
    />
</template>
