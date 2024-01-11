<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { Pane, Splitpanes } from 'splitpanes';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import MapTempMarker from '~/components/livemap/MapTempMarker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useLivemapStore } from '~/store/livemap';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';

useHead({
    title: 'common.dispatches',
});
definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
    showQuickButtons: false,
});

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { location, showLocationMarker } = storeToRefs(livemapStore);

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`centrum-dispatches-${offset.value}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: offset.value,
            },
            notStatus: [],
            status: [],
        };

        const call = $grpc.getCentrumClient().listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

onMounted(() => {
    showLocationMarker.value = true;
});

onBeforeUnmount(() => {
    showLocationMarker.value = false;
});
</script>

<template>
    <div class="relative h-full w-full">
        <Splitpanes class="h-full w-full">
            <Pane min-size="25">
                <BaseMap :map-options="{ zoomControl: false }">
                    <template #default>
                        <MapTempMarker />
                    </template>
                </BaseMap>
            </Pane>
            <Pane size="65">
                <div class="max-h-full overflow-y-auto flex flex-col">
                    <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.dispatches')])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.dispatches')])"
                        :retry="refresh"
                    />
                    <DataNoDataBlock v-else-if="data?.dispatches.length === 0" :type="$t('common.dispatches')" />

                    <template v-else>
                        <DispatchList
                            :show-button="false"
                            :hide-actions="true"
                            :always-show-day="true"
                            :dispatches="data?.dispatches"
                            @goto="location = $event"
                        />

                        <TablePagination :pagination="data?.pagination" :refresh="refresh" @offset-change="offset = $event" />
                    </template>
                </div>
            </Pane>
        </Splitpanes>
    </div>
</template>
