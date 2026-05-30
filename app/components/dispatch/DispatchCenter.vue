<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import DispatchList from '~/components/dispatch/dispatches/DispatchList.vue';
import Feed from '~/components/dispatch/Feed.vue';
import DispatchLayer from '~/components/dispatch/livemap/DispatchLayer.vue';
import UnitList from '~/components/dispatch/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useSettingsStore } from '~/stores/settings';
import DispatcherInfo from './dispatchers/DispatcherInfo.vue';

const { can } = useAuth();

const centrumStore = useCentrumStore();
const { error, feed, isCenter } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

type SplitpanesPane = {
    size: number;
};

type SplitpanesResizedEvent = {
    panes?: SplitpanesPane[];
};

const roundPaneSize = (size: number): number => Math.round(size * 100) / 100;
const defaultDispatchCenterPaneSizes = {
    map: 30,
    sidebar: 70,
    dispatchList: 58,
    unitList: 26,
    feed: 8,
};

centrum.value.dispatchCenterPaneSizes = {
    map: centrum.value.dispatchCenterPaneSizes?.map ?? defaultDispatchCenterPaneSizes.map,
    sidebar: centrum.value.dispatchCenterPaneSizes?.sidebar ?? defaultDispatchCenterPaneSizes.sidebar,
    dispatchList: centrum.value.dispatchCenterPaneSizes?.dispatchList ?? defaultDispatchCenterPaneSizes.dispatchList,
    unitList: centrum.value.dispatchCenterPaneSizes?.unitList ?? defaultDispatchCenterPaneSizes.unitList,
    feed: centrum.value.dispatchCenterPaneSizes?.feed ?? defaultDispatchCenterPaneSizes.feed,
};

function onOuterPanesResized({ panes }: SplitpanesResizedEvent): void {
    if (!panes || panes.length < 2) return;

    centrum.value.dispatchCenterPaneSizes.map = roundPaneSize(panes[0]!.size);
    centrum.value.dispatchCenterPaneSizes.sidebar = roundPaneSize(panes[1]!.size);
}

function onInnerPanesResized({ panes }: SplitpanesResizedEvent): void {
    if (!panes || panes.length < 3) return;

    centrum.value.dispatchCenterPaneSizes.dispatchList = roundPaneSize(panes[0]!.size);
    centrum.value.dispatchCenterPaneSizes.unitList = roundPaneSize(panes[1]!.size);
    centrum.value.dispatchCenterPaneSizes.feed = roundPaneSize(panes[2]!.size);
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
                    <ClientOnly>
                        <DispatcherInfo />
                    </ClientOnly>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div
                class="max-h-[calc(100dvh-var(--ui-header-height))] min-h-[calc(100dvh-var(--ui-header-height))] w-full overflow-hidden"
            >
                <Splitpanes @resized="onOuterPanesResized">
                    <Pane :min-size="25" :size="centrum.dispatchCenterPaneSizes.map">
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
                    </Pane>

                    <Pane :min-size="40" :size="centrum.dispatchCenterPaneSizes.sidebar">
                        <Splitpanes horizontal @resized="onInnerPanesResized">
                            <Pane :size="centrum.dispatchCenterPaneSizes.dispatchList" :min-size="2">
                                <DispatchList show-button />
                            </Pane>

                            <Pane :size="centrum.dispatchCenterPaneSizes.unitList" :min-size="2">
                                <UnitList />
                            </Pane>

                            <Pane :size="centrum.dispatchCenterPaneSizes.feed" :min-size="2">
                                <Feed :items="feed" />
                            </Pane>
                        </Splitpanes>
                    </Pane>
                </Splitpanes>
            </div>
        </template>
    </UDashboardPanel>
</template>

<style>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 2px;
    background-color: var(--ui-border);
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 2px;
    background-color: var(--ui-border);
}
</style>
