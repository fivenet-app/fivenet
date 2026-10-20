<script lang="ts" setup>
import {
    defaultDispatchCenterInnerPaneLayout,
    defaultDispatchCenterOuterPaneLayout,
    type DispatchCenterInnerPane,
    type DispatchCenterOuterPane,
    useSettingsStore,
} from '~/stores/settings';

const { t } = useI18n();

const props = withDefaults(
    defineProps<{
        hideInnerPanes?: boolean;
        autoSave?: boolean;
    }>(),
    {
        hideInnerPanes: false,
        autoSave: false,
    },
);

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

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

syncLayoutDraft();

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

function saveLayoutDraft(): void {
    centrum.value.dispatchCenterPaneLayout.outer = [...layoutDraft.outer];
    if (!props.hideInnerPanes) {
        centrum.value.dispatchCenterPaneLayout.inner = [...layoutDraft.inner];
    }
}

function moveOuterLayoutItem(index: number, delta: -1 | 1): void {
    layoutDraft.outer = moveLayoutItem(layoutDraft.outer, index, delta);
    if (props.autoSave) saveLayoutDraft();
}

function moveInnerLayoutItem(index: number, delta: -1 | 1): void {
    layoutDraft.inner = moveLayoutItem(layoutDraft.inner, index, delta);
    if (props.autoSave) saveLayoutDraft();
}

function resetLayout(): void {
    layoutDraft.outer = [...defaultDispatchCenterOuterPaneLayout];
    if (!props.hideInnerPanes) {
        layoutDraft.inner = [...defaultDispatchCenterInnerPaneLayout];
    }
    saveLayoutDraft();
}
</script>

<template>
    <div class="space-y-4">
        <UFormField
            class="flex items-center gap-4"
            name="dispatchListCardStyle"
            :label="$t('common.card_view')"
            :description="$t('components.auth.user_settings.dispatch_list_card_style.description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="centrum.dispatchListCardStyle" />
        </UFormField>

        <div class="space-y-4">
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
                        <UButton
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-arrow-left"
                            :disabled="index === 0"
                            :aria-label="$t('components.dispatch.layout_popover.move_left', [paneLabel(pane)])"
                            @click="moveOuterLayoutItem(index, -1)"
                        />
                        <UButton
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-arrow-right"
                            :disabled="index === layoutDraft.outer.length - 1"
                            :aria-label="$t('components.dispatch.layout_popover.move_right', [paneLabel(pane)])"
                            @click="moveOuterLayoutItem(index, 1)"
                        />
                    </div>
                </div>
            </div>
        </div>

        <div v-if="!props.hideInnerPanes" class="space-y-4">
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
                        <span class="text-xs text-muted">{{
                            $t('components.dispatch.layout_popover.position', [index + 1])
                        }}</span>
                    </div>

                    <div class="flex items-center gap-1">
                        <UButton
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-arrow-up"
                            :disabled="index === 0"
                            :aria-label="$t('components.dispatch.layout_popover.move_up', [paneLabel(pane)])"
                            @click="moveInnerLayoutItem(index, -1)"
                        />
                        <UButton
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-arrow-down"
                            :disabled="index === layoutDraft.inner.length - 1"
                            :aria-label="$t('components.dispatch.layout_popover.move_down', [paneLabel(pane)])"
                            @click="moveInnerLayoutItem(index, 1)"
                        />
                    </div>
                </div>
            </div>
        </div>

        <div class="flex items-center gap-2">
            <UButton color="neutral" variant="soft" :label="$t('common.reset')" @click="resetLayout" />
            <div class="flex-1" />
            <UButton v-if="!props.autoSave" :label="$t('common.save')" @click="saveLayoutDraft" />
        </div>
    </div>
</template>
