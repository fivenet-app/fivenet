<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import type { Category } from '~~/gen/ts/resources/documents/category/category';

const props = withDefaults(
    defineProps<{
        category: Category | undefined;
        size?: BadgeProps['size'];
    }>(),
    {
        size: 'md',
    },
);

const color = computed(() => (props.category?.color ? (props.category?.color as BadgeProps['color']) : 'primary'));

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip v-if="category" :text="category.description">
        <UBadge
            class="inline-flex flex-initial gap-1"
            :size="size"
            :color="color"
            :icon="category.icon ? convertComponentIconNameToDynamic(category.icon) : 'i-mdi-shape'"
            :label="category.name"
            v-bind="$attrs"
        />
    </UTooltip>
</template>
