<script lang="ts" setup>
import type { TypedRouteFromName } from '#build/typed-router';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import PagesList from '~/components/wiki/PagesList.vue';
import PageView from '~/components/wiki/PageView.vue';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'wiki.WikiService.ListPages',
    validate: async (route) => {
        route = route as TypedRouteFromName<'wiki-job-id-slug'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const { $grpc } = useNuxtApp();

const { activeChar } = useAuth();

const route = useRoute('wiki-job-id-slug');

const {
    data: pages,
    error: pagesError,
    refresh: pagesRefresh,
} = useLazyAsyncData(`wiki-pages:${route.path}`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    const job = route.params.job ?? activeChar.value?.job ?? '';
    try {
        const call = $grpc.wiki.wiki.listPages({
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
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`wiki-page:${route.path}`, () => getPage(parseInt(route.params.id)), {
    watch: [() => route.path],
});

async function getPage(id: number): Promise<Page | undefined> {
    try {
        const call = $grpc.wiki.wiki.getPage({
            id: id,
        });
        const { response } = await call;

        return response.page;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <PageView :loading="loading" :error="error" :refresh="refresh" :page="page" :pages="pages ?? []">
                <template #left>
                    <DataErrorBlock v-if="pagesError" :error="pagesError" :retry="pagesRefresh" />
                    <ClientOnly v-else>
                        <PagesList :pages="pages ?? []" />

                        <UTooltip :text="$t('common.refresh')">
                            <UButton class="-ml-2 mt-1" variant="link" icon="i-mdi-refresh" @click="pagesRefresh" />
                        </UTooltip>
                    </ClientOnly>
                </template>
            </PageView>
        </UDashboardPanel>
    </UDashboardPage>
</template>
