<script lang="ts" setup>
import type { TypedRouteFromName } from '#build/typed-router';
import type { ContentNavigationItem } from '@nuxt/content';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import PageList from '~/components/wiki/PageList.vue';
import PageView from '~/components/wiki/PageView.vue';
import { getWikiWikiClient } from '~~/gen/ts/clients';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'wiki.WikiService/ListPages',
    validate: async (route) => {
        route = route as TypedRouteFromName<'wiki-job-id-slug'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const { t } = useI18n();

const { activeChar } = useAuth();

const route = useRoute('wiki-job-id-slug');

const wikiWikiClient = await getWikiWikiClient();

const {
    data: pages,
    error: pagesError,
    status: pagesStatus,
    refresh: pagesRefresh,
} = useLazyAsyncData(`wiki-pages-${route.params.job}`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    const job = route.params.job ?? activeChar.value?.job ?? '';
    try {
        const call = wikiWikiClient.listPages({
            pagination: {
                offset: 0,
            },
            job: job,
            rootOnly: false,
        });
        const { response } = await call;

        return response.pages;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const {
    data: page,
    status,
    refresh,
    error,
} = useLazyAsyncData(`wiki-page:${route.path}`, () => getPage(parseInt(route.params.id)), {
    watch: [() => route.path],
});

async function getPage(id: number): Promise<Page | undefined> {
    try {
        const call = wikiWikiClient.getPage({
            id: id,
        });
        const { response } = await call;

        return response.page;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

useHead({
    title: () =>
        page.value?.meta?.title
            ? `${page.value.meta.title} - ${page.value.jobLabel} - ${t('pages.wiki.id.title')}`
            : t('pages.wiki.id.title'),
});

function mapNavItemToNavItem(page: PageShort): ContentNavigationItem {
    return {
        id: page.id.toString(),
        title: page.title !== '' ? page.title : `${page?.jobLabel ? page?.jobLabel + ': ' : ''}${t('common.wiki')}`,
        path: `/wiki/${page.job}/${page.id}/${page.slug ?? ''}`,
        icon: page.deletedAt !== undefined ? 'i-mdi-delete' : page.draft ? 'i-mdi-pencil' : undefined,
        children: page.children.map((p) => mapNavItemToNavItem(p)),
    };
}

const navItems = computed(() => pages.value?.map((p) => mapNavItemToNavItem(p)) ?? []);
</script>

<template>
    <PageView :status="status" :error="error" :refresh="refresh" :page="page" :pages="pages ?? []" :nav-items="navItems ?? []">
        <template #left>
            <DataErrorBlock v-if="pagesError" :error="pagesError" :retry="pagesRefresh" />
            <UPageAside v-else :ui="{ root: 'px-0 lg:pe-3' }">
                <ClientOnly>
                    <PageList :nav-items="navItems" />

                    <UTooltip :text="$t('common.refresh')">
                        <UButton
                            class="mt-1 -ml-2"
                            variant="link"
                            icon="i-mdi-refresh"
                            :loading="isRequestPending(pagesStatus)"
                            @click="() => pagesRefresh()"
                        />
                    </UTooltip>
                </ClientOnly>
            </UPageAside>
        </template>
    </PageView>
</template>
