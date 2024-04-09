<script lang="ts" setup>
import { withoutTrailingSlash } from 'ufo';

useHead({
    title: 'common.help',
});
definePageMeta({
    title: 'common.help',
    layout: 'help',
    requiresAuth: false,
    showCookieOptions: true,
});

const route = useRoute();

const { data: page } = await useAsyncData(route.path, () => queryContent(route.path).findOne());
if (!page.value) {
    throw createError({ statusCode: 404, statusMessage: 'Page not found', fatal: true });
}

const { data: surround } = await useAsyncData(`${route.path}-surround`, () =>
    queryContent()
        .where({ _extension: 'md', navigation: { $ne: false } })
        .only(['title', 'description', '_path'])
        .findSurround(withoutTrailingSlash(route.path)),
);

const headline = computed(() => findPageHeadline(page.value!));
</script>

<template>
    <UPage v-if="page">
        <UPageHeader
            :title="page.title?.includes('.') ? $t(page.title) : page?.title"
            :description="page.description"
            :headline="headline"
        />

        <UPageBody prose>
            <ContentRenderer v-if="page.body" :value="page" />

            <hr v-if="surround?.length" />

            <UContentSurround v-if="surround" :surround="surround" />
        </UPageBody>

        <template v-if="page.toc !== false" #right>
            <UContentToc :title="$t('common.toc')" :links="page.body?.toc?.links" />
        </template>
    </UPage>
</template>
