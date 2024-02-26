<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        url?: string;
        text?: string;
        size?: 'md' | 'lg' | 'xl';
    }>(),
    {
        url: undefined,
        text: '',
        size: 'lg',
    },
);

const size = computed(() => {
    switch (props.size) {
        case 'xl':
            return 'text-lg h-20 w-20';
        case 'lg':
            return 'text-lg h-12 w-12';
        case 'md':
        default:
            return 'text-base h-10 w-10';
    }
});
</script>

<template>
    <span class="flex items-center justify-center rounded-sm bg-gray-500 ring-2 ring-base-600" :class="size">
        <span v-if="!url" class="font-medium leading-none text-white">
            {{ text }}
        </span>
        <img v-else class="rounded-full" :class="size" :src="url" :alt="text" />
    </span>
</template>
