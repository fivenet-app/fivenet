<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import type { Category } from '~~/gen/ts/resources/documents/category';

const props = defineProps<{
    category: Category | undefined;
}>();

const color = computed(() => (props.category?.color ? (props.category?.color as BadgeProps['color']) : 'primary'));

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UBadge
        v-if="category"
        class="inline-flex flex-initial gap-1"
        size="md"
        :color="color"
        :icon="category.icon ? convertComponentIconNameToDynamic(category.icon) : 'i-mdi-shape'"
        v-bind="$attrs"
    >
        <UTooltip :text="category.description">
            {{ category.name }}
        </UTooltip>
    </UBadge>
</template>
