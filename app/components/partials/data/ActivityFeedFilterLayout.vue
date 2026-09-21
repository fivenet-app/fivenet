<script lang="ts" setup>
const props = defineProps<{
    fieldCount: 3 | 4;
}>();

const layoutClass = computed(() =>
    props.fieldCount === 4
        ? 'grid w-full gap-3 sm:grid-cols-2 lg:grid-cols-[repeat(3,minmax(0,1fr))_auto]'
        : 'grid w-full gap-3 sm:grid-cols-2 lg:grid-cols-[repeat(2,minmax(0,1fr))_auto]',
);
</script>

<template>
    <div :class="layoutClass">
        <slot />
    </div>
</template>
