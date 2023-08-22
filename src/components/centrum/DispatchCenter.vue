<script lang="ts" setup>
import { default as DispatchDetails } from '~/components/centrum/dispatches/Details.vue';
import { default as DispatchesList } from '~/components/centrum/dispatches/List.vue';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import { default as UnitsList } from '~/components/centrum/units/List.vue';
import Livemap from '~/components/livemap/Livemap.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import DisponentsInfo from './DisponentsInfo.vue';
import Feed from './Feed.vue';
import AssignDispatchModal from './dispatches/AssignDispatchModal.vue';
import { default as DispatchStatusUpdateModal } from './dispatches/StatusUpdateModal.vue';
import DispatchesLayer from './livemap/DispatchesLayer.vue';
import { setWaypoint } from './nui';
import AssignUnitModal from './units/AssignUnitModal.vue';
import { default as UnitStatusUpdateModal } from './units/StatusUpdateModal.vue';

const centrumStore = useCentrumStore();
const { error, abort, isDisponent, disponents, settings, feed } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;
const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

onMounted(() => {
    startStream();
});

onBeforeUnmount(() => {
    stopStream();
});

function goto(e: Coordinate) {
    location.value = { x: e.x, y: e.y };

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}

const selectedDispatch = ref<Dispatch | undefined>();
const openDispatchDetails = ref(false);
const openDispatchAssign = ref(false);
const openDispatchStatus = ref(false);

const selectedUnit = ref<Unit | undefined>();
const openUnitDetails = ref(false);
const openUnitAssign = ref(false);
const openUnitStatus = ref(false);
</script>

<template>
    <div class="flex-col h-full relative">
        <div v-if="error || abort === undefined" class="absolute inset-0 flex justify-center items-center z-20 bg-gray-600/70">
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
                                @select="
                                    selectedDispatch = $event;
                                    openDispatchDetails = true;
                                "
                                :show-all-dispatches="true"
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
                <div class="basis-3/5 max-h-[60%]">
                    <DispatchesList
                        @goto="goto($event)"
                        @details="
                            selectedDispatch = $event;
                            openDispatchDetails = true;
                        "
                        @assign-unit="
                            selectedDispatch = $event;
                            openDispatchAssign = true;
                        "
                        @status="
                            selectedDispatch = $event;
                            openDispatchStatus = true;
                        "
                    />

                    <template v-if="selectedDispatch">
                        <DispatchDetails
                            :dispatch="selectedDispatch"
                            :open="openDispatchDetails"
                            @close="openDispatchDetails = false"
                            @goto="goto($event)"
                            @assign-unit="
                                selectedDispatch = $event;
                                openDispatchAssign = true;
                            "
                            @status="
                                selectedDispatch = $event;
                                openDispatchStatus = true;
                            "
                        />
                        <AssignDispatchModal
                            :open="openDispatchAssign"
                            :dispatch="selectedDispatch"
                            @close="openDispatchAssign = false"
                        />
                        <DispatchStatusUpdateModal
                            :open="openDispatchStatus"
                            :dispatch="selectedDispatch"
                            @close="openDispatchStatus = false"
                        />
                    </template>
                </div>
                <div class="basis-1/5 max-h-[20%]">
                    <Feed :items="feed" />
                </div>
                <div class="basis-1/5 max-h-[20%]">
                    <UnitsList
                        @goto="goto($event)"
                        @details="
                            selectedUnit = $event;
                            openUnitDetails = true;
                        "
                    />

                    <template v-if="selectedUnit">
                        <UnitDetails
                            :unit="selectedUnit"
                            :open="openUnitDetails"
                            @close="openUnitDetails = false"
                            @goto="goto($event)"
                            @assign-users="
                                selectedUnit = $event;
                                openUnitAssign = true;
                            "
                            @status="
                                selectedUnit = $event;
                                openUnitStatus = true;
                            "
                        />
                        <AssignUnitModal :open="openUnitAssign" :unit="selectedUnit" @close="openUnitAssign = false" />
                        <UnitStatusUpdateModal :open="openUnitStatus" :unit="selectedUnit" @close="openUnitStatus = false" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
