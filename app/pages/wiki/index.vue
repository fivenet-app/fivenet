<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import PageSearch from '~/components/wiki/PageSearch.vue';
import { getWikiWikiClient } from '~~/gen/ts/clients';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'wiki.WikiService/ListPages',
});

const { activeChar, can } = useAuth();

const wikiWikiClient = await getWikiWikiClient();

const { data: pages, status, refresh, error } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = wikiWikiClient.listPages({
            pagination: {
                offset: 0,
            },
            rootOnly: true,
        });
        const { response } = await call;

        const pages = response.pages.sort((a, b) => a.job.localeCompare(b.job));
        if (pages.length > 0) {
            const ownPageIdx = pages.findIndex((p) => p.job === activeChar.value?.job);
            pages.unshift(pages.splice(ownPageIdx, 1)[0]!);
        }

        return pages;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(pages, async () => {
    if (!pages.value) {
        return;
    }

    if (pages.value.length === 1 && pages.value[0]?.job !== undefined && pages.value[0]?.job === activeChar.value?.job) {
        await navigateTo({
            name: 'wiki-job-id-slug',
            params: {
                job: pages.value[0]!.job,
                id: pages.value[0]!.id,
                slug: [pages.value[0]!.slug ?? ''],
            },
        });
    }
});

const wikiService = await useWikiWiki();
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('common.wiki')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #default>
                    <PageSearch />
                </template>

                <template #right>
                    <UTooltip v-if="can('wiki.WikiService/UpdatePage').value" :text="$t('common.create')">
                        <UButton color="neutral" trailing-icon="i-mdi-plus" @click="wikiService.createPage()">
                            {{ $t('common.page') }}
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar class="flex lg:hidden">
                <template #default>
                    <PageSearch />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.page')])" />
            <DataErrorBlock v-else-if="error" :retry="refresh" />
            <DataNoDataBlock
                v-else-if="!pages"
                icon="i-mdi-file-search"
                :title="$t('common.unable_to_load', [$t('common.wiki', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="pages.length === 0"
                icon="i-mdi-file-search"
                :title="$t('common.not_found', [$t('common.wiki', 2)])"
                :retry="refresh"
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
                    <template v-if="p.rootInfo?.logo?.filePath" #icon>
                        <GenericImg class="h-10 w-auto" :src="p.rootInfo?.logo?.filePath" :alt="$t('common.logo')" />
                    </template>
                </UPageCard>
            </UPageGrid>
        </template>
    </UDashboardPanel>
</template>
