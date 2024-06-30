<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import UnitList from '~/components/centrum/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import CentrumFeed from '~/components/centrum/CentrumFeed.vue';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import MarkersList from '~/components/centrum/MarkersList.vue';
import DisponentsInfo from './disponents/DisponentsInfo.vue';
import StreamControl from './StreamControl.vue';

const centrumStore = useCentrumStore();
const { error, abort, reconnecting, feed } = storeToRefs(centrumStore);
const { startStream } = centrumStore;

onMounted(async () => useTimeoutFn(async () => startStream(true), 250));

const mount = ref(false);
onMounted(async () => useTimeoutFn(() => (mount.value = true), 35));
</script>

<template>
    <UDashboardPanel grow>
        <UDashboardNavbar :title="$t('common.dispatch_center')">
            <template #center>
                <StreamControl />
            </template>

            <template #right>
                <DisponentsInfo />
            </template>
        </UDashboardNavbar>

        <div
            ref="splitpanesContainer"
            class="max-h-[calc(100vh-var(--header-height))] min-h-[calc(100vh-var(--header-height))] w-full overflow-hidden"
        >
            <Splitpanes v-if="mount" class="relative">
                <Pane :min-size="25">
                    <div class="relative z-0 size-full">
                        <div
                            v-if="error || (abort === undefined && !reconnecting)"
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

                        <LivemapBase :show-unit-names="true" :show-unit-status="true">
                            <template #default>
                                <DispatchesLayer v-if="can('CentrumService.Stream').value" :show-all-dispatches="true" />
                            </template>
                        </LivemapBase>
                    </div>
                </Pane>
                <Pane :min-size="40" :size="70">
                    <Splitpanes horizontal>
                        <Pane :size="58" :min-size="2">
                            <DispatchList :show-button="true" />
                        </Pane>
                        <Pane :size="26" :min-size="2">
                            <UnitList />
                        </Pane>
                        <Pane :size="8" :min-size="2">
                            <MarkersList />
                        </Pane>
                        <Pane :size="8" :min-size="2">
                            <CentrumFeed :items="feed" />
                        </Pane>
                    </Splitpanes>
                </Pane>
            </Splitpanes>
        </div>
    </UDashboardPanel>
</template>

<style>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 2px;
    background-color: rgb(var(--color-gray-800));
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 2px;
    background-color: rgb(var(--color-gray-800));
}
</style>
