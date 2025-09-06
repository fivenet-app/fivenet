<script lang="ts" setup>
import type { ContentNavigationItem } from '@nuxt/content';
import { mapContentNavigation } from '@nuxt/ui-pro/utils/content';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

const props = defineProps<{
    pages: PageShort[];
}>();

const { t } = useI18n();

function mapNavItemToNavItem(page: PageShort): ContentNavigationItem {
    return {
        id: page.id.toString(),
        title: page.title !== '' ? page.title : `${page?.jobLabel ? page?.jobLabel + ': ' : ''}${t('common.wiki')}`,
        path: `/wiki/${page.job}/${page.id}/${page.slug ?? ''}`,
        icon: page.deletedAt !== undefined ? 'i-mdi-delete' : page.draft ? 'i-mdi-pencil' : undefined,
        children: page.children.map((p) => mapNavItemToNavItem(p)),
    };
}

const navItems = computed(() => props.pages.map((p) => mapNavItemToNavItem(p)) ?? []);
</script>

<template>
    <UContentNavigation class="mt-1 sm:mt-0" :navigation="mapContentNavigation(navItems)" type="multiple" />
</template>
