<script lang="ts" setup>
import { useDraggable } from 'vue-draggable-plus';
import { storeToRefs } from 'pinia';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';
import type { OverviewFeature } from '~/composables/useOverviewFeatures';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';

const settingsStore = useSettingsStore();
const { overviewQuickAccess } = storeToRefs(settingsStore);
const { reorderOverviewQuickAccess } = settingsStore;

const items = useOverviewFeatures();

const quickAccessItems = computed(() => {
    const seen = new Set<string>();
    return overviewQuickAccess.value
        .map((path) => items.value.find((item) => item.to === path))
        .filter((item): item is NonNullable<(typeof items.value)[number]> => item !== undefined)
        .filter((item) => {
            if (seen.has(item.to)) return false;
            seen.add(item.to);
            return true;
        });
});

const reorderMode = ref<boolean>(false);
const sortableQuickAccessItems = ref<OverviewFeature[]>([]);
const listRef = useTemplateRef('listRef');
const activeDraggableRoot = shallowRef<HTMLElement | null>(null);

watch(quickAccessItems, (value) => (sortableQuickAccessItems.value = [...value]), { immediate: true });

const draggable = useDraggable(listRef, sortableQuickAccessItems, {
    immediate: false,
    animation: 150,
    handle: '.handle',
    draggable: '> .quick-access-item',
    onEnd: () => reorderOverviewQuickAccess(sortableQuickAccessItems.value.map((item) => item.to)),
});

const { moveUp, moveDown } = useListReorder(sortableQuickAccessItems, {
    onMove: () => reorderOverviewQuickAccess(sortableQuickAccessItems.value.map((item) => item.to)),
});

watch(
    [reorderMode, listRef],
    ([enabled, root]) => {
        if (!enabled || !root) {
            draggable.pause();
            activeDraggableRoot.value = null;
            return;
        }

        if (activeDraggableRoot.value !== root) {
            draggable.start(root);
            activeDraggableRoot.value = root;
        }

        draggable.resume();
    },
    { immediate: true, flush: 'post' },
);
</script>

<template>
    <div v-if="quickAccessItems.length > 0">
        <div class="mb-3 flex items-center justify-between gap-2">
            <div class="flex flex-col">
                <h2 class="text-base font-semibold text-highlighted">
                    {{ $t('common.quick_access') }}
                </h2>
            </div>

            <UButton
                :icon="reorderMode ? 'i-mdi-check' : 'i-mdi-sort'"
                color="gray"
                variant="ghost"
                :label="$t('common.change_order')"
                :disabled="quickAccessItems.length < 2"
                @click="() => (reorderMode = !reorderMode)"
            />
        </div>

        <div class="rounded-2xl border border-default/60 bg-elevated/70 p-2 shadow-sm backdrop-blur-sm">
            <div ref="listRef" class="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
                <ULink
                    v-for="(item, idx) in sortableQuickAccessItems"
                    :key="item.to"
                    :to="item.to"
                    class="quick-access-item group relative flex min-h-20 flex-col items-center justify-center gap-2 rounded-2xl border border-transparent bg-default/70 px-2 py-3 text-center transition hover:-translate-y-0.5 hover:border-primary-500/30 hover:bg-primary-50/60 dark:bg-default/30 dark:hover:bg-primary-950/20"
                >
                    <template v-if="reorderMode">
                        <DraggableHandle class="absolute top-1 left-1" size="xs" orientation="vertical" />
                        <div class="absolute top-1 right-1">
                            <ReorderButtons :idx="idx" :move-up="moveUp" :move-down="moveDown" orientation="horizontal" />
                        </div>
                    </template>

                    <UIcon class="size-8 text-primary transition group-hover:scale-110" :name="item.icon ?? 'i-mdi-star'" />
                    <span class="text-sm leading-tight font-medium text-highlighted">
                        {{ item.title }}
                    </span>
                </ULink>
            </div>
        </div>
    </div>
</template>
