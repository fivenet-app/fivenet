<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';

defineProps<{
    items: NavigationMenuItem[];
}>();
</script>

<template>
    <UNavigationMenu
        class="mt-1 sm:mt-0"
        orientation="vertical"
        :items="items"
        highlight
        color="primary"
        :ui="{ link: 'data-[active=]:font-bold' }"
    />
</template>
