<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import PageSearch from '~/components/wiki/PageSearch.vue';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: 'WikiService.ListPages',
});

const { $grpc } = useNuxtApp();

const { activeChar, can } = useAuth();

const { data: pages, pending: loading, refresh, error } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = $grpc.wiki.wiki.listPages({
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
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.wiki')">
                <template #center>
                    <PageSearch />
                </template>

                <template #right>
                    <UButton
                        v-if="can('WikiService.CreatePage').value"
                        color="gray"
                        trailing-icon="i-mdi-plus"
                        to="/wiki/create"
                    >
                        {{ $t('common.page') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar class="!flex lg:!hidden">
                <template #default>
                    <PageSearch />
                </template>
            </UDashboardToolbar>

            <UDashboardPanelContent>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page')])" />
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
                        <template v-if="p.rootInfo?.logo?.url" #icon>
                            <NuxtImg :src="p.rootInfo?.logo?.url" :alt="$t('common.logo')" class="h-10 w-10" loading="lazy" />
                        </template>
                    </UPageCard>
                </UPageGrid>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
