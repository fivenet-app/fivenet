<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getGRPCWikiClient } from '~/composables/grpc';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});
definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
});

const { can } = useAuth();

const { data: pages } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = getGRPCWikiClient().listPages({
            pagination: {
                offset: 0,
            },
            rootOnly: true,
        });
        const { response } = await call;

        return response.pages;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(pages, async () => {
    if (!pages.value) {
        return;
    }

    if (pages.value.length === 1) {
        await navigateTo({
            name: 'wiki-job-id-slug',
            params: {
                job: pages.value[0]!.job,
                id: pages.value[0]!.id,
                slug: pages.value[0]!.slug,
            },
        });
    }
});
</script>

<template>
    <UDashboardPage class="h-full">
        <UDashboardPanel
            class="h-full flex-shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <UDashboardNavbar :title="$t('common.wiki')">
                <template #right>
                    <UButton v-if="can('WikiService.CreatePage')" color="gray" trailing-icon="i-mdi-plus" to="/wiki/create">
                        {{ $t('common.page') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <DataNoDataBlock
                    v-if="!pages || pages.length === 0"
                    icon="i-mdi-file-search"
                    :message="$t('common.not_found', [$t('common.wiki', 2)])"
                />

                <UPageGrid v-else>
                    <UPageCard
                        v-for="p in pages"
                        :key="p.id"
                        :title="`${p.jobLabel} ${$t('common.wiki')}`"
                        :to="`/wiki/${p.job}/${p.id}/${p.slug ?? ''}`"
                        icon="i-mdi-brain"
                        :ui="{ title: 'text-xl' }"
                    >
                        <template v-if="p.rootInfo?.logo?.url" #icon>
                            <img :src="p.rootInfo?.logo?.url" class="h-10 w-10" :alt="$t('common.logo')" />
                        </template>
                    </UPageCard>
                </UPageGrid>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
