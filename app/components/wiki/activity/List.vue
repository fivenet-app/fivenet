<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import ListEntry from '~/components/wiki/activity/ListEntry.vue';
import { getPageActivityIcon, getPageActivityIndicatorColor } from '~/components/wiki/helpers';
import type { ListPageActivityResponse } from '~~/gen/ts/services/wiki/wiki';

const props = defineProps<{
    pageId: number;
}>();

const { listPageActivity: listWikiPageActivity } = await useWikiWiki();

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, refresh, error } = useAuthedLazyAsyncData(
    'userState',
    () => `wiki-page:${props.pageId}:activity:${page.value}`,
    ({ signal }) => listPageActivity(signal),
);

const timelineItems = computed(() =>
    (data.value?.activity ?? []).map((activity) => ({
        ...activity,
        icon: getPageActivityIcon(activity.activityType),
        ui: { indicator: getPageActivityIndicatorColor(activity.activityType) },
    })),
);

async function listPageActivity(signal: AbortSignal): Promise<ListPageActivityResponse> {
    return listWikiPageActivity(
        {
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            pageId: props.pageId,
        },
        { abort: signal },
    );
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.activity', 2)])" />
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

        <UTimeline v-else class="mb-1" :items="timelineItems" size="md" :ui="{ wrapper: '!mt-0 !pb-2' }">
            <template #indicator="{ item }">
                <UIcon :name="item.icon" class="size-6" />
            </template>

            <template #wrapper="{ item }">
                <ListEntry :entry="item" />
            </template>
        </UTimeline>

        <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
    </div>
</template>
