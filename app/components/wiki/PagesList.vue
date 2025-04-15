<script lang="ts" setup>
import type { ContentNavigationLink } from '@nuxt/ui-pro/runtime/components/content/ContentNavigation.vue';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

const props = defineProps<{
    pages: PageShort[];
}>();

const { t } = useI18n();

function mapNavItemToNavItem(page: PageShort): ContentNavigationLink {
    const fullPath = `/wiki/${page.job}/${page.id}/${page.slug ?? ''}`;
    return {
        _id: page.id.toString(),
        title: page.title !== '' ? page.title : `${page?.jobLabel ? page?.jobLabel + ': ' : ''}${t('common.wiki')}`,
        path: fullPath,
        children: page.children.map((p) => mapNavItemToNavItem(p)),
        icon: page.deletedAt !== undefined ? 'i-mdi-delete' : undefined,
    };
}

const navItems = computed(() => props.pages.map((p) => mapNavItemToNavItem(p)) ?? []);
</script>

<template>
    <UContentNavigation class="mt-2 sm:mt-0" :items="navItems" />
</template>
