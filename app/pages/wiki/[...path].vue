<script lang="ts" setup>
useHead({
    title: 'common.wiki',
});
definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
});

const body = ref('## TEST\n Hello world!');

const parsedBody = computedAsync(async () => await parseMarkdown(body.value));

const page = computed(() => ({
    _id: '123',
    title: 'Home',
    description: 'Home page of Wiki.',
    links: [],
    body: parsedBody.value?.body,
    toc: true,
}));

if (!page.value) {
    throw createError({ statusCode: 404, statusMessage: 'Page not found', fatal: true });
}

const surround = ref([]);
</script>

<template>
    <UDashboardPage class="h-full">
        <UDashboardPanel
            class="h-full flex-shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <UDashboardNavbar :title="$t('common.wiki')"> </UDashboardNavbar>

            <UDashboardPanelContent>
                <UPage>
                    <UPageHeader :title="page.title" :description="page.description" :links="page.links" />

                    <UPageBody prose>
                        <ContentRenderer v-if="page.body" :value="page" />

                        <hr v-if="surround?.length" />

                        <UContentSurround :surround="surround" />
                    </UPageBody>

                    <template v-if="page.toc !== false" #right>
                        <UContentToc :title="$t('common.toc')" :links="parsedBody?.toc?.links"> </UContentToc>
                    </template>
                </UPage>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
