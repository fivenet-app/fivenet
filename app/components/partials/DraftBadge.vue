<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';

defineProps</* @vue-ignore */ BadgeProps>();

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip :text="$t('common.draft')">
        <UBadge
            class="inline-flex grow-0 gap-1"
            color="info"
            size="md"
            icon="i-mdi-pencil"
            :label="$t('common.draft')"
            v-bind="$attrs"
        />
    </UTooltip>
</template>
