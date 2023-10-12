<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import { default as DispatchesList } from '~/components/centrum/dispatches/List.vue';
import DisponentsInfo from '~/components/centrum/disponents/DisponentsInfo.vue';
import { default as UnitsList } from '~/components/centrum/units/List.vue';
import Livemap from '~/components/livemap/Livemap.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import Feed from './Feed.vue';
import DispatchesLayer from './livemap/DispatchesLayer.vue';

const centrumStore = useCentrumStore();
const { error, abort, restarting, feed } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

onBeforeMount(async () => setTimeout(async () => startStream(true), 250));

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

<template>
    <div class="relative h-full w-full">
        <div
            v-if="error !== undefined || (abort === undefined && !restarting)"
            class="absolute inset-0 flex justify-center items-center z-30 bg-gray-600/70"
        >
            <DataErrorBlock
                v-if="error"
                :title="$t('components.centrum.dispatch_center.failed_datastream')"
                :retry="startStream"
            />
            <DataPendingBlock
                v-else-if="abort === undefined && !restarting"
                :message="$t('components.centrum.dispatch_center.starting_datastream')"
            />
        </div>
        <Splitpanes v-else class="h-full w-full">
            <Pane min-size="25">
                <Splitpanes horizontal>
                    <Pane min-size="80">
                        <Livemap :show-unit-names="true">
                            <template v-slot:default>
                                <DispatchesLayer
                                    v-if="can('CentrumService.Stream')"
                                    :show-all-dispatches="true"
                                    @goto="goto($event)"
                                />
                            </template>
                        </Livemap>
                    </Pane>
                    <Pane min-size="8" size="8">
                        <DisponentsInfo />
                    </Pane>
                </Splitpanes>
            </Pane>
            <Pane size="65">
                <Splitpanes horizontal>
                    <Pane size="55" min-size="2">
                        <DispatchesList @goto="goto($event)" />
                    </Pane>
                    <Pane size="33" min-size="2">
                        <UnitsList @goto="goto($event)" />
                    </Pane>
                    <Pane size="12" min-size="2">
                        <Feed :items="feed" />
                    </Pane>
                </Splitpanes>
            </Pane>
        </Splitpanes>
    </div>
</template>
