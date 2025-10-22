<script lang="ts" setup>
import type { ContentNavigationItem } from '@nuxt/content';

defineProps<{
    navItems: ContentNavigationItem[];
}>();
</script>

<template>
    <UContentNavigation class="mt-1 sm:mt-0" :navigation="navItems" type="multiple" highlight default-open />
</template>
