<script lang="ts" setup>
import { useSettingsStore, type DispatchCenterInnerPane, type DispatchCenterOuterPane } from '~/stores/settings';

const { t } = useI18n();

const props = withDefaults(
    defineProps<{
        hideInnerPanes?: boolean;
    }>(),
    {
        hideInnerPanes: false,
    },
);

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const defaultDispatchCenterPaneLayout = {
    outer: ['map', 'sidebar'] as DispatchCenterOuterPane[],
    inner: ['dispatchList', 'unitList', 'feed'] as DispatchCenterInnerPane[],
};

const open = ref(false);

const layoutDraft = reactive<{
    outer: DispatchCenterOuterPane[];
    inner: DispatchCenterInnerPane[];
}>({
    outer: [],
    inner: [],
});

function syncLayoutDraft(): void {
    layoutDraft.outer = [...centrum.value.dispatchCenterPaneLayout.outer];
    if (!props.hideInnerPanes) {
        layoutDraft.inner = [...centrum.value.dispatchCenterPaneLayout.inner];
    }
}

watch(
    open,
    (isOpen) => {
        if (isOpen) syncLayoutDraft();
    },
    { flush: 'sync' },
);

function paneLabel(pane: DispatchCenterOuterPane | DispatchCenterInnerPane): string {
    switch (pane) {
        case 'map':
            return t('components.dispatch.layout_popover.map');
        case 'sidebar':
            return t('components.dispatch.layout_popover.sidebar');
        case 'dispatchList':
            return t('components.dispatch.layout_popover.dispatches');
        case 'unitList':
            return t('components.dispatch.layout_popover.units');
        case 'feed':
            return t('components.dispatch.layout_popover.activity');
    }
}

function moveLayoutItem<T>(items: T[], index: number, delta: -1 | 1): T[] {
    const target = index + delta;
    if (target < 0 || target >= items.length) return items;

    const next = [...items];
    [next[index], next[target]] = [next[target]!, next[index]!];
    return next;
}

function moveOuterPane(index: number, delta: -1 | 1): void {
    layoutDraft.outer = moveLayoutItem(layoutDraft.outer, index, delta);
}

function moveInnerPane(index: number, delta: -1 | 1): void {
    layoutDraft.inner = moveLayoutItem(layoutDraft.inner, index, delta);
}

function saveLayoutDraft(): void {
    centrum.value.dispatchCenterPaneLayout.outer = [...layoutDraft.outer];
    if (!props.hideInnerPanes) {
        centrum.value.dispatchCenterPaneLayout.inner = [...layoutDraft.inner];
    }
    open.value = false;
}

function resetLayout(): void {
    layoutDraft.outer = [...defaultDispatchCenterPaneLayout.outer];
    if (!props.hideInnerPanes) {
        layoutDraft.inner = [...defaultDispatchCenterPaneLayout.inner];
    }
    saveLayoutDraft();
}

const contentClass = computed(() => (!props.hideInnerPanes ? 'w-[28rem]' : 'w-[24rem]'));
</script>

<template>
    <UPopover v-model:open="open" :content="{ align: 'end', sideOffset: 8 }">
        <UTooltip :text="$t('components.dispatch.layout_popover.trigger')">
            <UButton
                color="neutral"
                variant="ghost"
                icon="i-mdi-view-split-vertical"
                :aria-label="$t('components.dispatch.layout_popover.trigger')"
            />
        </UTooltip>

        <template #content>
            <div class="max-w-[calc(100vw-2rem)] p-4" :class="contentClass">
                <div class="space-y-6">
                    <div class="space-y-3">
                        <div>
                            <h3 class="text-sm font-semibold text-highlighted">
                                {{ $t('components.dispatch.layout_popover.map_sidebar_title') }}
                            </h3>
                            <p class="text-sm text-muted">
                                {{ $t('components.dispatch.layout_popover.map_sidebar_description') }}
                            </p>
                        </div>

                        <div class="grid gap-2 sm:grid-cols-2">
                            <div
                                v-for="(pane, index) in layoutDraft.outer"
                                :key="pane"
                                class="flex items-center justify-between gap-3 rounded-lg border border-default bg-elevated/50 px-3 py-2"
                            >
                                <div class="flex flex-col">
                                    <span class="font-medium text-highlighted">{{ paneLabel(pane) }}</span>
                                    <span class="text-xs text-muted">
                                        {{
                                            index === 0
                                                ? $t('components.dispatch.layout_popover.left')
                                                : $t('components.dispatch.layout_popover.right')
                                        }}
                                    </span>
                                </div>

                                <div class="flex items-center gap-1">
                                    <UTooltip :text="$t('components.dispatch.layout_popover.move_left', [paneLabel(pane)])">
                                        <UButton
                                            color="neutral"
                                            variant="ghost"
                                            icon="i-mdi-arrow-left"
                                            :disabled="index === 0"
                                            :aria-label="$t('components.dispatch.layout_popover.move_left', [paneLabel(pane)])"
                                            @click="moveOuterPane(index, -1)"
                                        />
                                    </UTooltip>

                                    <UTooltip :text="$t('components.dispatch.layout_popover.move_right', [paneLabel(pane)])">
                                        <UButton
                                            color="neutral"
                                            variant="ghost"
                                            icon="i-mdi-arrow-right"
                                            :disabled="index === layoutDraft.outer.length - 1"
                                            :aria-label="$t('components.dispatch.layout_popover.move_right', [paneLabel(pane)])"
                                            @click="moveOuterPane(index, 1)"
                                        />
                                    </UTooltip>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div v-if="!props.hideInnerPanes" class="space-y-3">
                        <div>
                            <h3 class="text-sm font-semibold text-highlighted">
                                {{ $t('components.dispatch.layout_popover.inner_title') }}
                            </h3>
                            <p class="text-sm text-muted">
                                {{ $t('components.dispatch.layout_popover.inner_description') }}
                            </p>
                        </div>

                        <div class="space-y-2">
                            <div
                                v-for="(pane, index) in layoutDraft.inner"
                                :key="pane"
                                class="flex items-center justify-between gap-3 rounded-lg border border-default bg-elevated/50 px-3 py-2"
                            >
                                <div class="flex flex-col">
                                    <span class="font-medium text-highlighted">{{ paneLabel(pane) }}</span>
                                    <span class="text-xs text-muted">
                                        {{ $t('components.dispatch.layout_popover.position', [index + 1]) }}
                                    </span>
                                </div>

                                <div class="flex items-center gap-1">
                                    <UTooltip :text="$t('components.dispatch.layout_popover.move_up', [paneLabel(pane)])">
                                        <UButton
                                            color="neutral"
                                            variant="ghost"
                                            icon="i-mdi-arrow-up"
                                            :disabled="index === 0"
                                            :aria-label="$t('components.dispatch.layout_popover.move_up', [paneLabel(pane)])"
                                            @click="moveInnerPane(index, -1)"
                                        />
                                    </UTooltip>

                                    <UTooltip :text="$t('components.dispatch.layout_popover.move_down', [paneLabel(pane)])">
                                        <UButton
                                            color="neutral"
                                            variant="ghost"
                                            icon="i-mdi-arrow-down"
                                            :disabled="index === layoutDraft.inner.length - 1"
                                            :aria-label="$t('components.dispatch.layout_popover.move_down', [paneLabel(pane)])"
                                            @click="moveInnerPane(index, 1)"
                                        />
                                    </UTooltip>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-4 flex items-center gap-2">
                    <UTooltip class="self-start" :text="$t('common.reset')">
                        <UButton color="neutral" icon="i-mdi-clear-box-outline" variant="soft" @click="resetLayout" />
                    </UTooltip>

                    <div class="flex-1" />

                    <UTooltip :text="$t('common.cancel')">
                        <UButton color="neutral" variant="soft" :label="$t('common.cancel')" @click="open = false" />
                    </UTooltip>

                    <UTooltip :text="$t('common.save')">
                        <UButton :label="$t('common.save')" @click="saveLayoutDraft" />
                    </UTooltip>
                </div>
            </div>
        </template>
    </UPopover>
</template>
