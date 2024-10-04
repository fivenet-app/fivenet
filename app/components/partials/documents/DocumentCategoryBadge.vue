<script lang="ts" setup>
import type { BadgeColor } from '#ui/types';
import type { Category } from '~~/gen/ts/resources/documents/category';

const props = defineProps<{
    category: Category | undefined;
}>();

const color = computed(() => (props.category?.color ? (props.category?.color as BadgeColor) : 'primary'));
</script>

<template>
    <UBadge v-if="category" class="inline-flex flex-initial gap-1" size="md" :color="color">
        <UIcon :name="category.icon ? convertComponentIconNameToDynamic(category.icon) : 'i-mdi-shape'" class="size-5" />
        <span :title="category.description ?? $t('common.na')">
            {{ category.name }}
        </span>
    </UBadge>
</template>
