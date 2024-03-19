<script lang="ts" setup>
import { useTimeoutFn } from '@vueuse/core';
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import DisponentsInfo from '~/components/centrum/disponents/DisponentsInfo.vue';
import UnitList from '~/components/centrum/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import CentrumFeed from '~/components/centrum/CentrumFeed.vue';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import MarkersList from '~/components/centrum/MarkersList.vue';

const centrumStore = useCentrumStore();
const { error, abort, reconnecting, feed } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

onBeforeMount(async () => useTimeoutFn(async () => startStream(true), 250));

onBeforeUnmount(async () => {
    stopStream();
    centrumStore.$reset();
});

function goto(e: Coordinate) {
    location.value = e;

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}
</script>

<template>
    <div class="relative h-full w-full">
        <div
            v-if="error !== undefined || (abort === undefined && !reconnecting)"
            class="absolute inset-0 z-30 flex items-center justify-center bg-gray-600/70"
        >
            <DataErrorBlock
                v-if="error"
                :title="$t('components.centrum.dispatch_center.failed_datastream')"
                :retry="startStream"
            />
            <DataPendingBlock
                v-else-if="abort === undefined && !reconnecting"
                :message="$t('components.centrum.dispatch_center.starting_datastream')"
            />
        </div>
        <Splitpanes v-else class="h-full w-full">
            <Pane min-size="25">
                <Splitpanes horizontal>
                    <Pane min-size="80" size="90">
                        <LivemapBase :show-unit-names="true" :show-unit-status="true" @goto="goto($event)">
                            <template #default>
                                <DispatchesLayer
                                    v-if="can('CentrumService.Stream')"
                                    :show-all-dispatches="true"
                                    @goto="goto($event)"
                                />
                            </template>
                        </LivemapBase>
                    </Pane>
                    <Pane min-size="8" size="10">
                        <DisponentsInfo />
                    </Pane>
                </Splitpanes>
            </Pane>
            <Pane size="65">
                <Splitpanes horizontal>
                    <Pane size="58" min-size="2">
                        <DispatchList :show-button="true" @goto="goto($event)" />
                    </Pane>
                    <Pane size="26" min-size="2">
                        <UnitList @goto="goto($event)" />
                    </Pane>
                    <Pane size="8" min-size="2">
                        <MarkersList @goto="goto($event)" />
                    </Pane>
                    <Pane size="8" min-size="2">
                        <CentrumFeed :items="feed" @goto="goto($event)" />
                    </Pane>
                </Splitpanes>
            </Pane>
        </Splitpanes>
    </div>
</template>

<style>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 3px;
    background-color: #898d9b;
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 3px;
    background-color: #898d9b;
}
</style>
