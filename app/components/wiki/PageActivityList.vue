<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import PageActivityListEntry from '~/components/wiki/PageActivityListEntry.vue';
import type { ListPageActivityResponse } from '~~/gen/ts/services/wiki/wiki';

const props = defineProps<{
    pageId: number;
}>();

const { $grpc } = useNuxtApp();

const page = useRouteQuery('page', '1', { transform: Number });

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`wiki-page:${props.pageId}-${page.value}-${page.value}`, () => listPageActivity());

async function listPageActivity(): Promise<ListPageActivityResponse> {
    try {
        const call = $grpc.wiki.wiki.listPageActivity({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            pageId: props.pageId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.activity', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.activity', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="!data || data.activity.length === 0"
            icon="i-mdi-ticket"
            :message="$t('common.not_found', [$t('common.activity')])"
        />

        <ul v-else class="mb-1 divide-y divide-gray-100 dark:divide-gray-800" role="list">
            <PageActivityListEntry v-for="item in data.activity" :key="item.id" :entry="item" />
        </ul>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
