<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import CentrumFeed from '~/components/centrum/CentrumFeed.vue';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import MarkersList from '~/components/centrum/MarkersList.vue';
import UnitList from '~/components/centrum/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import DispatchersInfo from './dispatchers/DispatchersInfo.vue';

const { can } = useAuth();

const centrumStore = useCentrumStore();
const { error, feed, isCenter } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

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
    <UDashboardPanel grow>
        <UDashboardNavbar :title="$t('common.dispatch_center')">
            <template #right>
                <ClientOnly>
                    <DispatchersInfo />
                </ClientOnly>
            </template>
        </UDashboardNavbar>

        <div class="max-h-[calc(100dvh-var(--header-height))] min-h-[calc(100dvh-var(--header-height))] w-full overflow-hidden">
            <Splitpanes>
                <Pane :min-size="25">
                    <div class="relative size-full">
                        <div v-if="error" class="absolute inset-0 z-30 flex items-center justify-center bg-gray-600/70">
                            <DataErrorBlock
                                :title="$t('components.centrum.dispatch_center.failed_datastream')"
                                :error="error"
                                :retry="startStream"
                            />
                        </div>

                        <LivemapBase :show-unit-names="true" :show-unit-status="true">
                            <template #default>
                                <DispatchesLayer
                                    v-if="can('centrum.CentrumService/Stream').value"
                                    :show-all-dispatches="true"
                                />
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
