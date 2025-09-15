<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import Feed from '~/components/centrum/Feed.vue';
import DispatchLayer from '~/components/centrum/livemap/DispatchLayer.vue';
import UnitList from '~/components/centrum/units/UnitList.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import DispatcherInfo from './dispatchers/DispatcherInfo.vue';

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
                <Splitpanes>
                    <Pane :min-size="25">
                        <div class="relative size-full">
                            <div v-if="error" class="absolute inset-0 z-30 flex items-center justify-center bg-default/70">
                                <DataErrorBlock
                                    :title="$t('components.centrum.dispatch_center.failed_datastream')"
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

                    <Pane :min-size="40" :size="70">
                        <Splitpanes horizontal>
                            <Pane :size="58" :min-size="2">
                                <DispatchList show-button />
                            </Pane>

                            <Pane :size="26" :min-size="2">
                                <UnitList />
                            </Pane>

                            <Pane :size="8" :min-size="2">
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
