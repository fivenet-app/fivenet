<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import PageListSlideoverNode from './PageListSlideoverNode.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';
import RefreshButton from '~/components/partials/RefreshButton.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const props = defineProps<{
    job: string;
    refresh: () => Promise<void>;
}>();

const notifications = useNotificationsStore();
const { listPages: listWikiPages, movePage: moveWikiPage } = await useWikiWiki();

const movingPageId = ref<number | undefined>(undefined);
const wikiPageChunkSize = 250;

const { data, status, refresh, error } = useLazyAsyncData(`wiki-pages-move-${props.job}`, () => listPages());

const pages = computed(() => data.value ?? []);

async function listPages(): Promise<PageShort[]> {
    const allPages: PageShort[] = [];
    let offset = 0;
    let totalCount = Number.POSITIVE_INFINITY;

    while (offset < totalCount) {
        const response = await listWikiPages({
            pagination: {
                offset: offset,
                pageSize: wikiPageChunkSize,
            },
            job: props.job,
            rootOnly: false,
        });

        allPages.push(...response.pages);
        totalCount = response.pagination?.totalCount ?? allPages.length;

        if (response.pages.length < wikiPageChunkSize) {
            break;
        }

        offset += wikiPageChunkSize;
    }

    return allPages;
}

async function movePage(payload: { pageId: number; beforeId?: number; afterId?: number }): Promise<void> {
    movingPageId.value = payload.pageId;

    try {
        await moveWikiPage(payload);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await refresh();
        void props.refresh().catch(() => undefined);
    } finally {
        movingPageId.value = undefined;
    }
}
</script>

<template>
    <USlideover :title="$t('common.change_order')" :overlay="false" :ui="{ content: 'max-w-4xl' }">
        <template #actions>
            <RefreshButton :loading="isRequestPending(status)" @click="() => refresh()" />
        </template>

        <template #body>
            <div class="flex h-full w-full flex-1 flex-col gap-4">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.page', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.page', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="pages.length === 0" :type="$t('common.page', 2)" icon="i-mdi-file-tree" />

                <template v-else>
                    <div class="space-y-3">
                        <PageListSlideoverNode
                            v-for="(page, index) in pages"
                            :key="page.id"
                            :page="page"
                            :siblings="pages"
                            :index="index"
                            :moving-page-id="movingPageId"
                            @move="movePage"
                        />
                    </div>
                </template>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </USlideover>
</template>
