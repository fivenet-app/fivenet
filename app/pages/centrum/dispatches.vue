<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import { z } from 'zod';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import Pagination from '~/components/partials/Pagination.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/stores/livemap';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';

useHead({
    title: 'common.dispatches',
});

definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'centrum.CentrumService/TakeControl',
});

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { showLocationMarker } = storeToRefs(livemapStore);

const schema = z.object({
    postal: z.string().trim().max(12).default(''),
    id: z.coerce.number().max(16).default(0),
});

const query = useSearchForm('centrum_dispatches_archive', schema);

const page = useRouteQuery('page', '1', { transform: Number });

const { data, pending: loading, refresh, error } = useLazyAsyncData(`centrum-dispatches-${page.value}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            notStatus: [],
            status: [],
            ids: [],
            postal: query.postal.replaceAll('-', '').replace(/\D/g, ''),
        };

        if (query.id && query.id > 0) {
            req.ids.push(query.id);
        }

        const call = $grpc.centrum.centrum.listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

onBeforeMount(() => (showLocationMarker.value = true));
onMounted(async () => useTimeoutFn(() => (mount.value = true), 35));

onBeforeUnmount(() => {
    showLocationMarker.value = false;
});

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});

const mount = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.dispatches')">
                <template #right>
                    <PartialsBackButton fallback-to="/centrum" />
                </template>
            </UDashboardNavbar>

            <div class="max-h-[calc(100dvh-var(--header-height))] min-h-[calc(100dvh-var(--header-height))] overflow-hidden">
                <Splitpanes v-if="mount" class="relative">
                    <Pane :min-size="25">
                        <ClientOnly>
                            <BaseMap :map-options="{ zoomControl: false }">
                                <template #default>
                                    <LazyLivemapMapTempMarker />

                                    <DispatchesLayer show-all-dispatches :dispatch-list="data?.dispatches ?? []" />
                                </template>
                            </BaseMap>
                        </ClientOnly>
                    </Pane>

                    <Pane :size="65">
                        <div class="max-h-full overflow-y-auto">
                            <div class="mb-2 px-2">
                                <UForm class="flex flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                                    <UFormGroup class="flex-1" name="postal" :label="$t('common.postal')">
                                        <UInput
                                            ref="input"
                                            v-model="query.postal"
                                            type="text"
                                            name="postal"
                                            :placeholder="$t('common.postal')"
                                        >
                                            <template #trailing>
                                                <UKbd value="/" />
                                            </template>
                                        </UInput>
                                    </UFormGroup>
                                    <UFormGroup class="flex-1" name="id" :label="$t('common.id')">
                                        <UInput
                                            v-model="query.id"
                                            type="text"
                                            name="id"
                                            :min="1"
                                            :max="99999999999"
                                            :placeholder="$t('common.id')"
                                        />
                                    </UFormGroup>
                                </UForm>
                            </div>

                            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.dispatches')])" />
                            <DataErrorBlock
                                v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.dispatches')])"
                                :error="error"
                                :retry="refresh"
                            />
                            <DataNoDataBlock v-else-if="data?.dispatches.length === 0" :type="$t('common.dispatches')" />

                            <div v-else class="relative overflow-x-auto">
                                <DispatchList
                                    :show-button="false"
                                    :hide-actions="true"
                                    :always-show-day="true"
                                    :dispatches="data?.dispatches"
                                />
                            </div>

                            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
                        </div>
                    </Pane>
                </Splitpanes>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>

<style scoped>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 2px;
    background-color: rgb(var(--color-gray-800));
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 2px;
    background-color: rgb(var(--color-gray-800));
}
</style>
