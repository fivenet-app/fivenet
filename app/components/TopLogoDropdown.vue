<script lang="ts" setup>
defineProps<{
    collapsed?: boolean;
}>();

const version = APP_VERSION;
</script>

<template>
    <UTooltip :text="version">
        <UButton
            class="w-full"
            color="neutral"
            variant="ghost"
            :avatar="{
                src: '/images/logo.webp',
                alt: 'FiveNet',
            }"
            :ui="{ base: collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)', leadingAvatar: 'rounded-none' }"
        >
            <span v-if="!collapsed" class="truncate font-semibold text-highlighted">FiveNet</span>
        </UButton>
    </UTooltip>
</template>
