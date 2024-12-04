<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import PageActivityListEntry from '~/components/wiki/PageActivityListEntry.vue';
import type { ListPageActivityResponse } from '~~/gen/ts/services/wiki/wiki';

const props = defineProps<{
    pageId: string;
}>();

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`wiki-page:${props.pageId}-${page.value}`, () => listPageActivity());

async function listPageActivity(): Promise<ListPageActivityResponse> {
    try {
        const call = getGRPCWikiClient().listPageActivity({
            pagination: {
                offset: offset.value,
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

watch(offset, async () => refresh());
</script>

<template>
    <div>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.activity', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.activity', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="!data || data.activity.length === 0"
            icon="i-mdi-ticket"
            :message="$t('common.not_found', [$t('common.activity')])"
        />

        <ul v-else role="list" class="mb-1 divide-y divide-gray-100 dark:divide-gray-800">
            <PageActivityListEntry v-for="item in data.activity" :key="item.id" :entry="item" />
        </ul>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
