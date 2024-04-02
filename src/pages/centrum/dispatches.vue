<script lang="ts" setup>
import { watchDebounced } from '@vueuse/core';
import { Pane, Splitpanes } from 'splitpanes';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import MapTempMarker from '~/components/livemap/MapTempMarker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/store/livemap';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';

useHead({
    title: 'common.dispatches',
});
definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
});

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { location, showLocationMarker } = storeToRefs(livemapStore);

const query = ref<{ postal?: string; id: string }>({
    postal: '',
    id: '',
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`centrum-dispatches-${page.value}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: offset.value,
            },
            notStatus: [],
            status: [],
            ids: [],
            postal: query.value.postal?.trim().replaceAll('-', '').replace(/\D/g, ''),
        };

        if (query.value.id) {
            const id = query.value.id?.trim().replaceAll('-', '').replace(/\D/g, '');
            if (id.length > 0) {
                req.ids.push(id);
            }
        }

        const call = $grpc.getCentrumClient().listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

watchDebounced(query.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});

onMounted(() => {
    showLocationMarker.value = true;
});

onBeforeUnmount(() => {
    showLocationMarker.value = false;
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.dispatches')"> </UDashboardNavbar>

            <Splitpanes class="relative min-h-[calc(100vh-var(--header-height))]">
                <Pane min-size="25">
                    <BaseMap :map-options="{ zoomControl: false }">
                        <template #default>
                            <MapTempMarker />
                        </template>
                    </BaseMap>
                </Pane>
                <Pane size="65">
                    <div class="py-2 pb-14">
                        <div class="px-1 sm:px-2 lg:px-4">
                            <div class="flex max-h-full flex-col overflow-y-auto">
                                <div class="border-b-2 border-neutral/20 pb-2 sm:flex sm:items-center">
                                    <div class="sm:flex-auto">
                                        <form @submit.prevent="refresh()">
                                            <div class="mx-auto flex flex-row gap-4">
                                                <div class="flex-1">
                                                    <label for="search" class="block text-sm font-medium leading-6">
                                                        {{ $t('common.postal') }}
                                                    </label>
                                                    <div class="relative mt-2">
                                                        <UInput
                                                            ref="searchInput"
                                                            v-model="query.postal"
                                                            type="text"
                                                            :placeholder="$t('common.postal')"
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        />
                                                    </div>
                                                </div>
                                                <div class="flex-1">
                                                    <label for="model" class="block text-sm font-medium leading-6">
                                                        {{ $t('common.id') }}
                                                    </label>
                                                    <div class="relative mt-2">
                                                        <UInput
                                                            v-model="query.id"
                                                            type="text"
                                                            name="model"
                                                            min="1"
                                                            max="99999999999"
                                                            :placeholder="$t('common.id')"
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        />
                                                    </div>
                                                </div>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                        </div>

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
                        </template>

                        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
                            <UPagination
                                v-model="page"
                                :page-count="data?.pagination?.pageSize ?? 0"
                                :total="data?.pagination?.totalCount ?? 0"
                            />
                        </div>
                    </div>
                </Pane>
            </Splitpanes>
        </UDashboardPanel>
    </UDashboardPage>
</template>
