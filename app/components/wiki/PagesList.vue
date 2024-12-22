<script lang="ts" setup>
import type { NavItem } from '@nuxt/content';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

const props = defineProps<{
    pages: PageShort[];
}>();

const { t } = useI18n();

function mapNavItemToNavItem(page: PageShort): NavItem {
    const fullPath = `/wiki/${page.job}/${page.id}/${page.slug ?? ''}`;
    return {
        _id: page.id,
        title: page.title !== '' ? page.title : `${page?.jobLabel ? page?.jobLabel + ': ' : ''}${t('common.wiki')}`,
        _path: fullPath,
        children: page.children.map((p) => mapNavItemToNavItem(p)),
        icon: page.deletedAt !== undefined ? 'i-mdi-trash-can' : undefined,
    };
}

const navItems = computed(() => props.pages.map((p) => mapNavItemToNavItem(p)) ?? []);
</script>

<template>
    <UNavigationTree class="mt-2 sm:mt-0" :links="mapContentNavigation(navItems)" />
</template>
