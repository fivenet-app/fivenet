<script lang="ts" setup>
import type { BadgeColor } from '#ui/types';
import type { Category } from '~~/gen/ts/resources/documents/category';

const props = defineProps<{
    category: Category | undefined;
}>();

const color = computed(() => (props.category?.color ? (props.category?.color as BadgeColor) : 'primary'));

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UBadge v-if="category" class="inline-flex flex-initial gap-1" size="md" :color="color" v-bind="$attrs">
        <UIcon class="size-5" :name="category.icon ? convertComponentIconNameToDynamic(category.icon) : 'i-mdi-shape'" />
        <UTooltip :text="category.description">
            {{ category.name }}
        </UTooltip>
    </UBadge>
</template>
