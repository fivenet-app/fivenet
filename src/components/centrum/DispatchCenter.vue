<script lang="ts" setup>
import { default as DispatchesList } from '~/components/centrum/dispatches/List.vue';
import { default as UnitsList } from '~/components/centrum/units/List.vue';
import Livemap from '~/components/livemap/Livemap.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { setWaypoint } from '~/composables/nui';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import DisponentsInfo from './DisponentsInfo.vue';
import Feed from './Feed.vue';
import DispatchesLayer from './livemap/DispatchesLayer.vue';

const centrumStore = useCentrumStore();
const { error, abort, restarting, isDisponent, disponents, settings, feed } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

onBeforeMount(async () => setTimeout(async () => startStream(true), 250));

onBeforeUnmount(async () => {
    stopStream();
    centrumStore.$reset();
});

function goto(e: Coordinate) {
    location.value = { x: e.x, y: e.y };

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}
</script>

<template>
    <div class="flex-col h-full relative">
        <div
            v-if="(error || abort === undefined) && !restarting"
            class="absolute inset-0 flex justify-center items-center z-20 bg-gray-600/70"
        >
            <DataErrorBlock
                v-if="error"
                :title="$t('components.centrum.dispatch_center.failed_datastream')"
                :retry="startStream"
            />
            <DataPendingBlock
                v-else-if="abort === undefined"
                :message="$t('components.centrum.dispatch_center.starting_datastream')"
            />
        </div>

        <div class="relative w-full h-full z-0 flex">
            <!-- Left column -->
            <div class="flex flex-col basis-1/3 divide-x divide-x-reverse divide-base-400 divide-y divide-base-400">
                <div class="basis-11/12">
                    <Livemap>
                        <template v-slot:default>
                            <DispatchesLayer
                                v-if="can('CentrumService.Stream')"
                                :show-all-dispatches="true"
                                @goto="goto($event)"
                            />
                        </template>
                    </Livemap>
                </div>
                <div class="basis-1/12">
                    <DisponentsInfo
                        :disponents="disponents"
                        :settings="settings"
                        :is-disponent="isDisponent"
                        :class="!isDisponent ? 'z-50' : ''"
                    />
                </div>
            </div>

            <!-- Right column -->
            <div class="flex flex-col basis-2/3 divide-y divide-base-400">
                <div class="basis-[55%] max-h-[55%]">
                    <DispatchesList @goto="goto($event)" />
                </div>

                <div class="basis-[35%] max-h-[35%]">
                    <UnitsList @goto="goto($event)" />
                </div>

                <div class="basis-[10%] max-h-[10%]">
                    <Feed :items="feed" />
                </div>
            </div>
        </div>
    </div>
</template>
