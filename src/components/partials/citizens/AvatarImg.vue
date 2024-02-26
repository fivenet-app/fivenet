<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        avatarUrl?: string;
        firstname: string;
        lastname: string;
        size?: 'md' | 'lg';
    }>(),
    {
        avatarUrl: undefined,
        size: 'lg',
    },
);

const size = computed(() => (props.size === 'lg' ? 'text-lg h-12 w-12' : 'text-base h-10 w-10'));
</script>

<template>
    <span class="flex items-center justify-center rounded-full bg-gray-500 ring-2 ring-base-600" :class="size">
        <span v-if="!avatarUrl" class="font-medium leading-none text-white">
            {{ getInitials(`${firstname} ${lastname}`) }}
        </span>
        <img v-else class="rounded-full" :class="size" :src="avatarUrl" :alt="$t('common.avatar')" />
    </span>
</template>
