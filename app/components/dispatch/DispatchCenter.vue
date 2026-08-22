<script lang="ts" setup>
import type { SplitterItem } from '@nuxt/ui';
import DispatchList from '~/components/dispatch/dispatches/DispatchList.vue';
import Feed from '~/components/dispatch/Feed.vue';
import DispatchLayer from '~/components/dispatch/livemap/DispatchLayer.vue';
import UnitList from '~/components/dispatch/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useSettingsStore, type DispatchCenterInnerPane, type DispatchCenterOuterPane } from '~/stores/settings';
import DispatchCenterLayoutPopover from './DispatchCenterLayoutPopover.vue';
import DispatcherInfo from './dispatchers/DispatcherInfo.vue';

const { can } = useAuth();

const centrumStore = useCentrumStore();
const { error, feed, isCenter } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const roundPaneSize = (size: number): number => Math.round(size * 100) / 100;
const defaultDispatchCenterPaneSizes = {
    map: 30,
    sidebar: 70,
    dispatchList: 58,
    unitList: 26,
    feed: 8,
};
const defaultDispatchCenterPaneLayout = {
    outer: ['map', 'sidebar'] as DispatchCenterOuterPane[],
    inner: ['dispatchList', 'unitList', 'feed'] as DispatchCenterInnerPane[],
};

centrum.value.dispatchCenterPaneSizes = {
    map: centrum.value.dispatchCenterPaneSizes?.map ?? defaultDispatchCenterPaneSizes.map,
    sidebar: centrum.value.dispatchCenterPaneSizes?.sidebar ?? defaultDispatchCenterPaneSizes.sidebar,
    dispatchList: centrum.value.dispatchCenterPaneSizes?.dispatchList ?? defaultDispatchCenterPaneSizes.dispatchList,
    unitList: centrum.value.dispatchCenterPaneSizes?.unitList ?? defaultDispatchCenterPaneSizes.unitList,
    feed: centrum.value.dispatchCenterPaneSizes?.feed ?? defaultDispatchCenterPaneSizes.feed,
};
centrum.value.dispatchCenterPaneLayout = {
    outer:
        centrum.value.dispatchCenterPaneLayout?.outer?.length === 2
            ? [...centrum.value.dispatchCenterPaneLayout.outer]
            : [...defaultDispatchCenterPaneLayout.outer],
    inner:
        centrum.value.dispatchCenterPaneLayout?.inner?.length === 3
            ? [...centrum.value.dispatchCenterPaneLayout.inner]
            : [...defaultDispatchCenterPaneLayout.inner],
};

const splitterUi = {
    panel: 'min-h-0 min-w-0 overflow-hidden',
    handle: 'data-[orientation=horizontal]:w-px data-[orientation=vertical]:h-px bg-border transition-colors data-[state=hover]:bg-primary data-[state=drag]:bg-primary',
};

const outerPaneItems = {
    map: {
        slot: 'map',
        minSize: 25,
    },
    sidebar: {
        slot: 'sidebar',
        minSize: 40,
    },
} satisfies Record<DispatchCenterOuterPane, Pick<SplitterItem, 'slot' | 'minSize'>>;

const innerPaneItems = {
    dispatchList: {
        slot: 'dispatchList',
        minSize: 2,
    },
    unitList: {
        slot: 'unitList',
        minSize: 2,
    },
    feed: {
        slot: 'feed',
        minSize: 2,
    },
} satisfies Record<DispatchCenterInnerPane, Pick<SplitterItem, 'slot' | 'minSize'>>;

const outerItems = computed<SplitterItem[]>(() =>
    centrum.value.dispatchCenterPaneLayout.outer.map((pane) => ({
        ...outerPaneItems[pane],
        defaultSize: centrum.value.dispatchCenterPaneSizes[pane],
    })),
);

const innerItems = computed<SplitterItem[]>(() =>
    centrum.value.dispatchCenterPaneLayout.inner.map((pane) => ({
        ...innerPaneItems[pane],
        defaultSize: centrum.value.dispatchCenterPaneSizes[pane],
    })),
);

const outerSplitterKey = computed(() => `dispatch-center-splitter-${centrum.value.dispatchCenterPaneLayout.outer.join('-')}`);
const innerSplitterKey = computed(
    () => `dispatch-center-splitter-inner-${centrum.value.dispatchCenterPaneLayout.inner.join('-')}`,
);

function onOuterLayout(sizes: number[]): void {
    const panes = centrum.value.dispatchCenterPaneLayout.outer;
    if (sizes.length < panes.length) return;

    panes.forEach((pane, index) => {
        const size = sizes[index];
        if (size === undefined) return;

        centrum.value.dispatchCenterPaneSizes[pane] = roundPaneSize(size);
    });
}

function onInnerLayout(sizes: number[]): void {
    const panes = centrum.value.dispatchCenterPaneLayout.inner;
    if (sizes.length < panes.length) return;

    panes.forEach((pane, index) => {
        const size = sizes[index];
        if (size === undefined) return;

        centrum.value.dispatchCenterPaneSizes[pane] = roundPaneSize(size);
    });
}

onBeforeMount(async () => {
    isCenter.value = true;
    useTimeoutFn(async () => {
        try {
            startStream();
        } catch (e) {
            logger.error('exception during start of centrum stream', e);
        }
    }, 500);
});

onBeforeRouteLeave(async (to) => {
    isCenter.value = false;

    // Don't end centrum stream if user is switching to livemap page
    if (to.path.startsWith('/livemap')) return;

    await stopStream(true);
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.dispatch_center')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <div class="flex items-center gap-2">
                        <DispatchCenterLayoutPopover />

                        <ClientOnly>
                            <DispatcherInfo />
                        </ClientOnly>
                    </div>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div
                class="max-h-[calc(100dvh-var(--ui-header-height)-var(--page-content-bottom-offset))] min-h-[calc(100dvh-var(--ui-header-height)-var(--page-content-bottom-offset))] w-full overflow-hidden"
            >
                <USplitter
                    id="dispatch-center-splitter"
                    :key="outerSplitterKey"
                    class="size-full"
                    :items="outerItems"
                    :ui="splitterUi"
                    @layout="onOuterLayout"
                >
                    <template #map>
                        <div class="relative size-full">
                            <div v-if="error" class="absolute inset-0 z-30 flex items-center justify-center bg-default/75">
                                <DataErrorBlock
                                    :title="$t('components.dispatch.dispatch_center.failed_datastream')"
                                    :error="error"
                                    :retry="startStream"
                                />
                            </div>

                            <LivemapBase show-unit-names show-unit-status>
                                <template #default>
                                    <DispatchLayer v-if="can('centrum.CentrumService/Stream').value" show-all-dispatches />
                                </template>
                            </LivemapBase>
                        </div>
                    </template>

                    <template #sidebar>
                        <USplitter
                            id="dispatch-center-splitter-inner"
                            :key="innerSplitterKey"
                            orientation="vertical"
                            class="size-full"
                            :items="innerItems"
                            :ui="splitterUi"
                            @layout="onInnerLayout"
                        >
                            <template #dispatchList>
                                <DispatchList show-button />
                            </template>

                            <template #unitList>
                                <UnitList />
                            </template>

                            <template #feed>
                                <Feed :items="feed" />
                            </template>
                        </USplitter>
                    </template>
                </USplitter>
            </div>
        </template>
    </UDashboardPanel>
</template>
