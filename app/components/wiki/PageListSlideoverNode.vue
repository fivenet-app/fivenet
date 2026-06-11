<script lang="ts" setup>
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import type { PageShort } from '~~/gen/ts/resources/wiki/page';
import { pageToURL, sameWikiMoveGroup } from './helpers';

const props = defineProps<{
    page: PageShort;
    siblings: PageShort[];
    index: number;
    depth?: number;
    movingPageId?: number;
}>();

const emit = defineEmits<{
    (
        e: 'move',
        v: {
            pageId: number;
            beforeId?: number;
            afterId?: number;
        },
    ): void;
}>();

const isBusy = computed(() => props.movingPageId === props.page.id);
const route = useRoute('wiki-job-id-slug');
const currentPageId = computed(() => Number(route.params.id));
const isCurrentPage = computed(() => props.page.id === currentPageId.value);

const canMoveUp = computed(
    () => props.index > 0 && !isBusy.value && sameWikiMoveGroup(props.page, props.siblings[props.index - 1]),
);
const canMoveDown = computed(
    () =>
        props.index < props.siblings.length - 1 &&
        !isBusy.value &&
        sameWikiMoveGroup(props.page, props.siblings[props.index + 1]),
);

const pagePath = computed(() => pageToURL(props.page));
const cardClass = computed(() =>
    isCurrentPage.value
        ? 'border-primary-300 bg-primary-50/80 ring-1 ring-inset ring-primary-200 dark:border-primary-800 dark:bg-primary-950/20 dark:ring-primary-900/60'
        : 'border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900',
);

function moveUp(idx: number): void {
    const beforeId = props.siblings[idx - 1]?.id;
    if (!beforeId || !canMoveUp.value) return;

    emit('move', {
        pageId: props.page.id,
        beforeId: beforeId,
    });
}

function moveDown(idx: number): void {
    const afterId = props.siblings[idx + 1]?.id;
    if (!afterId || !canMoveDown.value) return;

    emit('move', {
        pageId: props.page.id,
        afterId: afterId,
    });
}
</script>

<template>
    <div class="space-y-2">
        <UCard :class="cardClass" :ui="{ body: 'p-2 sm:p-2' }">
            <div class="flex items-start gap-3">
                <div class="min-w-0 flex-1">
                    <div class="flex flex-wrap items-center gap-2">
                        <UButton
                            class="!p-0 text-left"
                            variant="link"
                            :color="isCurrentPage ? 'primary' : 'neutral'"
                            :to="pagePath"
                            :label="page.title || $t('common.untitled')"
                        />

                        <UBadge
                            v-if="page.startpage"
                            color="neutral"
                            variant="soft"
                            icon="i-mdi-home"
                            :label="$t('common.startpage')"
                        />
                        <UBadge
                            v-if="page.draft"
                            color="warning"
                            variant="soft"
                            icon="i-mdi-pencil"
                            :label="$t('common.draft')"
                        />
                        <UBadge
                            v-if="page.deletedAt"
                            color="error"
                            variant="soft"
                            icon="i-mdi-delete"
                            :label="$t('common.deleted')"
                        />
                    </div>

                    <p v-if="page.description" class="mt-1 line-clamp-2 text-sm text-neutral-500 dark:text-neutral-400">
                        {{ page.description }}
                    </p>
                </div>

                <div class="flex shrink-0 items-center gap-1">
                    <ReorderButtons
                        :idx="index"
                        :move-up="moveUp"
                        :move-down="moveDown"
                        :disable-up="!canMoveUp"
                        :disable-down="!canMoveDown"
                        orientation="horizontal"
                        direction="vertical"
                    />
                </div>
            </div>
        </UCard>

        <div
            v-if="page.children.length > 0"
            class="space-y-2 border-l border-dashed border-neutral-200 pl-4 dark:border-neutral-800"
        >
            <PageListSlideoverNode
                v-for="(child, childIndex) in page.children"
                :key="child.id"
                :page="child"
                :siblings="page.children"
                :index="childIndex"
                :depth="(depth ?? 0) + 1"
                :moving-page-id="movingPageId"
                @move="emit('move', $event)"
            />
        </div>
    </div>
</template>
