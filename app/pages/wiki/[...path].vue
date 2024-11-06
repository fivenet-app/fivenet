<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getGRPCWikiClient } from '~/composables/grpc';
import type { Page } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});
definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
});

const { t } = useI18n();

const { can } = useAuth();

const route = useRoute();

const path = computed(() => route.path.slice(5));

const {
    data: page,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`wiki-page:${path.value}`, () => getPage(path.value), {
    watch: [path],
});

async function getPage(path: string): Promise<Page | undefined> {
    try {
        const call = getGRPCWikiClient().getPage({
            path: path,
        });
        const { response } = await call;

        if (response.page === undefined) {
            throw createError({ statusCode: 404, statusMessage: t('pages.notfound.page_not_found'), fatal: true });
        }

        return response.page;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const tocLinks = computed(() => []);

const surround = ref([]);
</script>

<template>
    <UDashboardPage class="h-full">
        <UDashboardPanel
            class="h-full flex-shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <UDashboardNavbar :title="$t('common.wiki')">
                <template #right>
                    <UButton v-if="can('WikiService.CreateOrUpdatePage')" color="gray" trailing-icon="i-mdi-plus">
                        {{ $t('common.page') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <div class="flex flex-1 flex-col px-8 py-0">
                <UPage>
                    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page')])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.page')])"
                        :retry="refresh"
                    />

                    <UPageHeader v-else-if="page?.meta" :title="page.meta.title" :description="page.meta.description">
                        <template #links>
                            <UTooltip v-if="can('WikiService.CreateOrUpdatePage').value" :text="$t('common.edit')">
                                <UButton color="white" icon="i-mdi-pencil" />
                            </UTooltip>
                            <UTooltip v-if="can('WikiService.DeletePage').value" :text="$t('common.delete')">
                                <UButton color="red" icon="i-mdi-trash-can" />
                            </UTooltip>
                        </template>
                    </UPageHeader>

                    <UPageBody prose>
                        <!-- eslint-disable vue/no-v-html -->
                        <div v-if="page" class="prose dark:prose-invert" v-html="page.content" />

                        <hr v-if="surround?.length" />

                        <UContentSurround :surround="surround" />
                    </UPageBody>

                    <template v-if="page?.meta?.toc === undefined || page?.meta?.toc !== false" #right>
                        <UContentToc :title="$t('common.toc')" :links="tocLinks"> </UContentToc>
                    </template>
                </UPage>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
