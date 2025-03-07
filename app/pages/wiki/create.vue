<script lang="ts" setup>
import PageEditor from '~/components/wiki/PageEditor.vue';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';

useHead({
    title: 'common.wiki',
});

definePageMeta({
    title: 'common.wiki',
    requiresAuth: true,
    permission: ['WikiService.CreatePage'],
});

const { activeChar } = useAuth();

const { data: pages } = useLazyAsyncData(`wiki-pages`, () => listPages());

async function listPages(): Promise<PageShort[]> {
    try {
        const call = getGRPCWikiClient().listPages({
            pagination: {
                offset: 0,
            },
            job: activeChar.value?.job ?? '',
        });
        const { response } = await call;

        return response.pages;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <PageEditor :model-value="undefined" :pages="pages ?? []" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
