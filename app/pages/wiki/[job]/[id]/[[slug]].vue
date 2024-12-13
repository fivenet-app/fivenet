<script lang="ts" setup>
import type { TypedRouteFromName } from '#build/typed-router';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import PageEditor from '~/components/wiki/PageEditor.vue';
import PagesList from '~/components/wiki/PagesList.vue';
import PageView from '~/components/wiki/PageView.vue';
import { getGRPCWikiClient } from '~/composables/grpc';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});
definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'WikiService.ListPages',
    validate: async (route) => {
        route = route as TypedRouteFromName<'wiki-job-id-slug'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return idParamRegex.test(route.params.id as string);
    },
});

const { activeChar } = useAuth();

const route = useRoute('wiki-job-id-slug');

const { data: pages, error: pagesError, refresh: pagesRefresh } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = getGRPCWikiClient().listPages({
            pagination: {
                offset: 0,
            },
            job: useRoute('wiki-job-id-slug').params.job ?? activeChar.value?.job ?? '',
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
} = useLazyAsyncData(`wiki-page:${route.path}`, () => getPage(route.params.id), {
    watch: [() => route.path],
});

async function getPage(id: string): Promise<Page | undefined> {
    try {
        const call = getGRPCWikiClient().getPage({
            id: id,
        });
        const { response } = await call;

        return response.page;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const editing = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <PageView
                v-if="!editing"
                :loading="loading"
                :error="error"
                :refresh="refresh"
                :page="page"
                :pages="pages ?? []"
                @edit="editing = !editing"
            >
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
            <PageEditor v-else v-model="page" :pages="pages ?? []" @close="editing = !editing" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
