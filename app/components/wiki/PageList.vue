<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';

defineProps<{
    items: NavigationMenuItem[];
}>();
</script>

<template>
    <UNavigationMenu
        orientation="vertical"
        :items="items"
        class="mt-1 sm:mt-0"
        highlight
        color="primary"
        :ui="{ link: 'data-[active=]:font-bold' }"
    />
</template>
