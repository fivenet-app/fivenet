<script lang="ts" setup>
import type { NavItem } from '@nuxt/content';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

function mapNavItemToNavItem(p: PageShort): NavItem {
    return {
        title: p.parentId === undefined && p.jobLabel !== undefined ? `${p.jobLabel}: ${p.title}` : p.title,
        _path: `/wiki/${p.job}/${p.id}/${p.slug ?? ''}`,
        children: p.children.map((p) => mapNavItemToNavItem(p)),
    } as NavItem;
}

const props = defineProps<{
    pages: PageShort[];
}>();

const navItems = computed(() => props.pages.map((p) => mapNavItemToNavItem(p)) ?? []);
</script>

<template>
    <UNavigationTree :links="mapContentNavigation(navItems)" />
</template>
