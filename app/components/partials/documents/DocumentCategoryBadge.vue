<script lang="ts" setup>
import type { Category } from '~~/gen/ts/resources/documents/category';

defineProps<{
    category: Category | undefined;
}>();
</script>

<template>
    <UBadge v-if="category" class="inline-flex flex-initial gap-1" size="md">
        <UIcon name="i-mdi-shape" class="size-5" />
        <span :title="category.description ?? $t('common.na')">
            {{ category.name }}
        </span>
    </UBadge>
</template>
