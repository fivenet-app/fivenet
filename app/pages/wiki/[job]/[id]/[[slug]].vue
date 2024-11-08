<script lang="ts" setup>
import type { TypedRouteFromName } from '#build/typed-router';
import PageEditor from '~/components/wiki/PageEditor.vue';
import PageView from '~/components/wiki/PageView.vue';
import { getGRPCWikiClient } from '~/composables/grpc';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});
definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
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

const { data: pages } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = getGRPCWikiClient().listPages({
            pagination: {
                offset: 0,
            },
            job: page.value?.job ?? activeChar.value?.job ?? '',
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
    <UDashboardPage class="h-full">
        <UDashboardPanel
            class="h-full flex-shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <PageView
                v-if="!editing"
                :loading="loading"
                :error="error"
                :refresh="refresh"
                :page="page"
                :pages="pages ?? []"
                @edit="editing = !editing"
            />
            <PageEditor v-else v-model="page" :pages="pages ?? []" @close="editing = !editing" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
