<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        title?: string;
        size?: 'sm' | 'md';
    }>(),
    {
        title: undefined,
        size: 'md',
    },
);

const titleClass = computed(() => (props.size === 'sm' ? 'text-base font-medium' : 'text-lg font-medium sm:text-xl'));
</script>

<template>
    <h2 :class="['line-clamp-2 min-w-0 break-words text-highlighted', titleClass]">
        <span v-if="!title" class="italic">
            {{ $t('common.untitled') }}
        </span>
        <span v-else>
            {{ title }}
        </span>
    </h2>
</template>
